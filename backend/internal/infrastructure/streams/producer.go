package streams

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Producer appends events to Redis Streams.
type Producer struct {
	client *redis.Client
}

// NewProducer creates a Producer for the given Redis URL.
func NewProducer(redisURL string) (*Producer, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("streams: parse redis url: %w", err)
	}
	return &Producer{client: redis.NewClient(opts)}, nil
}

// Publish appends event to the named stream as a JSON "data" field.
func (p *Producer) Publish(ctx context.Context, stream string, event any) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("streams: marshal event: %w", err)
	}
	return p.client.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: map[string]any{"data": string(payload)},
	}).Err()
}

// Close releases the underlying Redis connection.
func (p *Producer) Close() error {
	return p.client.Close()
}
