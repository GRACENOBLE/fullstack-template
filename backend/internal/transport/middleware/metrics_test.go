package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func init() { gin.SetMode(gin.TestMode) }

func TestPrometheusMiddleware_RecordsMetrics(t *testing.T) {
	r := gin.New()
	r.Use(PrometheusMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	count := testutil.ToFloat64(httpRequestsTotal.WithLabelValues("GET", "/test", "200"))
	if count != 1 {
		t.Errorf("expected counter to be 1, got %v", count)
	}
}

func TestPrometheusMiddleware_SkipsMetricsEndpoint(t *testing.T) {
	r := gin.New()
	r.Use(PrometheusMiddleware())
	r.GET("/metrics", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	before := testutil.ToFloat64(httpRequestsTotal.WithLabelValues("GET", "/metrics", "200"))

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	after := testutil.ToFloat64(httpRequestsTotal.WithLabelValues("GET", "/metrics", "200"))
	if before != after {
		t.Errorf("metrics endpoint should not be instrumented: before=%v after=%v", before, after)
	}
}

func TestPrometheusMiddleware_UnmatchedRoute(t *testing.T) {
	r := gin.New()
	r.Use(PrometheusMiddleware())
	// no routes registered — any request is unmatched

	req := httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Should be labelled "unmatched", not panic
	count := testutil.ToFloat64(httpRequestsTotal.WithLabelValues("GET", "unmatched", "404"))
	if count < 1 {
		t.Errorf("expected unmatched counter >= 1, got %v", count)
	}
}
