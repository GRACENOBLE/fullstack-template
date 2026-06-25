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
	req.RemoteAddr = "8.8.8.8:12345"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 from public IP, got %d", w.Code)
	}
}

func TestPprofIndex_XForwardedForSpoofingBlocked(t *testing.T) {
	h := &Handler{}
	handler := h.RegisterRoutes(0, 0, "", []string{"http://localhost:3000"})

	// Attacker connects from a public IP but spoofs X-Forwarded-For: 127.0.0.1.
	// LocalNetworkOnly uses RemoteAddr, not ClientIP(), so spoofing must not bypass the check.
	req := httptest.NewRequest(http.MethodGet, "/debug/pprof/", nil)
	req.RemoteAddr = "8.8.8.8:12345"
	req.Header.Set("X-Forwarded-For", "127.0.0.1")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 when X-Forwarded-For is spoofed, got %d (spoofing bypass!)", w.Code)
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
