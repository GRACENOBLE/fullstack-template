package r2_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"backend/internal/infrastructure/storage/r2"
)

// roundTripFunc allows a plain function to satisfy http.RoundTripper.
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestNew_MissingFields(t *testing.T) {
	cases := []struct {
		name                                         string
		accountID, accessKey, secretKey, bucket, pub string
	}{
		{"missing accountID", "", "ak", "sk", "bucket", "https://pub.example.com"},
		{"missing accessKey", "acct", "", "sk", "bucket", "https://pub.example.com"},
		{"missing secretKey", "acct", "ak", "", "bucket", "https://pub.example.com"},
		{"missing bucket", "acct", "ak", "sk", "", "https://pub.example.com"},
		{"missing publicBaseURL", "acct", "ak", "sk", "bucket", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := r2.New(tc.accountID, tc.accessKey, tc.secretKey, tc.bucket, tc.pub)
			if err == nil {
				t.Fatal("expected error for missing field, got nil")
			}
		})
	}
}

func TestNew_AllFieldsPresent(t *testing.T) {
	svc, err := r2.New("acct123", "AKID", "SECRET", "my-bucket", "https://pub.example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc == nil {
		t.Fatal("expected non-nil StorageService")
	}
}

func TestPresignUpload_ReturnsURLContainingKey(t *testing.T) {
	// PresignUpload generates the URL locally without an HTTP call, so we just
	// need a valid-looking service instance.
	svc, err := r2.New("acct123", "AKID", "SECRET", "my-bucket", "https://pub.example.com")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	url, err := svc.PresignUpload(context.Background(), "uploads/photo.jpg", "image/jpeg", 15*time.Minute)
	if err != nil {
		t.Fatalf("PresignUpload: %v", err)
	}
	if url == "" {
		t.Fatal("expected non-empty presigned URL")
	}
	if !strings.Contains(url, "photo.jpg") {
		t.Errorf("expected URL to contain key, got: %s", url)
	}
}

func TestDelete_CallsServerWithDELETE(t *testing.T) {
	var capturedMethod, capturedPath string

	// httptest server that records the request and responds 204.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedMethod = r.Method
		capturedPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	httpClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			// Rewrite the host to point at the test server.
			req.URL.Host = strings.TrimPrefix(srv.URL, "http://")
			req.URL.Scheme = "http"
			return http.DefaultTransport.RoundTrip(req)
		}),
	}

	svc, err := r2.NewWithHTTPClient("acct123", "AKID", "SECRET", "my-bucket", "https://pub.example.com", httpClient)
	if err != nil {
		t.Fatalf("NewWithHTTPClient: %v", err)
	}

	if err := svc.Delete(context.Background(), "uploads/photo.jpg"); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	if capturedMethod != http.MethodDelete {
		t.Errorf("expected DELETE request, got %s", capturedMethod)
	}
	if !strings.Contains(capturedPath, "photo.jpg") {
		t.Errorf("expected path to contain key, got: %s", capturedPath)
	}
}

func TestPublicURL_ReturnsBaseURLPlusKey(t *testing.T) {
	svc, err := r2.New("acct123", "AKID", "SECRET", "my-bucket", "https://pub.example.com")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	// url.PathEscape encodes both the space and the slash within the key,
	// because the entire key is a single path segment from the SDK's perspective.
	got := svc.PublicURL("uploads/my file.jpg")
	want := "https://pub.example.com/uploads%2Fmy%20file.jpg"
	if got != want {
		t.Errorf("PublicURL: got %q, want %q", got, want)
	}
}

func TestPublicURL_SimpleKey(t *testing.T) {
	svc, err := r2.New("acct123", "AKID", "SECRET", "my-bucket", "https://pub.example.com")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	got := svc.PublicURL("photo.jpg")
	want := "https://pub.example.com/photo.jpg"
	if got != want {
		t.Errorf("PublicURL: got %q, want %q", got, want)
	}
}
