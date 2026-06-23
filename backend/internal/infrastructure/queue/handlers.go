package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"

	"backend/internal/usecase"
)

// NewHandleWelcomeEmail returns an asynq handler that sends a welcome email via sender.
// If sender is nil, the task is acknowledged without sending (graceful degradation).
func NewHandleWelcomeEmail(sender usecase.EmailSender) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var p WelcomeEmailPayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("welcome email: unmarshal payload: %w", err)
		}
		if sender == nil {
			return nil
		}
		name := p.Name
		if name == "" {
			name = p.Email
		}
		return sender.SendWelcomeEmail(ctx, p.Email, name)
	}
}
