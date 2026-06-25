package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPprofIndex_LoopbackAllowed(t *testing.T) {
	h := &Handler{}
	handler := h.RegisterRoutes(0, 0, "", []string{"http://localhost:3000"})

	req := httptest.NewRequest(http.MethodGet, "/debug/pprof/", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 from loopback, got %d", w.Code)
	}
}

func TestPprofIndex_PublicIPForbidden(t *testing.T) {
	h := &Handler{}
	handler := h.RegisterRoutes(0, 0, "", []string{"http://localhost:3000"})

	req := httptest.NewRequest(http.MethodGet, "/debug/pprof/", nil)
	// Gin's ClientIP() respects X-Forwarded-For when RemoteAddr is trusted loopback.
	// Set RemoteAddr to a public IP so LocalNetworkOnly() blocks it directly.
	req.RemoteAddr = "8.8.8.8:12345"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 from public IP, got %d", w.Code)
	}
}

func TestPprofHeap_LoopbackAllowed(t *testing.T) {
	h := &Handler{}
	handler := h.RegisterRoutes(0, 0, "", []string{"http://localhost:3000"})

	req := httptest.NewRequest(http.MethodGet, "/debug/pprof/heap", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for /debug/pprof/heap from loopback, got %d", w.Code)
	}
}
