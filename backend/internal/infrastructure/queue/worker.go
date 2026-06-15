package queue

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hibiken/asynq"
	goredis "github.com/redis/go-redis/v9"
)

// Worker wraps an asynq.Server and routes tasks to registered handlers.
type Worker struct {
	server *asynq.Server
	mux    *asynq.ServeMux
}

// NewWorker creates a Worker using the given Redis URL.
func NewWorker(redisURL string) (*Worker, error) {
	opt, err := goredis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("queue: parse redis url: %w", err)
	}
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: opt.Addr, Password: opt.Password, DB: opt.DB},
		asynq.Config{
			Concurrency: 10,
			ErrorHandler: asynq.ErrorHandlerFunc(func(_ context.Context, t *asynq.Task, err error) {
				slog.Error("queue: task failed", "type", t.Type(), "err", err)
			}),
		},
	)
	return &Worker{server: srv, mux: asynq.NewServeMux()}, nil
}

// Register adds a handler for the given task type.
func (w *Worker) Register(taskType string, h asynq.Handler) {
	w.mux.Handle(taskType, h)
}

// Run starts processing until ctx is cancelled.
// Call this in its own goroutine — it blocks.
func (w *Worker) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		if err := w.server.Run(w.mux); err != nil {
			errCh <- err
		}
	}()
	select {
	case <-ctx.Done():
		w.server.Shutdown()
		return nil
	case err := <-errCh:
		return err
	}
}

// NewInspector returns a low-level asynq.Inspector for the same Redis instance.
// Used by Asynqmon to read queue state.
func NewInspector(redisURL string) (*asynq.Inspector, error) {
	opt, err := goredis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("queue: parse redis url: %w", err)
	}
	return asynq.NewInspector(asynq.RedisClientOpt{Addr: opt.Addr, Password: opt.Password, DB: opt.DB}), nil
}
