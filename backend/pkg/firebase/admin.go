package firebase

import (
	"context"
	"fmt"

	firebasesdk "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"

	"backend/internal/usecase"
)

// authClientAdapter wraps the Firebase Admin auth.Client and satisfies the
// usecase.FirebaseAdminClient interface without leaking the SDK into the application layer.
type authClientAdapter struct {
	client *auth.Client
}

// NewAuthClient initialises the Firebase Admin SDK and returns a usecase.FirebaseAdminClient
// ready for injection.
//
// credentialsJSON is the raw content of a service account JSON key
// (FIREBASE_SERVICE_ACCOUNT_JSON env var). When empty the SDK falls back to
// Application Default Credentials (ADC) — appropriate for GCP-hosted deployments.
func NewAuthClient(ctx context.Context, projectID, credentialsJSON string) (usecase.FirebaseAdminClient, error) {
	var opts []option.ClientOption
	if credentialsJSON != "" {
		opts = append(opts, option.WithCredentialsJSON([]byte(credentialsJSON)))
	}

	cfg := &firebasesdk.Config{ProjectID: projectID}
	app, err := firebasesdk.NewApp(ctx, cfg, opts...)
	if err != nil {
		return nil, fmt.Errorf("firebase: init app: %w", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("firebase: init auth client: %w", err)
	}

	return &authClientAdapter{client: client}, nil
}

// VerifyIDToken implements usecase.FirebaseTokenVerifier.
func (a *authClientAdapter) VerifyIDToken(ctx context.Context, idToken string) (*usecase.FirebaseToken, error) {
	tok, err := a.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("firebase: verify id token: %w", err)
	}

	email, _ := tok.Claims["email"].(string)
	name, _ := tok.Claims["name"].(string)
	photoURL, _ := tok.Claims["picture"].(string)

	return &usecase.FirebaseToken{
		UID:      tok.UID,
		Email:    email,
		Name:     name,
		PhotoURL: photoURL,
		Claims:   tok.Claims,
	}, nil
}

// GetUserByEmail implements usecase.FirebaseAdminClient.
func (a *authClientAdapter) GetUserByEmail(ctx context.Context, email string) (string, error) {
	user, err := a.client.GetUserByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("firebase: get user by email: %w", err)
	}
	return user.UID, nil
}

// UpdateUserPassword implements usecase.FirebaseAdminClient.
func (a *authClientAdapter) UpdateUserPassword(ctx context.Context, uid, newPassword string) error {
	params := (&auth.UserToUpdate{}).Password(newPassword)
	if _, err := a.client.UpdateUser(ctx, uid, params); err != nil {
		return fmt.Errorf("firebase: update user password: %w", err)
	}
	return nil
}
