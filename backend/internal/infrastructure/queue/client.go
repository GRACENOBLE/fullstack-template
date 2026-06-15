package queue

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
	goredis "github.com/redis/go-redis/v9"

	"backend/internal/usecase"
)

type client struct {
	asynq *asynq.Client
}

// NewClient creates a usecase.Enqueuer backed by Asynq.
// redisURL must be a valid redis:// URI (same value as REDIS_URL).
func NewClient(redisURL string) (usecase.Enqueuer, error) {
	opt, err := goredis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("queue: parse redis url: %w", err)
	}
	c := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     opt.Addr,
		Password: opt.Password,
		DB:       opt.DB,
	})
	return &client{asynq: c}, nil
}

// Enqueue submits a task of the given type with the provided payload.
func (c *client) Enqueue(ctx context.Context, taskType string, payload []byte) error {
	task := asynq.NewTask(taskType, payload)
	_, err := c.asynq.EnqueueContext(ctx, task)
	return err
}

// Close releases the underlying Asynq client connection.
func (c *client) Close() error {
	return c.asynq.Close()
}
