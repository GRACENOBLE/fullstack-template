---
topic: error-handling
last_verified: 2026-06-25
sources:
  - internal/infrastructure/database/postgres/health_repository.go
  - internal/transport/handlers/health_handler.go
  - internal/transport/handlers/validation.go
  - internal/transport/handlers/response.go
  - cmd/api/main.go
---

# Error Handling

## Response envelope

All handler responses use helpers defined in `internal/transport/handlers/response.go`. Never call `c.JSON` directly in a handler.

**Success shape:**
```json
{"data": <payload>}
```

**Error shape:**
```json
{"error": {"code": "snake_case_code", "message": "human readable message"}}
```

**Helper signatures:**
```go
// 200 with data wrapped in {"data": ...}
func JSON[T any](c *gin.Context, data T)

// Any status code with data wrapped in {"data": ...}
func JSONStatus[T any](c *gin.Context, status int, data T)

// Error response as {"error": {"code": "...", "message": "..."}}
func JSONError(c *gin.Context, status int, code, message string)
```

Use `JSON` for standard 200 responses, `JSONStatus` when a non-200 success status is needed (e.g., 201 Created), and `JSONError` for all error responses.

## General rule
Return errors up the call stack. Callers decide how to handle them.
Never use `log.Fatal` or `os.Exit` inside `internal/`.

## Documented exception (intentional)
| Location | Call | Reason |
|---|---|---|
| `cmd/api/main.go: main()` | `fmt.Fprintf(os.Stderr, ...) + os.Exit(1)` | `bootstrap.Run()` returned an error — process cannot start |

This is the only permitted early-exit path and it lives in `cmd/`, not `internal/`.
`server.NewServer` returns `(*http.Server, error)` — the caller in `cmd/api/main.go` checks the error and exits on failure. Fallible startup work is split between `bootstrap.Run` and `server.NewServer` (e.g. registering Prometheus collectors).

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
Handlers call use cases, check errors, and map them to HTTP status codes using the `JSONError` helper. The health handler returns 503 when the DB is unreachable:

```go
func (h *Handler) healthHandler(c *gin.Context) {
    stats, err := h.healthUC.GetHealth(c.Request.Context())
    if err != nil {
        log.Printf("health check failed: %v", err)
        JSONStatus(c, http.StatusServiceUnavailable, stats)
        return
    }
    JSON(c, stats)
}
```

For general handlers, map errors to status codes explicitly:

```go
func (h *Handler) getItemHandler(c *gin.Context) {
    item, err := h.itemUC.GetItem(c.Request.Context(), id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            JSONError(c, http.StatusNotFound, "not_found", "not found")
            return
        }
        JSONError(c, http.StatusInternalServerError, "internal_error", "internal error")
        return
    }
    JSON(c, item)
}
```

Never expose internal error messages to clients. Log the original error server-side.

## Request binding errors
Use the shared helpers in `internal/transport/handlers/validation.go` — do not call `c.ShouldBindJSON` / `c.ShouldBindQuery` directly in handlers.

```go
// internal/transport/handlers/validation.go
func bindJSON(c *gin.Context, dst any) bool   // writes 400 {"error":"invalid request body"} on failure
func bindQuery(c *gin.Context, dst any) bool  // writes 400 {"error":"invalid query parameters"} on failure
```

Both helpers return `false` and write the 400 response on failure, so the handler just returns immediately:

```go
var req updateMeRequest
if !bindJSON(c, &req) {
    return
}
```

The stable message strings (`"invalid request body"`, `"invalid query parameters"`) never expose raw validation errors to the client.

## Error wrapping
Use `fmt.Errorf("context: %w", err)` when adding context to returned errors so callers can use `errors.Is` / `errors.As`.

## Panic recovery
`gin.Recovery()` is applied explicitly in `RegisterRoutes()` — panics in handlers are recovered and return 500. Do not rely on this; handle errors explicitly.
