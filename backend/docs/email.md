---
title: Email (Mailjet)
last_verified: 2026-06-23
sources:
  - internal/usecase/email.go
  - internal/infrastructure/email/mailjet.go
  - internal/infrastructure/email/templates/welcome.html
  - internal/bootstrap/bootstrap.go
---

# Email (Mailjet)

## Overview

`usecase.EmailSender` is the interface all email sending goes through:

```go
// internal/usecase/email.go
type EmailSender interface {
    SendWelcomeEmail(ctx context.Context, toEmail, toName string) error
}
```

The interface lives in `usecase/` because that is the layer that depends on it — handlers and future use-case logic call `EmailSender` without knowing about Mailjet.

`App.EmailSender` is `nil` when `MAILJET_API_KEY` is omitted, so callers must nil-guard before use.

## Implementation

`MailjetSender` in `internal/infrastructure/email/mailjet.go` implements the interface using the Mailjet Send API v3.1.

### Constructors

```go
// Production — sends real email
func NewMailjetSender(apiKey, secretKey, fromEmail, fromName string, baseURL ...string) *MailjetSender

// Sandbox — Mailjet validates the payload but does not deliver; for integration tests
func NewSandboxSender(apiKey, secretKey, fromEmail, fromName string) *MailjetSender
```

The optional `baseURL` variadic on `NewMailjetSender` overrides the Mailjet API endpoint, which allows unit tests to point the sender at an `httptest.Server`.

### HTML templates

Templates are embedded at compile time via `//go:embed`:

```go
//go:embed templates/welcome.html
var templateFS embed.FS
```

`templates/welcome.html` is a Go `html/template` file. The only template data value currently used is `{{.Name}}` (the recipient's display name). `renderWelcomeTemplate` parses and executes the template on each call and returns the rendered HTML string.

### Wiring (bootstrap)

`bootstrap.Run` constructs the sender when both `MAILJET_API_KEY` and `MAILJET_SECRET_KEY` are non-empty:

```go
var emailSender usecase.EmailSender
if cfg.MailjetAPIKey != "" && cfg.MailjetSecretKey != "" {
    emailSender = email.NewMailjetSender(
        cfg.MailjetAPIKey,
        cfg.MailjetSecretKey,
        cfg.FromEmail,
        cfg.FromName,
    )
}
```

`emailSender` is stored on `bootstrap.App.EmailSender`. `server.NewServer` passes it to `handlers.NewHandler` as the last argument.

## Env vars

| Variable | Required | Description |
|---|---|---|
| `MAILJET_API_KEY` | No | Mailjet API key. Omit (or leave empty) to disable all email sending. |
| `MAILJET_SECRET_KEY` | No | Mailjet secret key. Must be set alongside `MAILJET_API_KEY`. |
| `FROM_EMAIL` | No | Verified sender address — e.g. `no-reply@example.com`. |
| `FROM_NAME` | No | Sender display name — e.g. `MyApp`. |

When `MAILJET_API_KEY` or `MAILJET_SECRET_KEY` is empty, `App.EmailSender` is `nil` and no email is sent. The other two vars are only read when the sender is initialised.

## How to add a new email type

1. **Add a method to the interface** in `internal/usecase/email.go`:
   ```go
   SendPasswordResetEmail(ctx context.Context, toEmail, resetURL string) error
   ```

2. **Add an HTML template** at `internal/infrastructure/email/templates/password_reset.html`. Use Go template syntax (`{{.FieldName}}`) for dynamic values.

3. **Implement the method** on `MailjetSender` in `internal/infrastructure/email/mailjet.go`. Follow the `SendWelcomeEmail` pattern: render the template, populate `mailjet.InfoMessagesV31`, call `client.SendMailV31`.

4. **Update the test double** in `internal/usecase/email_usecase_test.go` to add the new method stub so `mockEmailSender` continues to satisfy the interface at compile time.

5. **Write tests** — see the Testing section below.

## Testing

### Unit test — interface contract (`internal/usecase/`)

Test that a mock satisfies the interface and propagates errors correctly. No network calls.

```go
type mockEmailSender struct {
    calledWithEmail string
    err             error
}

func (m *mockEmailSender) SendWelcomeEmail(_ context.Context, toEmail, toName string) error {
    m.calledWithEmail = toEmail
    return m.err
}

var _ usecase.EmailSender = (*mockEmailSender)(nil) // compile-time check
```

See `internal/usecase/email_usecase_test.go` for the full example.

### Unit test — HTTP error propagation (`internal/infrastructure/email/`)

Use `httptest.NewServer` and pass its URL as the `baseURL` override to `NewMailjetSender`. No real credentials needed.

```go
srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
    w.WriteHeader(http.StatusUnauthorized)
    json.NewEncoder(w).Encode(map[string]any{"StatusCode": 401})
}))
defer srv.Close()

sender := email.NewMailjetSender("bad-key", "bad-secret", "no-reply@example.com", "Test", srv.URL+"/v3")
err := sender.SendWelcomeEmail(context.Background(), "user@example.com", "User")
// assert err != nil
```

### Integration test — live Mailjet sandbox

Uses `NewSandboxSender` which sets `SandBoxMode: true` on the Mailjet request. Mailjet validates the payload and returns a success response without delivering the email.

```go
func TestMailjetSender_SendWelcomeEmail_Sandbox(t *testing.T) {
    apiKey := os.Getenv("MAILJET_API_KEY")
    secretKey := os.Getenv("MAILJET_SECRET_KEY")
    if apiKey == "" || secretKey == "" {
        t.Skip("MAILJET_API_KEY and MAILJET_SECRET_KEY not set — skipping Mailjet sandbox integration test")
    }

    sender := email.NewSandboxSender(apiKey, secretKey, "no-reply@example.com", "MyApp Test")
    err := sender.SendWelcomeEmail(context.Background(), "sandbox@mailjet.com", "Sandbox User")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}
```

The skip guard means this test is silent in CI unless the credentials are injected as secrets.
