---
topic: auth
last_verified: 2026-06-15
sources:
  - internal/usecase/auth_usecase.go
  - internal/transport/middleware/auth.go
  - internal/transport/handlers/auth_handler.go
  - pkg/firebase/admin.go
  - internal/bootstrap/bootstrap.go
---

# Firebase Auth

## Domain type

`FirebaseToken` is defined in `internal/usecase/auth_usecase.go` and holds the claims extracted from a verified Firebase ID token:

```go
type FirebaseToken struct {
    UID      string         `json:"uid"`
    Email    string         `json:"email"`
    Name     string         `json:"name"`
    PhotoURL string         `json:"photoUrl"`
    Claims   map[string]any `json:"claims"`
}
```

`UID`, `Email`, `Name`, and `PhotoURL` are promoted from the raw Firebase JWT claims (`uid`, `email`, `name`, `picture`). `Claims` contains the full unmodified payload.

## Usecase interfaces

Both interfaces live in `internal/usecase/auth_usecase.go`.

```go
// FirebaseTokenVerifier can verify a raw Firebase ID token string.
type FirebaseTokenVerifier interface {
    VerifyIDToken(ctx context.Context, idToken string) (*FirebaseToken, error)
}

// FirebaseAdminClient extends FirebaseTokenVerifier with user-management operations.
type FirebaseAdminClient interface {
    FirebaseTokenVerifier
    GetUserByEmail(ctx context.Context, email string) (string, error)
    UpdateUserPassword(ctx context.Context, uid, newPassword string) error
}
```

`FirebaseTokenVerifier` is the narrow interface used by the `FirebaseAuth` middleware. `FirebaseAdminClient` is the broader interface stored on `bootstrap.App` and wired into the server.

## pkg/firebase/admin.go

`NewAuthClient` initialises the Firebase Admin SDK and returns a `usecase.FirebaseAdminClient`:

```go
func NewAuthClient(ctx context.Context, projectID, credentialsJSON string) (usecase.FirebaseAdminClient, error)
```

- When `credentialsJSON` is non-empty the SDK is initialised with `option.WithCredentialsJSON` (service account key).
- When `credentialsJSON` is empty the SDK falls back to Application Default Credentials (ADC) — used on GCP.

The returned value is `*authClientAdapter`, a private type that wraps `*auth.Client` from the Firebase Admin SDK. This adapter satisfies both `FirebaseTokenVerifier` and `FirebaseAdminClient` without leaking the SDK type into the application layer.

## FirebaseAuth middleware

Defined in `internal/transport/middleware/auth.go`.

```go
const FirebaseClaimsKey = "firebase_claims"

func FirebaseAuth(verifier usecase.FirebaseTokenVerifier) gin.HandlerFunc
```

For every request the middleware:
1. Reads the `Authorization` header; aborts with `401` if the header is missing or does not start with `Bearer `.
2. Calls `verifier.VerifyIDToken(ctx, idToken)`.
3. On success: stores `*usecase.FirebaseToken` in the Gin context under `FirebaseClaimsKey` and calls `c.Next()`.
4. On error: aborts with `401` and `{"error": "invalid or expired token"}`.

Retrieve claims inside a handler:
```go
val, _ := c.Get(middleware.FirebaseClaimsKey)
token, ok := val.(*usecase.FirebaseToken)
```

## MeHandler (GET /api/v1/me)

Defined in `internal/transport/handlers/auth_handler.go`.

```go
func (h *Handler) MeHandler(c *gin.Context)
```

- Reads `*usecase.FirebaseToken` from the Gin context (`FirebaseClaimsKey`).
- Returns `200 OK` with the token struct serialised as JSON.
- Returns `401 Unauthorized` with `{"error": "unauthorized"}` if the context value is missing or of the wrong type (should not happen when `FirebaseAuth` is applied to the group).

The handler is registered on the `/api/v1` group in `RegisterRoutes`:
```go
api := r.Group("/api/v1")
if verifier != nil {
    api.Use(middleware.FirebaseAuth(verifier))
}
api.GET("/me", h.MeHandler)
```

## Disabling auth in development

When `FIREBASE_PROJECT_ID` is not set `bootstrap.Run` skips Firebase initialisation and `app.Firebase` is `nil`. `server.NewServer` passes `app.Firebase` directly to `RegisterRoutes` as the `verifier` argument. When `verifier` is `nil` the `if verifier != nil` guard in `RegisterRoutes` skips `api.Use(middleware.FirebaseAuth(...))`, so `/api/v1/me` is reachable without a token.

To enable auth locally set both `FIREBASE_PROJECT_ID` and `FIREBASE_SERVICE_ACCOUNT_JSON` in `backend/.env`.
