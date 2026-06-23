package queue

// Task type constants used by both enqueuers and handlers.
const (
	TypeWelcomeEmail = "email:welcome"
)

// WelcomeEmailPayload is the JSON payload for TypeWelcomeEmail tasks.
type WelcomeEmailPayload struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}
