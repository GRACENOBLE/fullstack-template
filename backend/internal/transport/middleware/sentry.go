package middleware

import (
	"log/slog"

	sentry "github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

// SentryMiddleware reports panics and errors to Sentry.
// When dsn is empty it returns a no-op handler so the app works without Sentry configured.
// When dsn is set but sentry.Init fails (e.g. malformed DSN), the error is logged and
// the middleware degrades to a no-op rather than preventing startup.
func SentryMiddleware(dsn string) gin.HandlerFunc {
	if dsn == "" {
		return func(c *gin.Context) { c.Next() }
	}
	if err := sentry.Init(sentry.ClientOptions{Dsn: dsn}); err != nil {
		slog.Default().Warn("sentry: init failed, continuing without error tracking", "error", err)
		return func(c *gin.Context) { c.Next() }
	}
	return sentrygin.New(sentrygin.Options{Repanic: true})
}
