package usecase

import "context"

// EmailSender sends transactional email messages.
type EmailSender interface {
	SendWelcomeEmail(ctx context.Context, toEmail, toName string) error
}
