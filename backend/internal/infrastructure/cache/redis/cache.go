package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"backend/internal/usecase"
)

type cacheService struct {
	client *redis.Client
}

// New creates a CacheService backed by Redis at the given URL.
// The caller is responsible for calling Close when done.
func New(redisURL string) (usecase.CacheService, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	return &cacheService{client: redis.NewClient(opts)}, nil
}

func (c *cacheService) PingContext(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

func (c *cacheService) Get(ctx context.Context, key string) (string, bool, error) {
	val, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return val, true, nil
}

func (c *cacheService) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *cacheService) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *cacheService) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (c *cacheService) SetNX(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
	return c.client.SetNX(ctx, key, value, ttl).Result()
}

func (c *cacheService) Close() error {
	return c.client.Close()
}
