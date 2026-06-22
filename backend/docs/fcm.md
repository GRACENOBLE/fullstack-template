---
topic: fcm
last_verified: 2026-06-16
sources:
  - internal/domain/fcm_token.go
  - internal/usecase/notification.go
  - internal/infrastructure/database/postgres/fcm_token_repository.go
  - internal/infrastructure/database/postgres/fcm_token_repository_test.go
  - internal/transport/handlers/fcm_handler.go
  - internal/transport/handlers/fcm_handler_test.go
  - internal/transport/handlers/routes.go
  - internal/transport/handlers/handler.go
  - internal/server/server.go
  - internal/bootstrap/bootstrap.go
  - pkg/firebase/app.go
  - pkg/firebase/admin.go
  - pkg/firebase/messaging.go
  - internal/infrastructure/database/migrations/20260615220942_add_fcm_tokens.sql
---

# Firebase Cloud Messaging (FCM)

FCM sends push notifications to Android, iOS, and web clients. The backend stores per-user device tokens and sends notifications via the Firebase Admin SDK using the FCM HTTP v1 API.

## Domain entity

`internal/domain/fcm_token.go`:
```go
type FCMToken struct {
    ID        string
    UserID    string
    Token     string
    Platform  string    // "android" | "ios" | "web"
    CreatedAt time.Time
}
```

## Usecase interfaces

Both interfaces live in `internal/usecase/notification.go`.

```go
// NotificationSender sends FCM push notifications via the Admin SDK.
type NotificationSender interface {
    SendToToken(ctx context.Context, token, title, body string, data map[string]string) error
    SendMulticast(ctx context.Context, tokens []string, title, body string, data map[string]string) error
}

// FCMTokenRepository persists device registration tokens per user.
type FCMTokenRepository interface {
    SaveToken(ctx context.Context, userID, token, platform string) error
    GetTokensByUserID(ctx context.Context, userID string) ([]domain.FCMToken, error)
    DeleteToken(ctx context.Context, token string) error
}
```

`NotificationSender` is the narrow port used by any use case that needs to push a notification. `FCMTokenRepository` is the persistence port used by the HTTP handlers.

## pkg/firebase — shared app initialisation

`pkg/firebase/app.go` creates the single Firebase Admin SDK app instance:

```go
func NewApp(ctx context.Context, projectID, credentialsJSON string) (*firebasesdk.App, error)
```

The same `*firebasesdk.App` is passed to both `NewAuthClient` and `NewMessagingClient` so the SDK initialises only once per process. Calling `firebasesdk.NewApp` twice with default settings returns an error, hence the shared app pattern.

`pkg/firebase/admin.go` — `NewAuthClient(ctx, app)` now accepts the pre-created app instead of creating its own.

`pkg/firebase/messaging.go` — `NewMessagingClient(ctx, app)` wraps `*messaging.Client`:

```go
func NewMessagingClient(ctx context.Context, app *firebasesdk.App) (usecase.NotificationSender, error)
```

`SendMulticast` uses `client.SendEachForMulticast` and returns an error if any individual message fails.

## bootstrap.App

`internal/bootstrap/bootstrap.go` — `App` now carries `FCMSender usecase.NotificationSender`. Both Firebase clients share the same SDK app:

```go
if cfg.FirebaseProjectID != "" {
    fbApp, _  := firebase.NewApp(ctx, cfg.FirebaseProjectID, cfg.FirebaseServiceAccountJSON)
    authClient, _ := firebase.NewAuthClient(ctx, fbApp)
    msgClient, _  := firebase.NewMessagingClient(ctx, fbApp)
    app.Firebase  = authClient
    app.FCMSender = msgClient
}
```

When `FIREBASE_PROJECT_ID` is not set both fields are `nil` and FCM routes are not registered.

## Database — fcm_tokens table

Migration `20260615220942_add_fcm_tokens.sql`:
```sql
CREATE TABLE IF NOT EXISTS fcm_tokens (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     TEXT         NOT NULL,
    token       TEXT         NOT NULL UNIQUE,
    platform    TEXT         NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_fcm_tokens_user_id ON fcm_tokens(user_id);
```

`SaveToken` uses an upsert (`ON CONFLICT (token) DO UPDATE`) so a physical device re-registering or switching accounts is handled atomically without a duplicate token.

## HTTP endpoints

Both routes are registered inside the `/api/v1` group and inherit the `FirebaseAuth` middleware when `h.verifier != nil`. They are only registered when `h.fcmTokenRepo != nil`:

```go
if h.fcmTokenRepo != nil {
    api.POST("/fcm/register", h.RegisterFCMToken)
    api.DELETE("/fcm/unregister", h.UnregisterFCMToken)
}
```

### POST /api/v1/fcm/register

Request body:
```json
{ "token": "<fcm-registration-token>", "platform": "android|ios|web" }
```

Reads `UID` from `*usecase.FirebaseToken` in the Gin context (set by `FirebaseAuth` middleware), then calls `SaveToken`. Returns `200 {"message": "token registered"}`.

### DELETE /api/v1/fcm/unregister

Request body:
```json
{ "token": "<fcm-registration-token>" }
```

Calls `DeleteToken`. Returns `200 {"message": "token unregistered"}`. Designed to be called on logout.

## Wiring in server.go

```go
fcmTokenRepo := postgres.NewFCMTokenRepository(app.DB)
h := handlers.NewHandler(..., app.FCMSender, fcmTokenRepo)
```

`NewFCMTokenRepository` always returns a valid repository backed by `app.DB`; the FCM routes are gated by `h.fcmTokenRepo != nil` on the router side, but since `NewFCMTokenRepository` is always called, the routes are always registered when the server starts. The `FirebaseAuth` middleware is the actual auth gate.

## Sending notifications from other use cases

Inject `usecase.NotificationSender` into any use case that needs to push a notification:

```go
// Example: notify all user devices after a background job completes
tokens, _ := fcmTokenRepo.GetTokensByUserID(ctx, userID)
fcmTokens := make([]string, len(tokens))
for i, t := range tokens { fcmTokens[i] = t.Token }
_ = fcmSender.SendMulticast(ctx, fcmTokens, "Job done", "Your report is ready", nil)
```

## Testing

**Handler unit tests** (`internal/transport/handlers/fcm_handler_test.go`): `mockFCMTokenRepo` implements `usecase.FCMTokenRepository`; no database or Firebase SDK involved.

**Repository integration tests** (`internal/infrastructure/database/postgres/fcm_token_repository_test.go`): use Testcontainers (real PostgreSQL). Each test calls `setupFCMTokensTable` to create the table and registers a cleanup to drop it. Shares the `testDB` global and `TestMain` from `health_repository_test.go` — do not add another `TestMain` to this package.

## Environment variables

| Variable | Required | Description |
|---|---|---|
| `FIREBASE_PROJECT_ID` | No | Firebase project ID (e.g. `my-app-12345`). Omit to disable Firebase entirely. |
| `FIREBASE_SERVICE_ACCOUNT_JSON` | No | Raw service account JSON string. Omit to use Application Default Credentials. |

Both are already documented in `backend/.env.example`.
