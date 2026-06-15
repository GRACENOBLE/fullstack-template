package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/hibiken/asynq"
)

// HandleWelcomeEmail logs the welcome email payload.
// Real delivery via Mailjet is wired in issue #20.
func HandleWelcomeEmail(_ context.Context, t *asynq.Task) error {
	var p WelcomeEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("welcome email: unmarshal payload: %w", err)
	}
	slog.Info("queue: welcome email task received", "user_id", p.UserID, "email", p.Email)
	return nil
}
