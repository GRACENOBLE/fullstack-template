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

// mockUserRepo satisfies usecase.UserRepository for handler tests.
type mockUserRepo struct {
	upserted  *domain.User
	deleted   string
	upsertErr error
	deleteErr error
}

func (m *mockUserRepo) Upsert(_ context.Context, u *domain.User) (*domain.User, error) {
	if m.upsertErr != nil {
		return nil, m.upsertErr
	}
	out := *u
	out.ID = 1
	m.upserted = &out
	return &out, nil
}

func (m *mockUserRepo) DeleteByFirebaseUID(_ context.Context, firebaseUID string) error {
	m.deleted = firebaseUID
	return m.deleteErr
}

func newMeRouter(h *Handler) *gin.Engine {
	r := gin.New()
	injectClaims := func(c *gin.Context) {
		c.Set(middleware.FirebaseClaimsKey, &usecase.FirebaseToken{
			UID:      "uid-test",
			Email:    "test@example.com",
			PhotoURL: "https://example.com/photo.png",
		})
		c.Next()
	}
	r.PATCH("/api/v1/me", injectClaims, h.UpdateMeHandler)
	r.DELETE("/api/v1/me", injectClaims, h.DeleteMeHandler)
	return r
}

func TestUpdateMeHandler_Success(t *testing.T) {
	repo := &mockUserRepo{}
	h := &Handler{userRepo: repo}
	r := newMeRouter(h)

	body, _ := json.Marshal(map[string]string{"name": "Alice"})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPatch, "/api/v1/me", bytes.NewReader(body)))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	if repo.upserted == nil {
		t.Fatal("expected Upsert to be called")
	}
	if repo.upserted.Name != "Alice" {
		t.Errorf("Name: got %q, want %q", repo.upserted.Name, "Alice")
	}
	if repo.upserted.FirebaseUID != "uid-test" {
		t.Errorf("FirebaseUID: got %q, want %q", repo.upserted.FirebaseUID, "uid-test")
	}
}

func TestUpdateMeHandler_MissingBody(t *testing.T) {
	repo := &mockUserRepo{}
	h := &Handler{userRepo: repo}
	r := newMeRouter(h)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPatch, "/api/v1/me", bytes.NewReader([]byte(`{}`)))) // missing required name

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdateMeHandler_NoClaims(t *testing.T) {
	h := &Handler{userRepo: &mockUserRepo{}}
	r := gin.New()
	r.PATCH("/api/v1/me", h.UpdateMeHandler)

	body, _ := json.Marshal(map[string]string{"name": "Alice"})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPatch, "/api/v1/me", bytes.NewReader(body)))

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestDeleteMeHandler_Success(t *testing.T) {
	repo := &mockUserRepo{}
	h := &Handler{userRepo: repo}
	r := newMeRouter(h)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/api/v1/me", nil))

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", w.Code, w.Body.String())
	}
	if repo.deleted != "uid-test" {
		t.Errorf("expected DeleteByFirebaseUID(uid-test), got %q", repo.deleted)
	}
}

func TestDeleteMeHandler_NoClaims(t *testing.T) {
	h := &Handler{userRepo: &mockUserRepo{}}
	r := gin.New()
	r.DELETE("/api/v1/me", h.DeleteMeHandler)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/api/v1/me", nil))

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
