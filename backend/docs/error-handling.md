---
topic: error-handling
last_verified: 2026-06-14
sources:
  - internal/repository/postgres/health_repository.go
  - internal/handler/health_handler.go
  - cmd/api/main.go
---

# Error Handling

## General rule
Return errors up the call stack. Callers decide how to handle them.
Never use `log.Fatal` or `os.Exit` inside `internal/`.

## Documented exception (intentional)
| Location | Call | Reason |
|---|---|---|
| `cmd/api/main.go: main()` | `log.Fatalf(...)` | `server.NewServer()` returned an error — process cannot start |

This is the only permitted `log.Fatal` call and it lives in `cmd/`, not `internal/`.

## Repository errors
Repository methods return `(Result, error)`. On failure, wrap with context using `fmt.Errorf`:

```go
func (r *HealthRepository) Health(ctx context.Context) (domain.HealthStats, error) {
    if err := r.db.PingContext(pingCtx); err != nil {
        stats["status"] = "down"
        stats["error"] = fmt.Sprintf("db down: %v", err)
        return stats, fmt.Errorf("postgres: health ping: %w", err)
    }
    // ...
    return stats, nil
}
```

## Handler error responses
Handlers call use cases, check errors, and map them to HTTP status codes. The health handler returns 503 when the DB is unreachable:

```go
func (h *Handler) healthHandler(c *gin.Context) {
    stats, err := h.healthUC.GetHealth(c.Request.Context())
    if err != nil {
        log.Printf("health check failed: %v", err)
        c.JSON(http.StatusServiceUnavailable, stats)
        return
    }
    c.JSON(http.StatusOK, stats)
}
```

For general handlers, map errors to status codes explicitly:

```go
func (h *Handler) getItemHandler(c *gin.Context) {
    item, err := h.itemUC.GetItem(c.Request.Context(), id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
        return
    }
    c.JSON(http.StatusOK, item)
}
```

Never expose internal error messages to clients. Log the original error server-side.

## Request binding errors
Always validate and return 400 on bad input:

```go
var input MyRequest
if err := c.ShouldBindJSON(&input); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
}
```

## Error wrapping
Use `fmt.Errorf("context: %w", err)` when adding context to returned errors so callers can use `errors.Is` / `errors.As`.

## Panic recovery
`gin.Default()` includes the Recovery middleware — panics in handlers are recovered and return 500. Do not rely on this; handle errors explicitly.
