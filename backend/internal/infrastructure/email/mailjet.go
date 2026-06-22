package email

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"

	mailjet "github.com/mailjet/mailjet-apiv3-go/v4"
)

//go:embed templates/welcome.html
var templateFS embed.FS

// MailjetSender implements usecase.EmailSender using the Mailjet
// transactional email API (Send API v3.1).
type MailjetSender struct {
	client      *mailjet.Client
	fromEmail   string
	fromName    string
	sandboxMode bool
}

// NewMailjetSender constructs a MailjetSender with the supplied credentials.
// Provide a non-empty baseURL to override the Mailjet endpoint (useful in tests).
func NewMailjetSender(apiKey, secretKey, fromEmail, fromName string, baseURL ...string) *MailjetSender {
	client := mailjet.NewMailjetClient(apiKey, secretKey, baseURL...)
	return &MailjetSender{
		client:    client,
		fromEmail: fromEmail,
		fromName:  fromName,
	}
}

// NewSandboxSender constructs a MailjetSender that sends all messages in
// sandbox mode (Mailjet validates the request but does not deliver the email).
// Intended for integration tests against the live Mailjet API.
func NewSandboxSender(apiKey, secretKey, fromEmail, fromName string) *MailjetSender {
	return &MailjetSender{
		client:      mailjet.NewMailjetClient(apiKey, secretKey),
		fromEmail:   fromEmail,
		fromName:    fromName,
		sandboxMode: true,
	}
}

// SendWelcomeEmail sends a welcome email to toEmail rendered from the
// embedded welcome.html template.
func (s *MailjetSender) SendWelcomeEmail(_ context.Context, toEmail, toName string) error {
	html, err := renderWelcomeTemplate(toName)
	if err != nil {
		return fmt.Errorf("email: render welcome template: %w", err)
	}

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: s.fromEmail,
				Name:  s.fromName,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: toEmail,
					Name:  toName,
				},
			},
			Subject:  "Welcome to MyApp!",
			HTMLPart: html,
		},
	}

	req := &mailjet.MessagesV31{
		Info:        messagesInfo,
		SandBoxMode: s.sandboxMode,
	}

	if _, err := s.client.SendMailV31(req); err != nil {
		return fmt.Errorf("email: mailjet send: %w", err)
	}
	return nil
}

// renderWelcomeTemplate executes the embedded welcome.html template with the
// recipient name and returns the rendered HTML string.
func renderWelcomeTemplate(name string) (string, error) {
	raw, err := templateFS.ReadFile("templates/welcome.html")
	if err != nil {
		return "", fmt.Errorf("read welcome template: %w", err)
	}

	tmpl, err := template.New("welcome").Parse(string(raw))
	if err != nil {
		return "", fmt.Errorf("parse welcome template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]string{"Name": name}); err != nil {
		return "", fmt.Errorf("execute welcome template: %w", err)
	}
	return buf.String(), nil
}
