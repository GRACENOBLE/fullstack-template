package queue_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hibiken/asynq"

	"backend/internal/infrastructure/queue"
)

func TestHandleWelcomeEmail_ValidPayload(t *testing.T) {
	payload, err := json.Marshal(queue.WelcomeEmailPayload{UserID: "u1", Email: "test@example.com"})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	task := asynq.NewTask(queue.TypeWelcomeEmail, payload)
	if err := queue.HandleWelcomeEmail(context.Background(), task); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHandleWelcomeEmail_InvalidPayload(t *testing.T) {
	task := asynq.NewTask(queue.TypeWelcomeEmail, []byte("not-json"))
	if err := queue.HandleWelcomeEmail(context.Background(), task); err == nil {
		t.Error("expected error for invalid payload, got nil")
	}
}
