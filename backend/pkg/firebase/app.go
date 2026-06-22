package firebase

import (
	"context"
	"fmt"

	firebasesdk "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

// NewApp initialises the Firebase Admin SDK default app.
// credentialsJSON is the raw service account JSON (FIREBASE_SERVICE_ACCOUNT_JSON).
// When empty the SDK falls back to Application Default Credentials (ADC).
func NewApp(ctx context.Context, projectID, credentialsJSON string) (*firebasesdk.App, error) {
	var opts []option.ClientOption
	if credentialsJSON != "" {
		opts = append(opts, option.WithCredentialsJSON([]byte(credentialsJSON)))
	}
	cfg := &firebasesdk.Config{ProjectID: projectID}
	app, err := firebasesdk.NewApp(ctx, cfg, opts...)
	if err != nil {
		return nil, fmt.Errorf("firebase: init app: %w", err)
	}
	return app, nil
}
