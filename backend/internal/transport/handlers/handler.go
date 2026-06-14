package handlers

import (
	"backend/internal/usecase"
)

// Handler holds all use case dependencies for HTTP handlers.
type Handler struct {
	healthUC usecase.HealthUseCase
}

// NewHandler constructs a Handler with all required use cases.
func NewHandler(healthUC usecase.HealthUseCase) *Handler {
	return &Handler{healthUC: healthUC}
}
