package usecase

import "context"

// FirebaseToken holds the verified claims extracted from a Firebase ID token.
type FirebaseToken struct {
	UID      string         `json:"uid"`
	Email    string         `json:"email"`
	Name     string         `json:"name"`
	PhotoURL string         `json:"photoUrl"`
	Claims   map[string]any `json:"claims"`
}

// FirebaseTokenVerifier can verify a raw Firebase ID token string.
type FirebaseTokenVerifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (*FirebaseToken, error)
}

// FirebaseAdminClient extends FirebaseTokenVerifier with Firebase user-management operations.
type FirebaseAdminClient interface {
	FirebaseTokenVerifier
	GetUserByEmail(ctx context.Context, email string) (string, error)
	UpdateUserPassword(ctx context.Context, uid, newPassword string) error
}
