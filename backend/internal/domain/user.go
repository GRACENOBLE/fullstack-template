package domain

import "time"

// User represents an authenticated user's profile record.
type User struct {
	ID          int64     `json:"id"`
	FirebaseUID string    `json:"firebase_uid"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhotoURL    string    `json:"photo_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
