package middleware_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"backend/internal/domain"
	"backend/internal/infrastructure/ipgeo"
	"backend/internal/transport/middleware"
	"backend/internal/usecase"
)

// mockGeoLocator is a test double implementing usecase.GeoLocator.
type mockGeoLocator struct {
	geo *domain.GeoLocation
	err error
}

func (m *mockGeoLocator) Lookup(_ context.Context, _ string) (*domain.GeoLocation, error) {
	return m.geo, m.err
}

// Compile-time check.
var _ usecase.GeoLocator = (*mockGeoLocator)(nil)

func TestGeoFromRequest_AttachesGeoOnSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	want := &domain.GeoLocation{
		IP:          "1.2.3.4",
		CountryCode: "US",
		CountryName: "United States",
		City:        "Mountain View",
	}
	locator := &mockGeoLocator{geo: want}

	var captured *domain.GeoLocation
	r := gin.New()
	r.Use(middleware.GeoFromRequest(locator))
	r.GET("/", func(c *gin.Context) {
		val, exists := c.Get(middleware.GeoLocationKey)
		if !exists {
			t.Error("geo_location key not set in context")
			c.Status(http.StatusInternalServerError)
			return
		}
		geo, ok := val.(*domain.GeoLocation)
		if !ok {
			t.Errorf("geo_location wrong type: %T", val)
			c.Status(http.StatusInternalServerError)
			return
		}
		captured = geo
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "1.2.3.4:1234"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if captured == nil || captured.CountryCode != want.CountryCode {
		t.Errorf("captured geo mismatch: got %+v, want %+v", captured, want)
	}
}

func TestGeoFromRequest_SkipsOnError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	locator := &mockGeoLocator{err: errors.New("rate limited")}

	reached := false
	r := gin.New()
	r.Use(middleware.GeoFromRequest(locator))
	r.GET("/", func(c *gin.Context) {
		reached = true
		_, exists := c.Get(middleware.GeoLocationKey)
		if exists {
			t.Error("geo_location key should not be set on error")
		}
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "1.2.3.4:1234"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if !reached {
		t.Error("handler was not reached — middleware aborted the request")
	}
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGeoFromRequest_PrivateIPSkipped(t *testing.T) {
	gin.SetMode(gin.TestMode)

	locator := &mockGeoLocator{err: ipgeo.ErrPrivateIP}

	reached := false
	r := gin.New()
	r.Use(middleware.GeoFromRequest(locator))
	r.GET("/", func(c *gin.Context) {
		reached = true
		_, exists := c.Get(middleware.GeoLocationKey)
		if exists {
			t.Error("geo_location key should not be set for private IP")
		}
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "127.0.0.1:5678"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if !reached {
		t.Error("handler was not reached — middleware aborted the request")
	}
}

func TestRealIP_XForwardedFor(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Forwarded-For", "1.2.3.4, 10.0.0.1")
	req.RemoteAddr = "10.0.0.1:9999"

	got := middleware.RealIP(req)
	if got != "1.2.3.4" {
		t.Errorf("RealIP: got %q, want %q", got, "1.2.3.4")
	}
}

func TestRealIP_XRealIP(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Real-IP", "5.6.7.8")
	req.RemoteAddr = "10.0.0.1:9999"

	got := middleware.RealIP(req)
	if got != "5.6.7.8" {
		t.Errorf("RealIP: got %q, want %q", got, "5.6.7.8")
	}
}

func TestRealIP_RemoteAddr(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "9.10.11.12:4567"

	got := middleware.RealIP(req)
	if got != "9.10.11.12" {
		t.Errorf("RealIP: got %q, want %q", got, "9.10.11.12")
	}
}
