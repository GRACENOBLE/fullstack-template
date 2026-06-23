---
topic: streams
last_verified: 2026-06-23
sources:
  - internal/infrastructure/streams/events.go
  - internal/infrastructure/streams/producer.go
  - internal/infrastructure/streams/consumer.go
  - internal/usecase/streams.go
  - cmd/api/main.go
---

# Redis Streams

## Overview

Redis Streams provide an ordered, persistent event log for event-driven fan-out between
services. The same `REDIS_URL` env var used by the cache and queue layers is reused — no
additional env vars are required.

Streams complement Asynq rather than replace it:

| Concern | Mechanism |
|---|---|
| Discrete retriable background jobs | Asynq (`backend/docs/queue.md`) |
| Ordered event log / fan-out to multiple consumers | Redis Streams (`this file`) |

All Streams code lives in `internal/infrastructure/streams/`.

## StreamProducer interface

`internal/usecase/streams.go` defines the interface used by handlers and use cases:

```go
type StreamProducer interface {
    Publish(ctx context.Context, stream string, event any) error
    Close() error
}
```

`*streams.Producer` implements this interface. `app.StreamProducer` is closed during graceful
shutdown in `cmd/api/main.go` after the HTTP server and worker have stopped.

## Runtime wiring status

**Producer:** wired — `app.StreamProducer` is set when `REDIS_URL` is non-empty. It is
passed to `handlers.NewHandler` and is available for use in handlers and use cases.

**Consumer:** intentionally NOT started at runtime. The consumer infrastructure is fully
implemented but no consumer goroutine is registered in `cmd/api/main.go`. The commented-out
wiring block in `main.go` serves as the canonical example for adding a consumer:

```go
// streamCtx, streamCancel := context.WithCancel(context.Background())
// consumer, err := streams.NewConsumer(app.Config.RedisURL, streams.StreamUserCreated, "api", "api-1")
// if err != nil { ... }
// go func() {
//     _ = consumer.Run(streamCtx, func(ctx context.Context, data []byte) error {
//         var evt streams.UserCreatedEvent
//         if err := json.Unmarshal(data, &evt); err != nil { return err }
//         payload, _ := json.Marshal(queue.WelcomeEmailPayload{UserID: evt.UserID, Email: evt.Email, Name: evt.Name})
//         return app.Enqueuer.Enqueue(ctx, queue.TypeWelcomeEmail, payload)
//     })
//     _ = consumer.Close()
// }()
// // In shutdown: streamCancel()
```

## Stream names

Stream name constants and event payload structs are defined in
`internal/infrastructure/streams/events.go`:

```go
const (
    StreamUserCreated      = "stream:user.created"
    StreamNotificationSent = "stream:notification.sent"
)

type UserCreatedEvent struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    Name   string `json:"name"`
}
```

Naming convention: `stream:<domain>.<event>` — all lowercase, dot separator between
domain noun and past-tense verb.

## Publishing events

`internal/infrastructure/streams/producer.go`

```go
func NewProducer(redisURL string) (*Producer, error)
func (p *Producer) Publish(ctx context.Context, stream string, event any) error
func (p *Producer) Close() error
```

`Publish` marshals `event` to JSON and appends it to the named stream via `XADD` with
auto-generated IDs (`*`). The JSON is stored in a single `"data"` field:

```go
p.client.XAdd(ctx, &redis.XAddArgs{
    Stream: stream,
    Values: map[string]any{"data": string(payload)},
})
```

Call `producer.Close()` during application shutdown to release the Redis connection.

### Example — publishing from a use case

```go
producer, err := streams.NewProducer(redisURL)
// ...
err = producer.Publish(ctx, streams.StreamUserCreated, streams.UserCreatedEvent{
    UserID: user.ID,
    Email:  user.Email,
    Name:   user.Name,
})
```

## Consuming events

`internal/infrastructure/streams/consumer.go`

```go
func NewConsumer(redisURL, stream, group, consumer string) (*Consumer, error)
func (c *Consumer) Run(ctx context.Context, h Handler) error   // blocking; cancel ctx to stop
func (c *Consumer) Close() error

type Handler func(ctx context.Context, data []byte) error
```

`NewConsumer` takes four arguments: the Redis URL, the stream name constant, a consumer
group name, and a unique consumer name within the group.

### Run behaviour

`Run` blocks until `ctx` is cancelled. On each iteration it:

1. Calls `XGroupCreateMkStream` once at startup — creates the consumer group and the
   stream itself if either does not exist (`MKSTREAM`).
2. Issues `XReadGroup` with `">"` to fetch only new (undelivered) messages, blocking up
   to 2000 ms per call.
3. For each message, calls the `Handler` with the raw JSON bytes from the `"data"` field.
4. On handler success: acknowledges the message with `XACK`.
5. On handler error: logs via `slog.Error` and continues — the message is not
   acknowledged and will be redelivered on next startup (pending entries).
6. On `ctx` cancellation: returns `nil` cleanly.
7. On `redis.Nil` (timeout with no messages): continues the loop.
8. On other Redis errors: logs and continues.

### Idle-connection hardening

`NewConsumer` sets `MaxIdleConns=1` and `ConnMaxIdleTime=8s` on the Redis client options.
This proactively recycles idle connections before managed Redis providers (Upstash, Redis
Cloud) kill them at their own idle timeout (~10–60 s). Without this, long-blocking
`XReadGroup` calls surface as `i/o timeout` errors from the pool.

`Run` also handles `net.Error` timeouts explicitly: it logs a warning and backs off 2 s
before reconnecting, rather than exiting or spinning tightly.

### Example — registering a consumer goroutine in main.go

See the commented-out wiring block in `cmd/api/main.go` (reproduced in the Runtime wiring
status section above). The key pattern:

```go
streamCtx, streamCancel := context.WithCancel(context.Background())
consumer, _ := streams.NewConsumer(redisURL, streams.StreamUserCreated, "api", "api-1")
go func() {
    _ = consumer.Run(streamCtx, func(ctx context.Context, data []byte) error { ... })
    _ = consumer.Close()
}()
// In shutdown sequence: streamCancel()
```

## Adding a new event type

1. Add a `Stream<Domain><Event> = "stream:<domain>.<event>"` constant to
   `internal/infrastructure/streams/events.go`.
2. Add a `<Domain><Event>Event` payload struct with JSON tags in the same file.
3. Call `producer.Publish(ctx, streams.<StreamConst>, <Domain><Event>Event{...})` from
   the relevant domain action.
4. Register a `streams.NewConsumer` goroutine in `cmd/api/main.go` (or the consuming
   service's entry point), passing a `Handler` func that processes the raw JSON bytes.

## Testing

Integration tests use a Testcontainers Redis instance following the same `TestMain`
pattern as `internal/infrastructure/cache/redis/cache_test.go`.

Typical test flow:

```go
producer, _ := streams.NewProducer(redisURL)
producer.Publish(ctx, streams.StreamUserCreated, streams.UserCreatedEvent{
    UserID: "u1", Email: "a@b.com", Name: "Alice",
})
producer.Close()

// Verify the message was appended
msgs, _ := redisClient.XRange(ctx, streams.StreamUserCreated, "-", "+").Result()
// assert len(msgs) == 1 and msgs[0].Values["data"] contains expected JSON
```

Consumer handler logic is tested by invoking the `Handler` func directly with a
pre-marshalled `[]byte` payload — no Redis required for unit tests.
