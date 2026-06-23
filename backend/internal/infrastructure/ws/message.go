package ws

import (
	"context"
	"encoding/json"
)

// Envelope is the typed wire format for all WebSocket messages.
type Envelope struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// InboundMessage is a message received from a WebSocket client.
type InboundMessage struct {
	ClientID string
	Msg      Envelope
}

// InboundHandler processes an inbound WebSocket message.
type InboundHandler func(ctx context.Context, msg InboundMessage) error
