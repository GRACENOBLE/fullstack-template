---
topic: websocket
last_verified: 2026-06-23
sources:
  - internal/infrastructure/ws/message.go
  - internal/infrastructure/ws/hub.go
  - internal/infrastructure/ws/client.go
  - internal/infrastructure/ws/hub_test.go
  - internal/transport/handlers/ws_handler.go
  - internal/transport/handlers/routes.go
  - internal/server/server.go
  - cmd/api/main.go
---

# WebSocket

## Overview

Real-time bidirectional communication is provided via `github.com/gorilla/websocket`.
A `Hub` runs as a long-lived goroutine and fans out messages to all connected clients.
The `GET /ws` endpoint upgrades HTTP connections; a Firebase ID token is required as a
query parameter.

## Message types

`internal/infrastructure/ws/message.go` defines the wire format and inbound handler types:

```go
// Outbound (server → client) and inbound (client → server) wire format.
type Envelope struct {
    Type    string          `json:"type"`
    Payload json.RawMessage `json:"payload"`
}

// InboundMessage wraps an Envelope with the originating client ID.
type InboundMessage struct {
    ClientID string
    Msg      Envelope
}

// InboundHandler processes an inbound WebSocket message.
type InboundHandler func(ctx context.Context, msg InboundMessage) error
```

`Type` is a dot-separated event name (e.g. `"job.completed"`). `Payload` is arbitrary
JSON whose shape is determined by `Type`.

The same `Envelope` format is used for both directions — clients send JSON frames in this
shape, and the server broadcasts frames in this shape.

## Hub

`internal/infrastructure/ws/hub.go`

```go
func NewHub() *Hub
func (h *Hub) Run(ctx context.Context)                            // blocking; cancel ctx to stop
func (h *Hub) Publish(msgType string, payload any) error          // server → all clients
func (h *Hub) OnMessage(msgType string, handler InboundHandler)   // register inbound handler
```

`Run` must be called in its own goroutine and runs until `ctx` is cancelled.
`Publish` marshals `payload` into an `Envelope` and queues it for broadcast — safe to call
from any goroutine.
`OnMessage` registers a handler for a specific inbound message type. Safe to call before
or after `Run` starts (protected by `sync.RWMutex`).

### Goroutine model

The Hub is now fully bidirectional. Each WebSocket connection spawns `ReadPump` and `WritePump`.
`ReadPump` parses each inbound frame as an `Envelope` and forwards it to `hub.inbound`.
`Run` dispatches inbound messages to registered handlers via a 64-slot semaphore-bounded
goroutine pool so slow handlers do not stall the event loop.

```text
 WebSocket conn ──► client.ReadPump goroutine ──► hub.inbound chan
                                                        │
                                                  hub.Run goroutine ──► InboundHandler goroutine (≤64 concurrent)

 caller goroutine
       │  hub.Publish(...)
       ▼
  hub.broadcast chan
       │
   hub.Run goroutine ──► client.Send chan ──► client.WritePump goroutine ──► WebSocket conn
```

Slow clients are dropped: if `client.Send` is full, the Hub closes the channel and removes
the client without blocking the broadcast loop.

When `hub.inbound` is full, the message is dropped and a warning is logged (non-blocking).
When the semaphore is full (64 in-flight handlers), additional inbound messages are dropped.

### Lifecycle context

The `ctx` passed to `hub.Run` is forwarded to each `InboundHandler` goroutine. This means
handlers can respect cancellation and use it for downstream calls (e.g. `db.QueryContext`).

## Client

`internal/infrastructure/ws/client.go`

```go
type Client struct {
    hub  *Hub
    conn *websocket.Conn
    Send chan []byte   // exported for testing
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client   // registers with hub
func (c *Client) ReadPump()                              // must run in goroutine
func (c *Client) WritePump()                             // must run in goroutine
```

Ping/pong keepalive: `pingPeriod = 54s`, `pongWait = 60s`, `writeWait = 10s`.

## Route — GET /ws

```text
GET /ws?token=<firebase-id-token>
```

Defined in `internal/transport/handlers/ws_handler.go`. Registered in `RegisterRoutes`
outside the `/api/v1` auth group — auth is handled inline because WebSocket clients
cannot set `Authorization` headers.

**Auth flow:**
1. If `h.verifier != nil` (staging / production): reads `?token=` query param.
   Returns `401` when missing or when `VerifyIDToken` fails.
2. If `h.verifier == nil` (development): skips auth — connects immediately.

After successful auth, the connection is upgraded and `ReadPump` / `WritePump` are
started in separate goroutines. `ReadPump` now parses each inbound frame as an `Envelope`
and routes it to the Hub — malformed frames are logged and discarded (the connection stays open).

## Wiring in server.go and main.go

`server.go` accepts `*ws.Hub` as a second argument:

```go
func NewServer(app *bootstrap.App, hub *ws.Hub) (*http.Server, error)
```

`cmd/api/main.go` creates the Hub, starts `Run` with a child context, and cancels
it after the HTTP server shuts down (so all in-flight connections close first):

```go
hubCtx, hubCancel := context.WithCancel(context.Background())
hub := ws.NewHub()
go hub.Run(hubCtx)

srv, err := server.NewServer(app, hub)
// ...
<-done
hubCancel()   // stop hub after server drains connections
```

The canonical `Handler` struct definition and `NewHandler` signature (including all fields beyond `hub` and `verifier`) are documented in `backend/docs/routing.md`.

## Publishing events from workers

Call `hub.Publish` from any goroutine (Asynq workers, Redis Streams consumers, etc.):

```go
hub.Publish("job.completed", map[string]any{
    "jobId":  id,
    "status": "done",
})
```

## Handling inbound messages

Register an `InboundHandler` before or after calling `hub.Run`:

```go
hub.OnMessage("ping", func(ctx context.Context, msg ws.InboundMessage) error {
    // msg.ClientID is the remote address of the originating client
    // msg.Msg.Payload contains the raw JSON payload from the client
    return hub.Publish("pong", map[string]any{"client": msg.ClientID})
})
```

Handlers must be non-blocking or return quickly; the 64-slot semaphore limits concurrency.

## Testing

Unit tests in `internal/infrastructure/ws/hub_test.go` cover:
- `TestHub_RegisterAndBroadcast` — single client receives broadcast
- `TestHub_UnregisterRemovesClient` — Send channel is closed on unregister
- `TestHub_ContextCancelClosesSendChannels` — all channels closed on ctx cancel
- `TestHub_ConcurrentClientsAndBroadcast` — 10 concurrent clients all receive
- `TestHub_OnMessage` — `InboundHandler` is invoked with correct `InboundMessage`
- `TestHub_Publish` — JSON marshalling and delivery

Tests inject `*Client` with a nil `conn` and a buffered `Send` channel; the Hub only
touches `client.Send`, not the connection, so real WebSocket connections are not needed.
