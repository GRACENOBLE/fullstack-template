package ipgeo

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"backend/internal/usecase"
)

// mockCacheService is an in-memory implementation of usecase.CacheService for tests.
type mockCacheService struct {
	store map[string]string
}

func newMockCache() *mockCacheService {
	return &mockCacheService{store: make(map[string]string)}
}

func (m *mockCacheService) Get(_ context.Context, key string) (string, bool, error) {
	v, ok := m.store[key]
	return v, ok, nil
}

func (m *mockCacheService) Set(_ context.Context, key string, value string, _ time.Duration) error {
	m.store[key] = value
	return nil
}

func (m *mockCacheService) PingContext(_ context.Context) error { return nil }
func (m *mockCacheService) Delete(_ context.Context, _ string) error {
	return nil
}
func (m *mockCacheService) Exists(_ context.Context, _ string) (bool, error) { return false, nil }
func (m *mockCacheService) SetNX(_ context.Context, _ string, _ string, _ time.Duration) (bool, error) {
	return false, nil
}
func (m *mockCacheService) Close() error { return nil }

// Compile-time check: mockCacheService satisfies usecase.CacheService.
var _ usecase.CacheService = (*mockCacheService)(nil)

func TestClient_Lookup_ValidIP(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ipapiResponse{
			IP:          "8.8.8.8",
			CountryCode: "US",
			CountryName: "United States",
			Region:      "California",
			City:        "Mountain View",
			Timezone:    "America/Los_Angeles",
			Currency:    "USD",
			InEU:        false,
			Error:       false,
		}); err != nil {
			t.Errorf("encode response: %v", err)
		}
	}))
	defer srv.Close()

	client := NewWithBaseURL(nil, "", srv.URL)
	geo, err := client.Lookup(context.Background(), "8.8.8.8")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if geo.IP != "8.8.8.8" {
		t.Errorf("IP: got %q, want %q", geo.IP, "8.8.8.8")
	}
	if geo.CountryCode != "US" {
		t.Errorf("CountryCode: got %q, want %q", geo.CountryCode, "US")
	}
	if geo.CountryName != "United States" {
		t.Errorf("CountryName: got %q, want %q", geo.CountryName, "United States")
	}
	if geo.Region != "California" {
		t.Errorf("Region: got %q, want %q", geo.Region, "California")
	}
	if geo.City != "Mountain View" {
		t.Errorf("City: got %q, want %q", geo.City, "Mountain View")
	}
	if geo.Timezone != "America/Los_Angeles" {
		t.Errorf("Timezone: got %q, want %q", geo.Timezone, "America/Los_Angeles")
	}
	if geo.Currency != "USD" {
		t.Errorf("Currency: got %q, want %q", geo.Currency, "USD")
	}
	if geo.IsEU {
		t.Errorf("IsEU: got true, want false")
	}
}

func TestClient_Lookup_PrivateIP(t *testing.T) {
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := NewWithBaseURL(nil, "", srv.URL)

	privateIPs := []string{"127.0.0.1", "10.0.0.1", "192.168.1.1", "::1"}
	for _, ip := range privateIPs {
		t.Run(ip, func(t *testing.T) {
			_, err := client.Lookup(context.Background(), ip)
			if !errors.Is(err, ErrPrivateIP) {
				t.Errorf("ip %s: got error %v, want ErrPrivateIP", ip, err)
			}
		})
	}

	if called {
		t.Error("HTTP server was called for a private IP — should have short-circuited")
	}
}

func TestClient_Lookup_APIErrorInBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ipapiResponse{
			Error:  true,
			Reason: "invalid IP address",
		}); err != nil {
			t.Errorf("encode response: %v", err)
		}
	}))
	defer srv.Close()

	client := NewWithBaseURL(nil, "", srv.URL)
	_, err := client.Lookup(context.Background(), "1.2.3.4")
	if err == nil {
		t.Fatal("expected error when api body has error:true, got nil")
	}
}

func TestClient_Lookup_Non200Status(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()

	client := NewWithBaseURL(nil, "", srv.URL)
	_, err := client.Lookup(context.Background(), "1.2.3.4")
	if err == nil {
		t.Fatal("expected error for non-200 status, got nil")
	}
}

func TestClient_Lookup_CacheHit(t *testing.T) {
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ipapiResponse{
			IP:          "8.8.8.8",
			CountryCode: "US",
			CountryName: "United States",
			Region:      "California",
			City:        "Mountain View",
			Timezone:    "America/Los_Angeles",
			Currency:    "USD",
			InEU:        false,
		}); err != nil {
			t.Errorf("encode response: %v", err)
		}
	}))
	defer srv.Close()

	cache := newMockCache()
	client := NewWithBaseURL(cache, "", srv.URL)

	// First call — should hit HTTP.
	geo1, err := client.Lookup(context.Background(), "8.8.8.8")
	if err != nil {
		t.Fatalf("first lookup error: %v", err)
	}
	if callCount != 1 {
		t.Errorf("expected 1 HTTP call after first lookup, got %d", callCount)
	}

	// Second call — should be served from cache, no new HTTP call.
	geo2, err := client.Lookup(context.Background(), "8.8.8.8")
	if err != nil {
		t.Fatalf("second lookup error: %v", err)
	}
	if callCount != 1 {
		t.Errorf("expected still 1 HTTP call after cache hit, got %d", callCount)
	}

	if geo1.CountryCode != geo2.CountryCode {
		t.Errorf("cache returned different CountryCode: first=%q second=%q", geo1.CountryCode, geo2.CountryCode)
	}
}

func TestClient_Lookup_APIKey(t *testing.T) {
	var capturedURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedURL = r.URL.String()
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ipapiResponse{
			IP:          "1.2.3.4",
			CountryCode: "DE",
		}); err != nil {
			t.Errorf("encode response: %v", err)
		}
	}))
	defer srv.Close()

	client := NewWithBaseURL(nil, "my-secret-key", srv.URL)
	_, err := client.Lookup(context.Background(), "1.2.3.4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(capturedURL, "key=my-secret-key") {
		t.Errorf("expected ?key=my-secret-key in URL, got %q", capturedURL)
	}
}
