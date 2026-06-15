package usecase

import (
	"context"

	"backend/internal/domain"
)

// NotificationSender sends FCM push notifications.
type NotificationSender interface {
	SendToToken(ctx context.Context, token, title, body string, data map[string]string) error
	SendMulticast(ctx context.Context, tokens []string, title, body string, data map[string]string) error
}

// FCMTokenRepository persists device registration tokens.
type FCMTokenRepository interface {
	SaveToken(ctx context.Context, userID, token, platform string) error
	GetTokensByUserID(ctx context.Context, userID string) ([]domain.FCMToken, error)
	DeleteToken(ctx context.Context, token string) error
}
