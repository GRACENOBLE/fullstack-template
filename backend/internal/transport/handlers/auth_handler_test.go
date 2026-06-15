package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"backend/internal/transport/middleware"
	"backend/internal/usecase"
)

func TestMeHandler_WithClaims(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &Handler{}
	want := &usecase.FirebaseToken{UID: "abc123", Email: "user@example.com", Name: "Test User"}

	r := gin.New()
	r.GET("/me", func(c *gin.Context) {
		c.Set(middleware.FirebaseClaimsKey, want)
		h.MeHandler(c)
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/me", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var got usecase.FirebaseToken
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	if got.UID != want.UID || got.Email != want.Email || got.Name != want.Name {
		t.Errorf("response mismatch: got %+v, want %+v", got, *want)
	}
}

func TestMeHandler_WithoutClaims(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &Handler{}
	r := gin.New()
	r.GET("/me", h.MeHandler)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/me", nil))

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}
