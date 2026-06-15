package domain

import "time"

// FCMToken is a Firebase Cloud Messaging device registration token stored per user.
type FCMToken struct {
	ID        string
	UserID    string
	Token     string
	Platform  string
	CreatedAt time.Time
}
