---
topic: queue
last_verified: 2026-06-15
sources:
  - internal/usecase/enqueuer.go
  - internal/infrastructure/queue/tasks.go
  - internal/infrastructure/queue/client.go
  - internal/infrastructure/queue/worker.go
  - internal/infrastructure/queue/handlers.go
  - internal/transport/handlers/routes.go
  - cmd/api/main.go
---

# Queue (Asynq)

## Overview

Background jobs are processed by [Asynq](https://github.com/hibiken/asynq), which uses
Redis as its broker. The same `REDIS_URL` env var used by the cache layer is reused — no
additional env vars are required. When `REDIS_URL` is not set, `app.Enqueuer` is `nil`
and no worker is started; the rest of the application continues normally.

Asynq handles discrete, retriable background jobs (welcome emails, notifications,
webhooks). For ordered event fan-out between services see `backend/docs/streams.md`.

## Task definitions

Task type constants and payload structs are co-located in
`internal/infrastructure/queue/tasks.go`:

```go
const (
    TypeWelcomeEmail = "email:welcome"
)

type WelcomeEmailPayload struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
}
```

Both the enqueuer (caller side) and the handler (worker side) import these constants so
the string is never duplicated.

## Enqueuer interface

`internal/usecase/enqueuer.go` defines the port that use cases and handlers depend on:

```go
type Enqueuer interface {
    Enqueue(ctx context.Context, taskType string, payload []byte) error
    Close() error
}
```

The concrete implementation is created with `queue.NewClient`:

```go
// internal/infrastructure/queue/client.go
func NewClient(redisURL string) (usecase.Enqueuer, error)
```

`NewClient` parses `redisURL` using `go-redis/v9`'s `ParseURL`, constructs an
`asynq.Client`, and returns it wrapped as a `usecase.Enqueuer`.

### Enqueueing from a use case or handler

```go
payload, err := json.Marshal(queue.WelcomeEmailPayload{
    UserID: userID,
    Email:  email,
})
if err != nil {
    return err
}
if app.Enqueuer != nil {
    if err := app.Enqueuer.Enqueue(ctx, queue.TypeWelcomeEmail, payload); err != nil {
        return err
    }
}
```

Guard with `!= nil` because `app.Enqueuer` is nil when `REDIS_URL` is unset.

### Bootstrap wiring

`app.Enqueuer` is populated by `bootstrap.Run` when `REDIS_URL` is set (see
`internal/bootstrap/bootstrap.go`). After the HTTP server shuts down, `main.go` calls
`app.Enqueuer.Close()` to release the connection.

## Worker setup

`internal/infrastructure/queue/worker.go`

```go
func NewWorker(redisURL string) (*Worker, error)
func (w *Worker) Register(taskType string, h asynq.Handler)
func (w *Worker) Run(ctx context.Context) error   // blocking; cancel ctx to stop
```

`NewWorker` creates an `asynq.Server` with `Concurrency: 10`. Failed tasks are logged via
`slog.Error` through the server's `ErrorHandler`.

`Run` starts the server in an inner goroutine and blocks until either `ctx` is cancelled
(clean shutdown via `w.server.Shutdown()`) or the server returns an error.

### Goroutine wiring in cmd/api/main.go

The worker is started only when `REDIS_URL` is set, using a child context so it can be
cancelled independently of the hub:

```go
var workerCancel context.CancelFunc
if app.Config.RedisURL != "" {
    workerCtx, wCancel := context.WithCancel(context.Background())
    workerCancel = wCancel
    worker, err := queue.NewWorker(app.Config.RedisURL)
    // ...
    worker.Register(queue.TypeWelcomeEmail, asynq.HandlerFunc(queue.HandleWelcomeEmail))
    go func() {
        if err := worker.Run(workerCtx); err != nil {
            slog.Error("queue: worker error", "err", err)
        }
    }()
}
```

Shutdown order after `<-done`:

```go
if workerCancel != nil {
    workerCancel() // stop worker before hub (in-flight jobs drain first)
}
hubCancel()
```

The worker is cancelled before the hub so that any task handler that calls `hub.Publish`
can still reach the hub during drain.

## Task handlers

Handler functions have the signature `func(context.Context, *asynq.Task) error` and live
in `internal/infrastructure/queue/handlers.go`:

```go
func HandleWelcomeEmail(_ context.Context, t *asynq.Task) error {
    var p WelcomeEmailPayload
    if err := json.Unmarshal(t.Payload(), &p); err != nil {
        return fmt.Errorf("welcome email: unmarshal payload: %w", err)
    }
    slog.Info("queue: welcome email task received", "user_id", p.UserID, "email", p.Email)
    return nil
}
```

Returning a non-nil error causes Asynq to retry the task (up to its configured retry limit).

## Adding a new task

1. Add a `Type<Name> = "<category>:<action>"` constant and `<Name>Payload` struct to
   `internal/infrastructure/queue/tasks.go`.
2. Write a `Handle<Name>(ctx context.Context, t *asynq.Task) error` function in
   `internal/infrastructure/queue/handlers.go`.
3. Register the handler in `cmd/api/main.go`:
   ```go
   worker.Register(queue.Type<Name>, asynq.HandlerFunc(queue.Handle<Name>))
   ```
4. Enqueue from the relevant use case or handler using `app.Enqueuer.Enqueue(ctx, queue.Type<Name>, payload)`.

## Asynqmon UI

When running in `gin.DebugMode` and `REDIS_URL` is set, the Asynqmon job-monitoring UI is
available at:

```
http://localhost:8080/admin/queues
```

Routes are registered in `RegisterRoutes` only under these conditions:

```go
if gin.Mode() == gin.DebugMode && h.queueUI != nil {
    r.GET("/admin/queues", gin.WrapH(h.queueUI))
    r.Any("/admin/queues/*path", gin.WrapH(h.queueUI))
}
```

The UI is not mounted in staging or production (`gin.ReleaseMode`).

## Testing

**Unit tests** call handler functions directly with `asynq.NewTask` — no Redis required:

```go
payload, _ := json.Marshal(queue.WelcomeEmailPayload{UserID: "u1", Email: "a@b.com"})
task := asynq.NewTask(queue.TypeWelcomeEmail, payload)
err := queue.HandleWelcomeEmail(context.Background(), task)
// assert err == nil
```

**Integration tests** use Testcontainers Redis (same `TestMain` pattern as
`internal/infrastructure/cache/redis/cache_test.go`). The `TestMain` function skips
integration tests gracefully when Docker is unavailable so that unit tests still run.

Enqueuer integration tests construct a real `queue.NewClient`, enqueue a task, and use
`asynq.NewInspector` (via `queue.NewInspector`) to verify the task appears in the queue.
