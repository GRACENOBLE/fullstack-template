package streams

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

// Handler processes a single stream message payload (raw JSON bytes).
type Handler func(ctx context.Context, data []byte) error

// Consumer reads from a Redis Stream via a consumer group.
type Consumer struct {
	client   *redis.Client
	stream   string
	group    string
	consumer string
}

// NewConsumer creates a Consumer for the given stream / group / consumer name.
func NewConsumer(redisURL, stream, group, consumer string) (*Consumer, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("streams: parse redis url: %w", err)
	}
	return &Consumer{
		client:   redis.NewClient(opts),
		stream:   stream,
		group:    group,
		consumer: consumer,
	}, nil
}

// Run reads and dispatches messages until ctx is cancelled.
// Creates the consumer group (and stream) if either does not exist.
func (c *Consumer) Run(ctx context.Context, h Handler) error {
	// MKSTREAM ensures the stream exists before the first publish.
	_ = c.client.XGroupCreateMkStream(ctx, c.stream, c.group, "0").Err()

	for {
		entries, err := c.client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    c.group,
			Consumer: c.consumer,
			Streams:  []string{c.stream, ">"},
			Count:    10,
			Block:    2000,
		}).Result()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			if err == redis.Nil {
				continue
			}
			slog.Error("streams: read error", "stream", c.stream, "err", err)
			continue
		}
		for _, entry := range entries {
			for _, msg := range entry.Messages {
				data, _ := msg.Values["data"].(string)
				if err := h(ctx, []byte(data)); err != nil {
					slog.Error("streams: handler error", "msg_id", msg.ID, "err", err)
				} else {
					_ = c.client.XAck(ctx, c.stream, c.group, msg.ID)
				}
			}
		}
	}
}

// Close releases the underlying Redis connection.
func (c *Consumer) Close() error {
	return c.client.Close()
}
