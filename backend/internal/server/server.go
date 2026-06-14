package server

import (
	"fmt"
	"net/http"
	"time"

	"backend/internal/bootstrap"
	"backend/internal/handler"
	"backend/internal/repository/postgres"
	"backend/internal/usecase"
)

// NewServer wires all layers and returns a configured *http.Server.
func NewServer(app *bootstrap.App) *http.Server {
	healthRepo := postgres.NewHealthRepository(app.DB)
	healthUC := usecase.NewHealthUseCase(healthRepo)
	h := handler.NewHandler(healthUC)

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.Port),
		Handler:      h.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
