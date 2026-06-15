package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/internal/transport/middleware"
	"github.com/gin-gonic/gin"
)

func init() { gin.SetMode(gin.TestMode) }

func newTestRouter(rps float64, burst int) *gin.Engine {
	r := gin.New()
	r.Use(middleware.RateLimit(rps, burst))
	r.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })
	return r
}

func TestRateLimit_Disabled(t *testing.T) {
	// rps=0 → no-op, all requests pass
	r := newTestRouter(0, 0)
	for i := 0; i < 20; i++ {
		req, _ := http.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("request %d: got %d, want 200", i, rr.Code)
		}
	}
}

func TestRateLimit_AllowsUnderLimit(t *testing.T) {
	// burst=5, rps=100 → first 5 requests all pass (burst absorbs them)
	r := newTestRouter(100, 5)
	for i := 0; i < 5; i++ {
		req, _ := http.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("request %d: got %d, want 200", i, rr.Code)
		}
	}
}

func TestRateLimit_BlocksOverLimit(t *testing.T) {
	// burst=1, rps=1 → second request immediately blocked
	r := newTestRouter(1, 1)
	req1, _ := http.NewRequest("GET", "/", nil)
	rr1 := httptest.NewRecorder()
	r.ServeHTTP(rr1, req1)
	if rr1.Code != http.StatusOK {
		t.Fatalf("first request: got %d, want 200", rr1.Code)
	}

	req2, _ := http.NewRequest("GET", "/", nil)
	rr2 := httptest.NewRecorder()
	r.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusTooManyRequests {
		t.Fatalf("second request: got %d, want 429", rr2.Code)
	}
}

func TestRateLimit_FractionalRPS_AllowsFirst(t *testing.T) {
	// rps=0.1 → int(0.1)*5 == 0; without the burst clamp this would block every request.
	// With clamp to 1, the first request must pass.
	r := newTestRouter(0.1, 0)
	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("fractional rps first request: got %d, want 200", rr.Code)
	}
}

func TestRateLimit_PerIP(t *testing.T) {
	// Each IP gets its own limiter; exhausting one IP does not affect another.
	// RemoteAddr is what Gin's ClientIP() resolves to when no trusted proxy is
	// configured (the default in tests), so set it directly to simulate distinct
	// clients rather than relying on X-Forwarded-For which requires proxy trust.
	r := newTestRouter(1, 1)

	sendFrom := func(ip string) int {
		req, _ := http.NewRequest("GET", "/", nil)
		req.RemoteAddr = ip + ":12345"
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		return rr.Code
	}

	// Exhaust IP-A
	if code := sendFrom("10.0.0.1"); code != http.StatusOK {
		t.Fatalf("IP-A first request: got %d, want 200", code)
	}
	if code := sendFrom("10.0.0.1"); code != http.StatusTooManyRequests {
		t.Fatalf("IP-A second request: got %d, want 429", code)
	}

	// IP-B should still be allowed
	if code := sendFrom("10.0.0.2"); code != http.StatusOK {
		t.Fatalf("IP-B first request: got %d, want 200", code)
	}
}
