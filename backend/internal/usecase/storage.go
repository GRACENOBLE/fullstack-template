package usecase

import (
	"context"
	"time"
)

// StorageService is the application-layer interface for object storage.
type StorageService interface {
	PresignUpload(ctx context.Context, key string, contentType string, ttl time.Duration) (string, error)
	Delete(ctx context.Context, key string) error
	PublicURL(key string) string
}
