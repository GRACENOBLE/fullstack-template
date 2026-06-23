package ws

import (
	"context"
	"sync"
	"testing"
	"time"
)

func newTestClient(hub *Hub) *Client {
	return &Client{hub: hub, Send: make(chan []byte, 16)}
}

func TestHub_RegisterAndBroadcast(t *testing.T) {
	hub := NewHub()
	go hub.Run(t.Context())

	c := newTestClient(hub)
	hub.Register <- c
	time.Sleep(5 * time.Millisecond)

	msg := []byte(`{"type":"ping","payload":null}`)
	hub.broadcast <- msg

	select {
	case got := <-c.Send:
		if string(got) != string(msg) {
			t.Errorf("got %q, want %q", got, msg)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for broadcast")
	}
}

func TestHub_UnregisterRemovesClient(t *testing.T) {
	hub := NewHub()
	go hub.Run(t.Context())

	c := newTestClient(hub)
	hub.Register <- c
	time.Sleep(5 * time.Millisecond)

	hub.Unregister <- c
	time.Sleep(5 * time.Millisecond)

	select {
	case _, ok := <-c.Send:
		if ok {
			t.Error("expected Send channel to be closed")
		}
	default:
		t.Error("expected Send channel to be closed but it was still open")
	}
}

func TestHub_ContextCancelClosesSendChannels(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	hub := NewHub()
	go hub.Run(ctx)

	c := newTestClient(hub)
	hub.Register <- c
	time.Sleep(5 * time.Millisecond)

	cancel()
	time.Sleep(10 * time.Millisecond)

	select {
	case _, ok := <-c.Send:
		if ok {
			t.Error("expected Send channel to be closed after context cancel")
		}
	default:
		t.Error("expected Send channel to be closed but it was still open")
	}
}

func TestHub_ConcurrentClientsAndBroadcast(t *testing.T) {
	hub := NewHub()
	go hub.Run(t.Context())

	const n = 10
	clients := make([]*Client, n)
	for i := range clients {
		clients[i] = newTestClient(hub)
		hub.Register <- clients[i]
	}
	time.Sleep(10 * time.Millisecond)

	msg := []byte(`{"type":"hello","payload":null}`)

	var wg sync.WaitGroup
	for _, c := range clients {
		wg.Add(1)
		go func(cl *Client) {
			defer wg.Done()
			select {
			case got := <-cl.Send:
				if string(got) != string(msg) {
					t.Errorf("client got %q, want %q", got, msg)
				}
			case <-time.After(time.Second):
				t.Error("client timed out waiting for broadcast")
			}
		}(c)
	}

	hub.broadcast <- msg
	wg.Wait()
}

func TestHub_OnMessage(t *testing.T) {
	hub := NewHub()
	go hub.Run(t.Context())

	called := make(chan InboundMessage, 1)
	hub.OnMessage("ping", func(_ context.Context, msg InboundMessage) error {
		called <- msg
		return nil
	})

	im := InboundMessage{ClientID: "client-1", Msg: Envelope{Type: "ping"}}
	select {
	case hub.inbound <- im:
	case <-time.After(time.Second):
		t.Fatal("timed out sending to hub.inbound")
	}

	select {
	case got := <-called:
		if got.ClientID != "client-1" {
			t.Errorf("ClientID: got %q, want %q", got.ClientID, "client-1")
		}
		if got.Msg.Type != "ping" {
			t.Errorf("Msg.Type: got %q, want %q", got.Msg.Type, "ping")
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for OnMessage handler to be called")
	}
}

func TestHub_Publish(t *testing.T) {
	hub := NewHub()
	go hub.Run(t.Context())

	c := newTestClient(hub)
	hub.Register <- c
	time.Sleep(5 * time.Millisecond)

	if err := hub.Publish("job.completed", map[string]string{"id": "123"}); err != nil {
		t.Fatalf("Publish() error: %v", err)
	}

	select {
	case got := <-c.Send:
		if len(got) == 0 {
			t.Error("expected non-empty message from Publish")
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for published message")
	}
}
