---
topic: error-handling
last_verified: 2026-06-14
sources:
  - internal/database/database.go
  - internal/server/routes.go
  - cmd/api/main.go
---

# Error Handling

## General rule
Return errors up the call stack. Callers decide how to handle them.
Never use `log.Fatal` or `os.Exit` inside `internal/` except for the two documented exceptions below.

## Documented exceptions (intentional)
| Location | Call | Reason |
|---|---|---|
| `database.go: New()` | `log.Fatal(err)` | Startup failure — if the DB can't connect at boot, the process should not continue |
| `database.go: Health()` | `log.Fatalf(...)` | Unrecoverable DB loss detected during health check — intentional termination |

These are startup/health-check paths only. All other paths in `internal/` return errors.

## Handler error responses
In Gin handlers, map errors to HTTP status codes explicitly:

```go
func (s *Server) getItemHandler(c *gin.Context) {
    item, err := s.db.GetItem(c.Request.Context(), id)
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
