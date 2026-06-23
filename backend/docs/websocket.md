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

## Message envelope

All messages use a typed JSON envelope defined in `internal/infrastructure/ws/message.go`:

```go
type Envelope struct {
    Type    string          `json:"type"`
    Payload json.RawMessage `json:"payload"`
}
```

`Type` is a dot-separated event name (e.g. `"job.completed"`). `Payload` is arbitrary
JSON whose shape is determined by `Type`.

## Hub

`internal/infrastructure/ws/hub.go`

```go
type Hub struct {
    clients    map[*Client]struct{}
    broadcast  chan []byte    // buffered, capacity 256
    Register   chan *Client
    Unregister chan *Client
}

func NewHub() *Hub
func (h *Hub) Run(ctx context.Context)          // blocking; cancel ctx to stop
func (h *Hub) Publish(msgType string, payload any) error
```

`Run` must be called in its own goroutine and runs until `ctx` is cancelled.
`Publish` marshals `payload` into an `Envelope` and queues it for broadcast —
safe to call from any goroutine (e.g. an Asynq worker in future #18).

### Goroutine model

Each WebSocket connection spawns two goroutines: `ReadPump` and `WritePump` (on `Client`).
The Hub serialises all mutations (register / unregister / broadcast) through a `select` loop
so no locking is needed on its internal `clients` map.

```text
 caller goroutine
       │  hub.Publish(...)
       ▼
  hub.broadcast chan
       │
   hub.Run goroutine ──► client.Send chan ──► client.WritePump goroutine ──► WebSocket conn
                                              client.ReadPump goroutine  ──► (discards incoming / handles pings)
```

Slow clients are dropped: if `client.Send` is full, the Hub closes the channel and
removes the client without blocking the broadcast loop.

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
started in separate goroutines.

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

## Publishing events from workers (future #18)

Call `hub.Publish` from any goroutine:

```go
hub.Publish("job.completed", map[string]any{
    "jobId":  id,
    "status": "done",
})
```

When #18 (Asynq + Redis Streams) lands, the Asynq task handlers and Redis Streams
consumers will call this method to push domain events to connected clients.

## Testing

Unit tests in `internal/infrastructure/ws/hub_test.go` cover:
- `TestHub_RegisterAndBroadcast` — single client receives broadcast
- `TestHub_UnregisterRemovesClient` — Send channel is closed on unregister
- `TestHub_ContextCancelClosesSendChannels` — all channels closed on ctx cancel
- `TestHub_ConcurrentClientsAndBroadcast` — 10 concurrent clients all receive
- `TestHub_Publish` — JSON marshalling and delivery

Tests inject `*Client` with a nil `conn` and a buffered `Send` channel; the Hub only
touches `client.Send`, not the connection, so real WebSocket connections are not needed.
