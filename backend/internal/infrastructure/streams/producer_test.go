package streams

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	goredis "github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testRedisURL string

func mustStartRedisContainerForStreams() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
	ctx := context.Background()
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "redis:7-alpine",
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor:   wait.ForLog("Ready to accept connections"),
		},
		Started: true,
	})
	if err != nil {
		return nil, fmt.Errorf("start redis container: %w", err)
	}
	host, err := container.Host(ctx)
	if err != nil {
		return container.Terminate, fmt.Errorf("container host: %w", err)
	}
	port, err := container.MappedPort(ctx, "6379/tcp")
	if err != nil {
		return container.Terminate, fmt.Errorf("mapped port: %w", err)
	}
	testRedisURL = fmt.Sprintf("redis://%s:%s", host, port.Port())
	return container.Terminate, nil
}

func TestMain(m *testing.M) {
	teardown, err := mustStartRedisContainerForStreams()
	if err != nil {
		log.Fatalf("could not start redis container: %v", err)
	}
	m.Run()
	if teardown != nil {
		if err := teardown(context.Background()); err != nil {
			log.Fatalf("could not teardown redis container: %v", err)
		}
	}
}

func TestProducer_Publish(t *testing.T) {
	ctx := context.Background()
	p, err := NewProducer(testRedisURL)
	if err != nil {
		t.Fatalf("NewProducer: %v", err)
	}
	defer p.Close()

	event := UserCreatedEvent{UserID: "u1", Email: "test@example.com"}
	if err := p.Publish(ctx, StreamUserCreated, event); err != nil {
		t.Fatalf("Publish: %v", err)
	}

	// Read the message back directly to verify it was appended.
	opts, _ := goredis.ParseURL(testRedisURL)
	rc := goredis.NewClient(opts)
	defer rc.Close()

	entries, err := rc.XRange(ctx, StreamUserCreated, "-", "+").Result()
	if err != nil {
		t.Fatalf("XRange: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 stream entry, got %d", len(entries))
	}

	data, ok := entries[0].Values["data"].(string)
	if !ok {
		t.Fatal("stream entry missing 'data' field")
	}

	var got UserCreatedEvent
	if err := json.Unmarshal([]byte(data), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.UserID != event.UserID || got.Email != event.Email {
		t.Errorf("got %+v, want %+v", got, event)
	}
}
