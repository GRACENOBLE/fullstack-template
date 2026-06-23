package usecase

import "context"

// StreamProducer publishes domain events to a named stream.
type StreamProducer interface {
	Publish(ctx context.Context, stream string, event any) error
}
