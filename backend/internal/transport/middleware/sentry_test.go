package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"backend/internal/transport/middleware"
)

func TestSentryMiddleware_NoDSN_IsNoOp(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.SentryMiddleware(""))
	r.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestSentryMiddleware_WithDSN_ReturnsHandler(t *testing.T) {
	// Use a syntactically valid but non-connectable DSN
	h := middleware.SentryMiddleware("https://key@o0.ingest.sentry.io/0")
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}
