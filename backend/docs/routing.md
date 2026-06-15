---
topic: routing
last_verified: 2026-06-15
sources:
  - internal/transport/handlers/handler.go
  - internal/transport/handlers/routes.go
  - internal/transport/handlers/hello_handler.go
  - internal/transport/handlers/health_handler.go
  - internal/transport/middleware/logger.go
  - internal/server/server.go
  - cmd/api/main.go
---

# Routing

## Handler struct
```go
// internal/handler/handler.go
type Handler struct {
    healthUC usecase.HealthUseCase
}

func NewHandler(healthUC usecase.HealthUseCase) *Handler {
    return &Handler{healthUC: healthUC}
}
```
The `Handler` struct holds use case interfaces — not `*sql.DB` directly. Add new use case fields here as features are added.

## Wiring (server.go)
`internal/server/server.go` contains `NewServer(app *bootstrap.App) *http.Server` — wiring only, no logic.
It receives the already-validated `*bootstrap.App` (which holds `*sql.DB` and `Config`), constructs the repository, use case, and handler in order, then returns a configured `*http.Server`. It does not read env vars or return an error.

```go
healthRepo := postgres.NewHealthRepository(app.DB)
healthUC := usecase.NewHealthUseCase(healthRepo)
h := handler.NewHandler(healthUC)

return &http.Server{
    Addr:         fmt.Sprintf(":%d", app.Config.Port),
    Handler:      h.RegisterRoutes(app.Config.RateLimitRPS, app.Config.RateLimitBurst),
    IdleTimeout:  time.Minute,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 30 * time.Second,
}
```

## Route registration
All routes registered in `RegisterRoutes()` on `*Handler`, which returns `http.Handler`.
`rps` and `burst` come from `bootstrap.Config` (env vars `RATE_LIMIT_RPS` / `RATE_LIMIT_BURST`); pass `rps=0` to disable.

```go
func (h *Handler) RegisterRoutes(rps float64, burst int) http.Handler {
    r := gin.New()

    // Gin's colorful logger locally; structured slog logger in staging/production.
    if gin.Mode() == gin.DebugMode {
        r.Use(gin.Recovery(), gin.Logger())
    } else {
        r.Use(gin.Recovery(), middleware.Logger())
    }

    r.Use(middleware.RateLimit(rps, burst))

    r.Use(cors.New(cors.Config{ ... }))

    r.GET("/path", h.myHandler)
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
| Method | Path | Handler | File |
|---|---|---|---|
| GET | `/` | `HelloWorldHandler` — returns `{"message": "Hello World"}` | `hello_handler.go` |
| GET | `/health` | `healthHandler` — returns `HealthStats`; 503 when DB is down | `health_handler.go` |

## Graceful shutdown
Wired in `cmd/api/main.go` via `signal.NotifyContext` for SIGINT/SIGTERM.
5-second shutdown timeout. Server notifies `done chan bool` when complete.
`main()` calls `bootstrap.Run(ctx)` first; on failure it writes to stderr and calls `os.Exit(1)`.
`server.NewServer` does not return an error — all fallible startup work is in bootstrap.
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
