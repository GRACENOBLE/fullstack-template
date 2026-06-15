package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"backend/internal/bootstrap"
	"backend/internal/infrastructure/database/postgres"
	"backend/internal/transport/handlers"
	"backend/internal/usecase"
)

// NewServer wires all layers and returns a configured *http.Server.
func NewServer(app *bootstrap.App) *http.Server {
	switch app.Config.Env {
	case "staging", "production":
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	// Suppress all [GIN-debug] output (banner + route table); docs URL is logged in main instead.
	gin.DebugPrintFunc = func(_ string, _ ...any) {}
	gin.DebugPrintRouteFunc = func(_, _, _ string, _ int) {}

	healthRepo := postgres.NewHealthRepository(app.DB)
	healthUC := usecase.NewHealthUseCase(healthRepo)
	h := handlers.NewHandler(healthUC)

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.Port),
		Handler:      h.RegisterRoutes(app.Config.RateLimitRPS, app.Config.RateLimitBurst, app.Firebase),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
