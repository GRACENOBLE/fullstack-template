package usecase

import "context"

// Enqueuer is the task-queue port that use cases depend on.
// The Asynq-backed implementation lives in internal/infrastructure/queue/.
// nil disables background job processing (when REDIS_URL is not set).
type Enqueuer interface {
	Enqueue(ctx context.Context, taskType string, payload []byte) error
	Close() error
}
