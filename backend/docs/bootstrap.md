---
topic: bootstrap
last_verified: 2026-06-23
sources:
  - internal/bootstrap/bootstrap.go
  - internal/server/server.go
  - cmd/api/main.go
---

# Bootstrap

## Purpose
`internal/bootstrap` owns the startup lifecycle: load config, validate it, initialise shared dependencies, and probe services for readiness. It runs before the HTTP server starts and aborts the process cleanly if anything is wrong.

## App struct
```go
type App struct {
    DB          *sql.DB
    Cache       usecase.CacheService        // nil when REDIS_URL is not set
    Enqueuer    usecase.Enqueuer            // nil when REDIS_URL is not set
    Firebase    usecase.FirebaseAdminClient // nil when FIREBASE_PROJECT_ID is not set
    FCMSender   usecase.NotificationSender  // nil when FIREBASE_PROJECT_ID is not set
    EmailSender usecase.EmailSender         // nil when MAILJET_API_KEY/SECRET_KEY are not set
    Config      Config
    Log         *slog.Logger
}
```
`App` is constructed once by `Run` and passed to `server.NewServer`. Nothing re-initialises dependencies after this point. Optional fields (`Cache`, `Enqueuer`, `Firebase`, `FCMSender`, `EmailSender`) are nil when their corresponding env vars are absent.

## Config struct
```go
type Config struct {
    Port                       int
    Env                        string
    DB                         postgres.DBConfig
    RedisURL                   string
    RateLimitRPS               float64
    RateLimitBurst             int
    FirebaseProjectID          string
    FirebaseServiceAccountJSON string
    SentryDSN                  string
}
```
`loadConfig()` reads all values from environment variables. `PORT` defaults to `8080`; `BLUEPRINT_DB_SCHEMA` defaults to `public`; `BLUEPRINT_DB_SSLMODE` defaults to `disable`. `RateLimitBurst` is derived as `int(RPS)*5` when omitted and RPS is set. Optional fields (`RedisURL`, `FirebaseProjectID`, `FirebaseServiceAccountJSON`) default to empty string — their respective services are skipped when empty.

## Run sequence
`bootstrap.Run(ctx)` executes these steps in order:

1. Build `Config` from env vars via `loadConfig()`
2. Validate required fields via `validateConfig()` — fast, no I/O
3. Open `*sql.DB` via `postgres.NewPostgresDB(cfg.DB)`
4. Probe Postgres with `probeWithRetry` under a 60-second total timeout
5. Init Redis cache via `redis.New(cfg.RedisURL)` and probe it — skipped when `REDIS_URL` is empty
6. Init Asynq enqueuer via `queue.NewClient(cfg.RedisURL)` — skipped when `REDIS_URL` is empty
7. Init Firebase app via `firebase.NewApp(ctx, ...)`, then init Auth client (`firebase.NewAuthClient`) and FCM messaging client (`firebase.NewMessagingClient`) from the same app instance — all skipped when `FIREBASE_PROJECT_ID` is empty
8. Init Mailjet email sender via `email.NewMailjetSender(...)` — skipped when `MAILJET_API_KEY` or `MAILJET_SECRET_KEY` is empty; startup fails if only a partial Mailjet config is provided
9. Return `*App` on success; return a non-nil error on any failure

```go
ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
app, err := bootstrap.Run(ctx)
stop()
if err != nil {
    fmt.Fprintf(os.Stderr, "startup failed: %v\n", err)
    os.Exit(1)
}
```
The signal-aware context means Ctrl-C during startup cancels probes immediately rather than waiting for timeouts to expire.

## Config validation
`validateConfig` checks that all five required DB env vars are non-empty before attempting a connection:

```go
requireNonEmpty("BLUEPRINT_DB_HOST", cfg.DB.Host)
requireNonEmpty("BLUEPRINT_DB_PORT", cfg.DB.Port)
requireNonEmpty("BLUEPRINT_DB_DATABASE", cfg.DB.Database)
requireNonEmpty("BLUEPRINT_DB_USERNAME", cfg.DB.Username)
requireNonEmpty("BLUEPRINT_DB_PASSWORD", cfg.DB.Password)
```

Failures are collected and returned as `*ConfigError` — all issues reported at once, not just the first.

## Service probing
The `Pinger` interface is satisfied by `*sql.DB` natively:

```go
type Pinger interface {
    PingContext(ctx context.Context) error
}
```

`probeWithRetry` attempts up to 5 pings. Between failures it sleeps for a full-jitter exponential backoff: random duration in `[0, min(16s, 500ms × 2^attempt)]`. Each attempt has a 15-second deadline (sized to accommodate Neon cold starts). The total probe budget is 60 seconds.

Log output during probing:
```
bootstrap: probing service  service=postgres  attempt=1  max_attempts=5
bootstrap: service not ready  service=postgres  attempt=1  error=...
bootstrap: waiting before retry  service=postgres  attempt=2  delay=347ms
bootstrap: service ready  service=postgres  attempts=2
```

## Adding a new dependency
1. Add an initialisation function `initFoo(cfg Config, log *slog.Logger) *FooClient` — return nil when the dependency is optional and not configured.
2. Add the field to `App`.
3. If the dependency supports `PingContext`, add it to the probes slice in a `probeAll`-style helper; otherwise just initialise and nil-check.
4. Pass the field through in `server.NewServer(app)`.
