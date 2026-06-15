---
topic: observability
last_verified: 2026-06-15
sources:
  - internal/transport/middleware/sentry.go
  - internal/bootstrap/bootstrap.go
  - internal/transport/handlers/routes.go
---

# Observability

## Sentry SDK

| Package | Version |
|---|---|
| `github.com/getsentry/sentry-go` | v0.46.2 |
| `github.com/getsentry/sentry-go/gin` | v0.46.2 |

## SentryMiddleware

`internal/transport/middleware/sentry.go` exports a single function:

```go
func SentryMiddleware(dsn string) gin.HandlerFunc
```

Behavior:
- When `dsn` is empty, returns `func(c *gin.Context) { c.Next() }` — a no-op that adds no overhead.
- When `dsn` is non-empty, calls `sentry.Init(sentry.ClientOptions{Dsn: dsn})` and returns `sentrygin.New(sentrygin.Options{Repanic: true})`.
- `Repanic: true` means the middleware re-panics after capturing, allowing Gin's `Recovery()` to handle the panic response normally.

## Middleware registration order

`RegisterRoutes` in `internal/transport/handlers/routes.go` registers middleware in this order:

1. `SentryMiddleware(sentryDSN)` — first, so it wraps all subsequent handlers
2. `gin.Recovery()` + `gin.Logger()` (debug) or `gin.Recovery()` + `middleware.Logger()` (non-debug)
3. `middleware.RateLimit(rps, burst)`
4. CORS

`RegisterRoutes` signature:

```go
func (h *Handler) RegisterRoutes(rps float64, burst int, sentryDSN string) http.Handler
```

The `sentryDSN` parameter is forwarded directly from `Config.SentryDSN`.

## Environment variable

| Variable | Required | Default |
|---|---|---|
| `SENTRY_DSN` | No | `""` (Sentry disabled) |

Loaded in `loadConfig()` in `internal/bootstrap/bootstrap.go`:

```go
SentryDSN: os.Getenv("SENTRY_DSN"),
```

Stored on `Config.SentryDSN`. Not validated — an empty value disables Sentry without error.

## Supplying the DSN

**Local development** — add to `backend/.env`:
```dotenv
SENTRY_DSN=https://<key>@o<org>.ingest.sentry.io/<project>
```

**Production** — set `SENTRY_DSN` as an environment variable in your deployment platform. The app reads it at startup via `godotenv/autoload` (dev) or the process environment (production).

Leave `SENTRY_DSN` empty (or omit it) to run without Sentry. The app starts and serves normally in both cases.
