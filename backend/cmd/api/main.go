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
func main() {
	// Signal-aware context so SIGINT/SIGTERM cancels bootstrap probes immediately.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	app, err := bootstrap.Run(ctx)
	stop() // release signal handler; gracefulShutdown re-registers it

	if err != nil {
		fmt.Fprintf(os.Stderr, "startup failed: %v\n", err)
		os.Exit(1)
	}

	srv := server.NewServer(app)
	slog.Info("API docs", "url", fmt.Sprintf("http://localhost%s/swagger/index.html", srv.Addr))

	done := make(chan bool, 1)
	go gracefulShutdown(srv, done)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	<-done
	if app.Cache != nil {
		if err := app.Cache.Close(); err != nil {
			slog.Error("cache close error", "error", err)
		}
	}
	slog.Info("graceful shutdown complete")
}
