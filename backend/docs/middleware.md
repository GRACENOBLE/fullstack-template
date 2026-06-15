---
topic: middleware
last_verified: 2026-06-15
sources:
  - internal/transport/middleware/logger.go
  - internal/transport/middleware/ratelimit.go
  - internal/transport/middleware/auth.go
  - internal/transport/handlers/routes.go
---

# Middleware

All middleware lives in `internal/transport/middleware/` and follows the Gin `HandlerFunc` pattern. Middleware is registered in `RegisterRoutes()` inside `internal/transport/handlers/routes.go`.

## Registration order

```go
// 1. Sentry error reporting
r.Use(middleware.SentryMiddleware(sentryDSN))
// 2. Recovery + logger (debug: gin.Logger, release: middleware.Logger)
r.Use(gin.Recovery(), middleware.Logger())
// 3. Prometheus metrics collection
r.Use(middleware.PrometheusMiddleware())
// 4. Rate limiter (no-op when RPS <= 0)
r.Use(middleware.RateLimit(rps, burst))
// 5. CORS
r.Use(cors.New(...))

// Global routes (no auth):
r.GET("/", h.HelloWorldHandler)
r.GET("/health", h.HealthHandler)
r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

// Protected group — FirebaseAuth applied when h.verifier != nil:
api := r.Group("/api/v1")
if h.verifier != nil {
    api.Use(middleware.FirebaseAuth(h.verifier))
}
api.GET("/me", h.MeHandler)
```

## Logger

`Logger() gin.HandlerFunc` emits one structured `slog` record per request after `c.Next()` returns. Fields: `status`, `method`, `path`, `latency`, `ip`, and optionally `query` and `errors`.

In debug mode (`ENV` not set to `staging`/`production`) Gin's built-in colorful logger is used instead.

## Rate limiter

`RateLimit(rps float64, burst int) gin.HandlerFunc` limits each client IP to `rps` requests per second using a token-bucket algorithm (`golang.org/x/time/rate`). Each IP gets its own `rate.Limiter` stored in a mutex-guarded map.

- **Disabled** when `rps <= 0` (no-op middleware returned).
- **429 Too Many Requests** returned when the bucket is empty; body: `{"error": "rate limit exceeded"}`.
- Configured via env vars `RATE_LIMIT_RPS` and `RATE_LIMIT_BURST` (see [environment](environment.md)).

## FirebaseAuth

`FirebaseAuth(verifier usecase.FirebaseTokenVerifier) gin.HandlerFunc` validates a Firebase ID token on every request to the routes it guards.

```go
const FirebaseClaimsKey = "firebase_claims"

func FirebaseAuth(verifier usecase.FirebaseTokenVerifier) gin.HandlerFunc
```

Behaviour:
- Expects `Authorization: Bearer <firebase-id-token>` header.
- Calls `verifier.VerifyIDToken(ctx, idToken)` — the concrete implementation is `pkg/firebase.authClientAdapter`.
- On success: stores `*usecase.FirebaseToken` in the Gin context under `FirebaseClaimsKey` and calls `c.Next()`.
- On failure (missing header, malformed header, or token verification error): aborts with `401 Unauthorized` and a JSON body `{"error": "..."}`.

Retrieve verified claims inside a handler:
```go
val, _ := c.Get(middleware.FirebaseClaimsKey)
token, ok := val.(*usecase.FirebaseToken)
```

Pass `nil` as the `verifier` to `NewHandler` to skip Firebase auth entirely (development without credentials). `RegisterRoutes` reads `h.verifier` from the struct — it is not a parameter of `RegisterRoutes`.

## Adding new middleware

1. Create `internal/transport/middleware/<name>.go` with a function returning `gin.HandlerFunc`.
2. Register it in `RegisterRoutes()` in `internal/transport/handlers/routes.go` at the appropriate position in the chain.
3. If it requires configuration, add fields to `bootstrap.Config` and read from env in `loadConfig()`.
