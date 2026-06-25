---
topic: middleware
last_verified: 2026-06-25
sources:
  - internal/transport/middleware/logger.go
  - internal/transport/middleware/ratelimit.go
  - internal/transport/middleware/auth.go
  - internal/transport/middleware/metrics.go
  - internal/transport/middleware/local_network.go
  - internal/transport/middleware/geo.go
  - internal/transport/middleware/request_id.go
  - internal/transport/handlers/routes.go
---

# Middleware

All middleware lives in `internal/transport/middleware/` and follows the Gin `HandlerFunc` pattern. Middleware is registered in `RegisterRoutes()` inside `internal/transport/handlers/routes.go`.

## Registration order

```go
// 1. Request ID — must be first so every subsequent middleware can read the ID
r.Use(middleware.RequestID())
// 2. Sentry error reporting
r.Use(middleware.SentryMiddleware(sentryDSN))
// 3. Recovery + logger (debug: gin.Logger, release: middleware.Logger)
r.Use(gin.Recovery(), middleware.Logger())
// 4. Prometheus metrics collection
r.Use(middleware.PrometheusMiddleware())
// 5. Rate limiter (no-op when RPS <= 0)
r.Use(middleware.RateLimit(rps, burst))
// 6. CORS
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

## RequestID

`RequestID() gin.HandlerFunc` assigns a unique identifier to every request. It is registered as the first middleware in `RegisterRoutes` so all subsequent middleware (including logger and Sentry) have access to the ID.

```go
const RequestIDKey    = "request_id"
const RequestIDHeader = "X-Request-ID"

func RequestID() gin.HandlerFunc
```

Behaviour:
- Reads the `X-Request-ID` request header. If present and non-empty, uses that value (allows callers to propagate their own trace IDs).
- If absent or empty, generates a random 16-byte hex string (`crypto/rand`).
- Stores the ID in the Gin context under `RequestIDKey` via `c.Set`.
- Echoes the ID back in the `X-Request-ID` response header.

Reading the ID inside a handler or middleware:
```go
requestID := c.GetString(middleware.RequestIDKey)
```

The structured `Logger()` middleware appends `"request_id"` to every slog record automatically.

## Logger

`Logger() gin.HandlerFunc` emits one structured `slog` record per request after `c.Next()` returns. Fields: `status`, `method`, `path`, `latency`, `ip`, `request_id`, and optionally `query` and `errors`.

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

## PrometheusMiddleware

`PrometheusMiddleware() gin.HandlerFunc` records two metrics for every request except `/metrics` itself:

- `http_requests_total` — counter vector with labels `method`, `path`, `status`.
- `http_request_duration_seconds` — histogram vector with labels `method`, `path`.

Unmatched routes (404s with no Gin `FullPath()`) are recorded under the path label `"unmatched"`.

## LocalNetworkOnly

`LocalNetworkOnly() gin.HandlerFunc` aborts with `403 Forbidden` when the client IP is neither a loopback address nor an RFC 1918 private address. `RegisterRoutes` applies it in two places:

1. `/metrics` — in release mode only, so the Prometheus scrape endpoint is reachable from the internal network but not from external clients.
2. `/debug/pprof/*` — unconditionally (both debug and release modes), applied as a group middleware so all pprof endpoints are always restricted to loopback/private addresses.

```go
// release mode only:
r.GET("/metrics", middleware.LocalNetworkOnly(), gin.WrapH(promhttp.Handler()))

// all modes:
debug := r.Group("/debug/pprof", middleware.LocalNetworkOnly())
```

## GeoFromRequest

`GeoFromRequest(locator usecase.GeoLocator) gin.HandlerFunc` resolves the request's originating IP to geographic metadata and stores the result in the Gin context. It is best-effort: any error from `locator.Lookup` (private IP, rate-limit, network failure) is silently dropped and the request continues without geo data.

```go
const GeoLocationKey = "geo_location"

func GeoFromRequest(locator usecase.GeoLocator) gin.HandlerFunc
```

Context key: `middleware.GeoLocationKey` (`"geo_location"`). Value type: `*domain.GeoLocation`.

Reading geo data in a handler:
```go
val, exists := c.Get(middleware.GeoLocationKey)
if !exists {
    // geo unavailable
    return
}
geo, ok := val.(*domain.GeoLocation)
if !ok || geo == nil {
    return
}
```

### RealIP helper

```go
func RealIP(r *http.Request) string
```

IP extraction precedence:
1. First address in `X-Forwarded-For` (proxy/Railway deploys).
2. `X-Real-IP` header.
3. `r.RemoteAddr` with port stripped via `net.SplitHostPort`.

`RealIP` is exported for direct use in tests and other packages.

## Adding new middleware

1. Create `internal/transport/middleware/<name>.go` with a function returning `gin.HandlerFunc`.
2. Register it in `RegisterRoutes()` in `internal/transport/handlers/routes.go` at the appropriate position in the chain.
3. If it requires configuration, add fields to `bootstrap.Config` and read from env in `loadConfig()`.
