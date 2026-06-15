package streams

// Stream name constants for Redis Streams event log.
const (
	StreamUserCreated      = "stream:user.created"
	StreamNotificationSent = "stream:notification.sent"
)

// UserCreatedEvent is published when a new user account is created.
type UserCreatedEvent struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}
