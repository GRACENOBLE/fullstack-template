package middleware_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"backend/internal/transport/middleware"
	"backend/internal/usecase"
)

type mockVerifier struct {
	token *usecase.FirebaseToken
	err   error
}

func (m *mockVerifier) VerifyIDToken(_ context.Context, _ string) (*usecase.FirebaseToken, error) {
	return m.token, m.err
}

func TestFirebaseAuth_MissingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.FirebaseAuth(&mockVerifier{}))
	r.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestFirebaseAuth_NonBearerHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.FirebaseAuth(&mockVerifier{}))
	r.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Basic dXNlcjpwYXNz")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestFirebaseAuth_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mv := &mockVerifier{err: errors.New("token expired")}
	r := gin.New()
	r.Use(middleware.FirebaseAuth(mv))
	r.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer badtoken")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestFirebaseAuth_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	want := &usecase.FirebaseToken{UID: "uid123", Email: "test@example.com"}
	mv := &mockVerifier{token: want}

	var capturedClaims *usecase.FirebaseToken
	r := gin.New()
	r.Use(middleware.FirebaseAuth(mv))
	r.GET("/", func(c *gin.Context) {
		val, exists := c.Get(middleware.FirebaseClaimsKey)
		if !exists {
			t.Error("firebase_claims not set in context")
			c.Status(http.StatusInternalServerError)
			return
		}
		tok, ok := val.(*usecase.FirebaseToken)
		if !ok {
			t.Errorf("claims wrong type: %T", val)
			c.Status(http.StatusInternalServerError)
			return
		}
		capturedClaims = tok
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer validtoken")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if capturedClaims == nil || capturedClaims.UID != want.UID || capturedClaims.Email != want.Email {
		t.Errorf("claims mismatch: got %+v", capturedClaims)
	}
}
