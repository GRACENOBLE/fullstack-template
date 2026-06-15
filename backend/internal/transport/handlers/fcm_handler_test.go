package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"backend/internal/domain"
	"backend/internal/transport/middleware"
	"backend/internal/usecase"
)

type mockFCMTokenRepo struct {
	savedUserID   string
	savedToken    string
	savedPlatform string
	deletedToken  string
	tokens        []domain.FCMToken
	saveErr       error
	deleteErr     error
}

func (m *mockFCMTokenRepo) SaveToken(_ context.Context, userID, token, platform string) error {
	m.savedUserID, m.savedToken, m.savedPlatform = userID, token, platform
	return m.saveErr
}
func (m *mockFCMTokenRepo) GetTokensByUserID(_ context.Context, _ string) ([]domain.FCMToken, error) {
	return m.tokens, nil
}
func (m *mockFCMTokenRepo) DeleteToken(_ context.Context, token string) error {
	m.deletedToken = token
	return m.deleteErr
}

func newFCMRouter(h *Handler) *gin.Engine {
	r := gin.New()
	injectUID := func(c *gin.Context) {
		c.Set(middleware.FirebaseClaimsKey, &usecase.FirebaseToken{UID: "uid123"})
		c.Next()
	}
	r.POST("/api/v1/fcm/register", injectUID, h.RegisterFCMToken)
	r.DELETE("/api/v1/fcm/unregister", injectUID, h.UnregisterFCMToken)
	return r
}

func TestRegisterFCMToken_Success(t *testing.T) {
	repo := &mockFCMTokenRepo{}
	h := &Handler{fcmTokenRepo: repo}
	r := newFCMRouter(h)

	body, _ := json.Marshal(map[string]string{"token": "test-token", "platform": "android"})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/api/v1/fcm/register", bytes.NewReader(body)))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	if repo.savedToken != "test-token" || repo.savedPlatform != "android" || repo.savedUserID != "uid123" {
		t.Errorf("unexpected saved values: token=%s platform=%s uid=%s", repo.savedToken, repo.savedPlatform, repo.savedUserID)
	}
}

func TestRegisterFCMToken_InvalidPlatform(t *testing.T) {
	repo := &mockFCMTokenRepo{}
	h := &Handler{fcmTokenRepo: repo}
	r := newFCMRouter(h)

	body, _ := json.Marshal(map[string]string{"token": "test-token", "platform": "unknown"})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/api/v1/fcm/register", bytes.NewReader(body)))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestRegisterFCMToken_NoAuth(t *testing.T) {
	h := &Handler{fcmTokenRepo: &mockFCMTokenRepo{}}
	r := gin.New()
	r.POST("/api/v1/fcm/register", h.RegisterFCMToken)

	body, _ := json.Marshal(map[string]string{"token": "t", "platform": "web"})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/api/v1/fcm/register", bytes.NewReader(body)))

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestUnregisterFCMToken_Success(t *testing.T) {
	repo := &mockFCMTokenRepo{}
	h := &Handler{fcmTokenRepo: repo}
	r := newFCMRouter(h)

	body, _ := json.Marshal(map[string]string{"token": "del-token"})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/api/v1/fcm/unregister", bytes.NewReader(body)))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	if repo.deletedToken != "del-token" {
		t.Errorf("expected del-token to be deleted, got %q", repo.deletedToken)
	}
}

func TestUnregisterFCMToken_NoAuth(t *testing.T) {
	h := &Handler{fcmTokenRepo: &mockFCMTokenRepo{}}
	r := gin.New()
	r.DELETE("/api/v1/fcm/unregister", h.UnregisterFCMToken)

	body, _ := json.Marshal(map[string]string{"token": "t"})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/api/v1/fcm/unregister", bytes.NewReader(body)))

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
