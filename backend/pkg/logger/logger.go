package logger

import (
	"log/slog"
	"os"
)

// New returns a structured logger configured for the given environment.
// JSON format is used for staging and production (suitable for log aggregators).
// Human-readable text format is used for all other values (local dev).
func New(env string) *slog.Logger {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	switch env {
	case "staging", "production":
		return slog.New(slog.NewJSONHandler(os.Stdout, opts))
	default:
		return slog.New(slog.NewTextHandler(os.Stdout, opts))
	}
}
