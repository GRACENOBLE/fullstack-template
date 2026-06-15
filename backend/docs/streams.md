---
topic: streams
last_verified: 2026-06-15
sources:
  - internal/infrastructure/streams/events.go
  - internal/infrastructure/streams/producer.go
  - internal/infrastructure/streams/consumer.go
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

### Example — registering a consumer goroutine in main.go

```go
consumer, err := streams.NewConsumer(app.Config.RedisURL,
    streams.StreamUserCreated, "api-group", "api-consumer-1")
if err != nil {
    // handle
}
go func() {
    if err := consumer.Run(ctx, handleUserCreated); err != nil {
        slog.Error("streams: consumer error", "err", err)
    }
}()
```

Cancel `ctx` (the same child context used for the worker) to stop the consumer during
shutdown, then call `consumer.Close()` to release the connection.

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
    UserID: "u1", Email: "a@b.com",
})
producer.Close()

// Verify the message was appended
msgs, _ := redisClient.XRange(ctx, streams.StreamUserCreated, "-", "+").Result()
// assert len(msgs) == 1 and msgs[0].Values["data"] contains expected JSON
```

Consumer handler logic is tested by invoking the `Handler` func directly with a
pre-marshalled `[]byte` payload — no Redis required for unit tests.
