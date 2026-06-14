---
topic: routing
last_verified: 2026-06-14
sources:
  - internal/server/server.go
  - internal/server/routes.go
  - cmd/api/main.go
---

# Routing

## Server struct
```go
type Server struct {
    port int
    db   database.Service
}
```
`NewServer()` reads `PORT` from env, wires up the `Server`, and returns a configured `*http.Server` with timeouts:
- `IdleTimeout`: 1 minute
- `ReadTimeout`: 10 seconds
- `WriteTimeout`: 30 seconds

## Route registration
All routes registered in `RegisterRoutes()` on `*Server`, which returns `http.Handler`.
`RegisterRoutes()` creates the Gin engine, applies middleware, registers routes, and returns.

```go
func (s *Server) RegisterRoutes() http.Handler {
    r := gin.Default()
    r.Use(cors.New(cors.Config{ ... }))
    r.GET("/path", s.myHandler)
    return r
}
```

## Handler pattern
All handlers are methods on `*Server`. Always use `*gin.Context`.

```go
func (s *Server) myHandler(c *gin.Context) {
    // access DB via s.db
    c.JSON(http.StatusOK, gin.H{"key": "value"})
}
```

## CORS configuration
Pre-configured in `RegisterRoutes()` via `github.com/gin-contrib/cors`.
Current allowed origin: `http://localhost:5173` — **update to `http://localhost:3000`** for this project's Next.js frontend.
Allowed methods: GET, POST, PUT, DELETE, OPTIONS, PATCH.
`AllowCredentials: true` — cookies and auth headers pass through.

## Existing routes
| Method | Path | Handler |
|---|---|---|
| GET | `/` | `HelloWorldHandler` — returns `{"message": "Hello World"}` |
| GET | `/health` | `healthHandler` — returns `s.db.Health()` map |

## Graceful shutdown
Wired in `cmd/api/main.go` via `signal.NotifyContext` for SIGINT/SIGTERM.
5-second shutdown timeout. Server notifies `done chan bool` when complete.
Do not add shutdown logic to `internal/` — it belongs in `cmd/`.

## Adding a new route — checklist
1. Register in `RegisterRoutes()`: `r.METHOD("/path", s.handlerName)`
2. Add handler as `func (s *Server) handlerName(c *gin.Context)`
3. If DB is needed, call `s.db.MethodName(c.Request.Context(), ...)` — always pass context
4. For request body: bind with `c.ShouldBindJSON(&input)`, return 400 on error
5. For path params: `c.Param("id")`
6. For query params: `c.Query("key")`
