---
topic: routing
last_verified: 2026-06-14
sources:
  - internal/handler/handler.go
  - internal/handler/routes.go
  - internal/handler/hello_handler.go
  - internal/handler/health_handler.go
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
`internal/server/server.go` contains `NewServer() (*http.Server, error)` — wiring only, no logic.
It builds `DBConfig` from env, calls `NewPostgresDB`, constructs the repository, use case, and handler in order, then returns a configured `*http.Server`.

```go
db, err := postgres.NewPostgresDB(cfg)
healthRepo := postgres.NewHealthRepository(db)
healthUC := usecase.NewHealthUseCase(healthRepo)
h := handler.NewHandler(healthUC)

srv := &http.Server{
    Addr:         fmt.Sprintf(":%d", port),
    Handler:      h.RegisterRoutes(),
    IdleTimeout:  time.Minute,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 30 * time.Second,
}
```

## Route registration
All routes registered in `RegisterRoutes()` on `*Handler`, which returns `http.Handler`.

```go
func (h *Handler) RegisterRoutes() http.Handler {
    r := gin.Default()
    r.Use(cors.New(cors.Config{ ... }))
    r.GET("/path", h.myHandler)
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
`main()` handles the error returned by `server.NewServer()` with `log.Fatalf`.
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
