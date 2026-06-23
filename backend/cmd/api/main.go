package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/internal/bootstrap"
	"backend/internal/infrastructure/queue"
	"backend/internal/infrastructure/ws"
	"backend/internal/server"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	slog.Info("shutting down gracefully, press Ctrl+C again to force")
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	slog.Info("server exiting")
	done <- true
}

// @title			Blueprint API
// @version		1.0
// @description	Fullstack template REST API.
//
// @host		localhost:8080
// @BasePath	/
//
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Firebase ID token — prefix with "Bearer "
func main() {
	// Signal-aware context so SIGINT/SIGTERM cancels bootstrap probes immediately.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	app, err := bootstrap.Run(ctx)
	stop() // release signal handler; gracefulShutdown re-registers it

	if err != nil {
		fmt.Fprintf(os.Stderr, "startup failed: %v\n", err)
		os.Exit(1)
	}

	hubCtx, hubCancel := context.WithCancel(context.Background())
	hub := ws.NewHub()
	go hub.Run(hubCtx)

	var workerCancel context.CancelFunc
	if app.Config.RedisURL != "" {
		workerCtx, wCancel := context.WithCancel(context.Background())
		workerCancel = wCancel
		worker, err := queue.NewWorker(app.Config.RedisURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "startup failed: %v\n", err)
			os.Exit(1)
		}
		worker.Register(queue.TypeWelcomeEmail, queue.NewHandleWelcomeEmail(app.EmailSender))
		go func() {
			if err := worker.Run(workerCtx); err != nil {
				slog.Error("queue: worker error", "err", err)
			}
		}()
	}

	srv, err := server.NewServer(app, hub)
	if err != nil {
		fmt.Fprintf(os.Stderr, "startup failed: %v\n", err)
		os.Exit(1)
	}
	slog.Info("API docs", "url", fmt.Sprintf("http://localhost%s/swagger/index.html", srv.Addr))

	done := make(chan bool, 1)
	go gracefulShutdown(srv, done)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	<-done

	if workerCancel != nil {
		workerCancel() // stop worker before hub (in-flight jobs drain first)
	}
	hubCancel() // stop hub after all WS connections have been closed by server shutdown

	if app.Cache != nil {
		if err := app.Cache.Close(); err != nil {
			slog.Error("cache close error", "error", err)
		}
	}
	if app.Enqueuer != nil {
		if err := app.Enqueuer.Close(); err != nil {
			slog.Error("enqueuer close error", "error", err)
		}
	}
	slog.Info("graceful shutdown complete")
}
