package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"backend/internal/usecase"
)

var testEnqueuer usecase.Enqueuer

func mustStartRedisContainerForQueue() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
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
	redisURL := fmt.Sprintf("redis://%s:%s", host, port.Port())
	eq, err := NewClient(redisURL)
	if err != nil {
		return container.Terminate, fmt.Errorf("new queue client: %w", err)
	}
	testEnqueuer = eq
	return container.Terminate, nil
}

func TestMain(m *testing.M) {
	teardown, err := mustStartRedisContainerForQueue()
	if err != nil {
		// Docker unavailable — skip integration tests but let unit tests in this
		// package run (handlers_test.go has no container dependency).
		log.Printf("queue integration tests skipped: %v", err)
	}
	code := m.Run()
	if teardown != nil {
		if err := teardown(context.Background()); err != nil {
			log.Printf("teardown redis container: %v", err)
		}
	}
	os.Exit(code)
}

func TestEnqueue_WelcomeEmail(t *testing.T) {
	if testEnqueuer == nil {
		t.Skip("requires Docker")
	}
	payload, _ := json.Marshal(WelcomeEmailPayload{UserID: "u1", Email: "test@example.com"})
	err := testEnqueuer.Enqueue(context.Background(), TypeWelcomeEmail, payload)
	if err != nil {
		t.Errorf("Enqueue() returned unexpected error: %v", err)
	}
}
