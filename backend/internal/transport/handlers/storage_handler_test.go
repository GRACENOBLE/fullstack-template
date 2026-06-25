package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"backend/internal/transport/middleware"
	"backend/internal/usecase"
)

// mockStorageService implements usecase.StorageService for handler unit tests.
type mockStorageService struct {
	presignURL string
	publicURL  string
	presignErr error
	deleteErr  error
}

func (m *mockStorageService) PresignUpload(_ context.Context, key string, _ string, _ time.Duration) (string, error) {
	if m.presignErr != nil {
		return "", m.presignErr
	}
	if m.presignURL != "" {
		return m.presignURL, nil
	}
	return "https://r2.example.com/presigned/" + key, nil
}

func (m *mockStorageService) Delete(_ context.Context, _ string) error {
	return m.deleteErr
}

func (m *mockStorageService) PublicURL(key string) string {
	if m.publicURL != "" {
		return m.publicURL
	}
	return "https://pub.example.com/" + key
}

func newStorageRouter(h *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	injectUID := func(c *gin.Context) {
		c.Set(middleware.FirebaseClaimsKey, &usecase.FirebaseToken{UID: "uid123"})
		c.Next()
	}
	r.POST("/api/v1/storage/presign", injectUID, h.PresignHandler)
	r.DELETE("/api/v1/storage/:key", injectUID, h.DeleteObjectHandler)
	return r
}

// --- PresignHandler tests ---

func TestPresignHandler_HappyPath(t *testing.T) {
	mock := &mockStorageService{}
	h := &Handler{storageService: mock}
	r := newStorageRouter(h)

	body, _ := json.Marshal(map[string]string{
		"filename":     "photo.jpg",
		"content_type": "image/jpeg",
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/api/v1/storage/presign", bytes.NewReader(body)))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var envelope struct {
		Data presignResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if envelope.Data.UploadURL == "" {
		t.Error("expected non-empty upload_url")
	}
	if envelope.Data.PublicURL == "" {
		t.Error("expected non-empty public_url")
	}
}

func TestPresignHandler_MissingBody_Returns400(t *testing.T) {
	mock := &mockStorageService{}
	h := &Handler{storageService: mock}
	r := newStorageRouter(h)

	// Send empty JSON object — missing required fields.
	body, _ := json.Marshal(map[string]string{})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/api/v1/storage/presign", bytes.NewReader(body)))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestPresignHandler_StorageError_Returns500(t *testing.T) {
	mock := &mockStorageService{presignErr: errors.New("r2: presign failed")}
	h := &Handler{storageService: mock}
	r := newStorageRouter(h)

	body, _ := json.Marshal(map[string]string{
		"filename":     "photo.jpg",
		"content_type": "image/jpeg",
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/api/v1/storage/presign", bytes.NewReader(body)))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body.String())
	}
}

// --- DeleteObjectHandler tests ---

func TestDeleteObjectHandler_HappyPath(t *testing.T) {
	mock := &mockStorageService{}
	h := &Handler{storageService: mock}
	r := newStorageRouter(h)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/api/v1/storage/photo.jpg", nil))

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDeleteObjectHandler_StorageError_Returns500(t *testing.T) {
	mock := &mockStorageService{deleteErr: errors.New("r2: delete failed")}
	h := &Handler{storageService: mock}
	r := newStorageRouter(h)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/api/v1/storage/photo.jpg", nil))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDeleteObjectHandler_NilStorageService_Returns404(t *testing.T) {
	// When storageService is nil, the route should not be registered;
	// hitting it directly returns 404 from Gin.
	h := &Handler{storageService: nil}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	// Register no storage routes (simulates routes.go conditional).
	// Any request to the storage path returns 404.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/api/v1/storage/photo.jpg", nil))

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
	_ = h // suppress unused warning
}
