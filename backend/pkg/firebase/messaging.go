package firebase

import (
	"context"
	"fmt"

	firebasesdk "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"

	"backend/internal/usecase"
)

type messagingAdapter struct {
	client *messaging.Client
}

// NewMessagingClient returns a usecase.NotificationSender backed by FCM HTTP v1.
// Use NewApp to create the app so the same SDK instance is shared with NewAuthClient.
func NewMessagingClient(ctx context.Context, app *firebasesdk.App) (usecase.NotificationSender, error) {
	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("firebase: init messaging client: %w", err)
	}
	return &messagingAdapter{client: client}, nil
}

func (m *messagingAdapter) SendToToken(ctx context.Context, token, title, body string, data map[string]string) error {
	msg := &messaging.Message{
		Token:        token,
		Notification: &messaging.Notification{Title: title, Body: body},
		Data:         data,
	}
	if _, err := m.client.Send(ctx, msg); err != nil {
		return fmt.Errorf("fcm: send to token: %w", err)
	}
	return nil
}

func (m *messagingAdapter) SendMulticast(ctx context.Context, tokens []string, title, body string, data map[string]string) error {
	msg := &messaging.MulticastMessage{
		Tokens:       tokens,
		Notification: &messaging.Notification{Title: title, Body: body},
		Data:         data,
	}
	br, err := m.client.SendEachForMulticast(ctx, msg)
	if err != nil {
		return fmt.Errorf("fcm: send multicast: %w", err)
	}
	if br.FailureCount > 0 {
		return fmt.Errorf("fcm: %d/%d messages failed", br.FailureCount, len(tokens))
	}
	return nil
}
