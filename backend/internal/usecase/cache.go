package usecase

import (
	"context"
	"time"
)

// CacheService is the cache port that use cases depend on.
// Infrastructure implementations live in internal/infrastructure/cache/.
type CacheService interface {
	PingContext(ctx context.Context) error
	Get(ctx context.Context, key string) (string, bool, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	SetNX(ctx context.Context, key string, value string, ttl time.Duration) (bool, error)
	Close() error
}
