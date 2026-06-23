package ws

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

// Client couples a single WebSocket connection to the Hub.
// Send is exported so hub_test.go can inject a client without a real conn.
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	Send chan []byte
}

// NewClient allocates a Client and registers it with the hub.
// Call WritePump and ReadPump in separate goroutines after this returns.
func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	c := &Client{
		hub:  hub,
		conn: conn,
		Send: make(chan []byte, 256),
	}
	hub.Register <- c
	return c
}

// ReadPump reads from the WebSocket and handles connection liveness.
// It unregisters the client and closes the connection when done.
func (c *Client) ReadPump() {
	defer func() {
		// Non-blocking: if Hub.Run has already exited (ctx cancelled) there is no
		// receiver on Unregister. The Hub already closed c.Send in that case.
		select {
		case c.hub.Unregister <- c:
		default:
		}
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		return
	}
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("websocket read error", "error", err)
			}
			break
		}
		var env Envelope
		if jsonErr := json.Unmarshal(raw, &env); jsonErr != nil {
			slog.Warn("ws: malformed inbound frame, discarding", "client", c.conn.RemoteAddr(), "err", jsonErr)
			continue
		}
		select {
		case c.hub.inbound <- InboundMessage{ClientID: c.conn.RemoteAddr().String(), Msg: env}:
		default:
			slog.Warn("ws: inbound channel full, message dropped")
		}
	}
}

// WritePump writes queued messages to the WebSocket and sends periodic pings.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.Send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return
			}
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			if _, err := w.Write(msg); err != nil {
				return
			}
			// Drain any queued messages into the same WebSocket frame.
			for i := len(c.Send); i > 0; i-- {
				if _, err := w.Write([]byte{'\n'}); err != nil {
					break
				}
				if _, err := w.Write(<-c.Send); err != nil {
					break
				}
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
