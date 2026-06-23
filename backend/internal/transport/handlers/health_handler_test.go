package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"backend/internal/domain"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// mockHealthUC is a local test double implementing usecase.HealthUseCase.
type mockHealthUC struct {
	stats domain.HealthStats
	err   error
}

func (m *mockHealthUC) GetHealth(_ context.Context) (domain.HealthStats, error) {
	return m.stats, m.err
}

func TestHealthHandler_Success(t *testing.T) {
	want := domain.HealthStats{
		Status:  "up",
		Message: "It's healthy",
	}
	h := NewHandler(&mockHealthUC{stats: want}, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	r := gin.New()
	r.GET("/health", h.HealthHandler)

	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestHealthHandler_ServiceUnavailable(t *testing.T) {
	h := NewHandler(&mockHealthUC{err: errors.New("connection refused")}, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	r := gin.New()
	r.GET("/health", h.HealthHandler)

	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status 503, got %d", rr.Code)
	}
}
