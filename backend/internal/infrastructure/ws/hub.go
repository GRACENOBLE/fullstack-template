package ws

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"sync"
)

const inboundWorkers = 64 // max concurrent inbound handler goroutines

// Hub maintains the set of active WebSocket clients and broadcasts messages to them.
// All mutations to the clients map are serialised through the Run goroutine.
// msgHandlers is protected by mu and may be updated concurrently.
type Hub struct {
	clients     map[*Client]struct{}
	broadcast   chan []byte
	Register    chan *Client
	Unregister  chan *Client
	inbound     chan InboundMessage
	msgHandlers map[string]InboundHandler
	mu          sync.RWMutex
	sem         chan struct{} // bounds concurrent inbound handler goroutines
}

// NewHub allocates a Hub with buffered channels.
func NewHub() *Hub {
	return &Hub{
		clients:     make(map[*Client]struct{}),
		broadcast:   make(chan []byte, 256),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		inbound:     make(chan InboundMessage, 256),
		msgHandlers: make(map[string]InboundHandler),
		sem:         make(chan struct{}, inboundWorkers),
	}
}

// OnMessage registers a handler for the given message type.
// Safe to call before or after Run starts.
func (h *Hub) OnMessage(msgType string, handler InboundHandler) {
	h.mu.Lock()
	h.msgHandlers[msgType] = handler
	h.mu.Unlock()
}

// Run processes register, unregister, broadcast, and inbound events until ctx is cancelled.
// Call this in its own goroutine.
func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			for c := range h.clients {
				close(c.Send)
			}
			return
		case c := <-h.Register:
			h.clients[c] = struct{}{}
		case c := <-h.Unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.Send)
			}
		case msg := <-h.broadcast:
			for c := range h.clients {
				select {
				case c.Send <- msg:
				default:
					// Slow client: drop and disconnect.
					close(c.Send)
					delete(h.clients, c)
				}
			}
		case im := <-h.inbound:
			h.mu.RLock()
			fn := h.msgHandlers[im.Msg.Type]
			h.mu.RUnlock()
			if fn == nil {
				continue
			}
			select {
			case h.sem <- struct{}{}:
				go func(msg InboundMessage) {
					defer func() { <-h.sem }()
					if err := fn(ctx, msg); err != nil {
						slog.Error("ws: inbound handler error", "type", msg.Msg.Type, "err", err)
					}
				}(im)
			default:
				slog.Warn("ws: inbound worker pool full, message dropped", "type", im.Msg.Type)
			}
		}
	}
}

// Publish marshals msgType + payload into an Envelope and queues it for broadcast.
// Safe to call from any goroutine (e.g. an Asynq worker or Redis Streams consumer).
func (h *Hub) Publish(msgType string, payload any) error {
	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	env, err := json.Marshal(Envelope{Type: msgType, Payload: raw})
	if err != nil {
		return err
	}
	select {
	case h.broadcast <- env:
		return nil
	default:
		return errors.New("ws: broadcast channel full, message dropped")
	}
}
