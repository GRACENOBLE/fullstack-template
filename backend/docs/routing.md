---
topic: routing
last_verified: 2026-06-15
sources:
  - internal/transport/handlers/handler.go
  - internal/transport/handlers/routes.go
  - internal/transport/handlers/hello_handler.go
  - internal/transport/handlers/health_handler.go
  - internal/transport/handlers/auth_handler.go
  - internal/transport/middleware/logger.go
  - internal/server/server.go
  - cmd/api/main.go
---

# Routing

## Handler struct
```go
// internal/transport/handlers/handler.go
type Handler struct {
    healthUC usecase.HealthUseCase
    verifier usecase.FirebaseTokenVerifier  // nil disables auth (dev only)
    hub      *ws.Hub
}

func NewHandler(healthUC usecase.HealthUseCase, verifier usecase.FirebaseTokenVerifier, hub *ws.Hub) *Handler {
    return &Handler{healthUC: healthUC, verifier: verifier, hub: hub}
}
```
The `Handler` struct holds use case interfaces and infrastructure dependencies — not `*sql.DB` directly. `verifier` is stored on the struct (not passed to `RegisterRoutes`) so the WebSocket handler can read it inline for query-param auth.

## Wiring (server.go)
`internal/server/server.go` contains `NewServer(app *bootstrap.App, hub *ws.Hub) (*http.Server, error)` — wiring only, no logic.
It receives the already-validated `*bootstrap.App` (which holds `*sql.DB`, `Cache`, `Firebase`, and `Config`) and a `*ws.Hub`, constructs the repository, use case, and handler in order, then returns a configured `*http.Server`. Errors from initialisation steps are returned to the caller.

```go
healthRepo := postgres.NewHealthRepository(app.DB)
healthUC := usecase.NewHealthUseCase(healthRepo)
h := handlers.NewHandler(healthUC, app.Firebase, hub)

return &http.Server{
    Addr:         fmt.Sprintf(":%d", app.Config.Port),
    Handler:      h.RegisterRoutes(app.Config.RateLimitRPS, app.Config.RateLimitBurst, app.Config.SentryDSN),
    IdleTimeout:  time.Minute,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 30 * time.Second,
}, nil
```

## Route registration
All routes registered in `RegisterRoutes()` on `*Handler`, which returns `http.Handler`.
`rps` and `burst` come from `bootstrap.Config` (env vars `RATE_LIMIT_RPS` / `RATE_LIMIT_BURST`); pass `rps=0` to disable.
`h.verifier` (set via `NewHandler`) controls Firebase auth — the verifier is read from the struct, not passed to `RegisterRoutes`; a `nil` verifier skips Firebase auth (development only — see [auth](auth.md)).

```go
func (h *Handler) RegisterRoutes(rps float64, burst int, sentryDSN string) http.Handler {
    r := gin.New()

    // Gin's colorful logger locally; structured slog logger in staging/production.
    if gin.Mode() == gin.DebugMode {
        r.Use(gin.Recovery(), gin.Logger())
    } else {
        r.Use(gin.Recovery(), middleware.Logger())
    }

    r.Use(middleware.RateLimit(rps, burst))

    r.Use(cors.New(cors.Config{ ... }))

    r.GET("/", h.HelloWorldHandler)
    r.GET("/health", h.HealthHandler)
    r.GET("/ws", h.WsHandler)          // WebSocket upgrade — auth via ?token= query param

    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    api := r.Group("/api/v1")
    if h.verifier != nil {
        api.Use(middleware.FirebaseAuth(h.verifier))
    }
    api.GET("/me", h.MeHandler)

    return r
}
```

## Handler pattern
All handlers are methods on `*Handler`. Always use `*gin.Context`.

```go
func (h *Handler) myHandler(c *gin.Context) {
    result, err := h.someUC.DoSomething(c.Request.Context(), ...)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
        return
    }
    c.JSON(http.StatusOK, result)
}
```

## CORS configuration
Pre-configured in `RegisterRoutes()` via `github.com/gin-contrib/cors`.
Current allowed origin: `http://localhost:3000`.
Allowed methods: GET, POST, PUT, DELETE, OPTIONS, PATCH.
`AllowCredentials: true` — cookies and auth headers pass through.

## Existing routes
| Method | Path | Auth | Handler | File |
|---|---|---|---|---|
| GET | `/` | none | `HelloWorldHandler` — returns `{"message": "Hello World"}` | `hello_handler.go` |
| GET | `/health` | none | `HealthHandler` — returns `HealthStats`; 503 when DB is down | `health_handler.go` |
| GET | `/ws` | `?token=` query param | `WsHandler` — upgrades to WebSocket; 401 when token missing/invalid | `ws_handler.go` |
| GET | `/metrics` | `LocalNetworkOnly()` | Prometheus scrape endpoint; restricted to loopback/RFC 1918 in staging/production | `metrics_handler.go` |
| GET | `/api/v1/me` | FirebaseAuth header | `MeHandler` — returns verified `FirebaseToken` claims | `auth_handler.go` |

## Graceful shutdown
Wired in `cmd/api/main.go` via `signal.NotifyContext` for SIGINT/SIGTERM.
5-second shutdown timeout. Server notifies `done chan bool` when complete.
`main()` calls `bootstrap.Run(ctx)` first; on failure it writes to stderr and calls `os.Exit(1)`.
`server.NewServer` returns `(*http.Server, error)` — the caller in `cmd/api/main.go` checks the error and exits on failure.
Do not add shutdown logic to `internal/` — it belongs in `cmd/`.

## Adding a new route — checklist
1. Define domain type in `internal/domain/` if needed.
2. Add use case interface + implementation in `internal/usecase/`.
3. Add repository interface in the use case package; implement in `internal/repository/postgres/`.
4. Add use case field to `Handler` struct in `handler.go`; update `NewHandler` signature.
5. Wire the new repository → use case → handler in `server.NewServer()`.
6. Register the route in `RegisterRoutes()`: `r.METHOD("/path", h.handlerName)`.
7. Add handler as `func (h *Handler) handlerName(c *gin.Context)` in its own file.
8. Always pass `c.Request.Context()` to use case calls.
9. For request body: bind with `c.ShouldBindJSON(&input)`, return 400 on error.
10. For path params: `c.Param("id")`. For query params: `c.Query("key")`.
