package handlers

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "backend/docs/swagger"
	"backend/internal/transport/middleware"
	"backend/internal/usecase"
)

// RegisterRoutes creates the Gin engine, applies middleware, and registers all routes.
// rps and burst configure IP-based rate limiting; pass rps<=0 to disable.
// verifier enables Firebase token auth on protected routes; pass nil to skip auth (dev only).
// sentryDSN enables Sentry error tracking; pass empty string to disable.
func (h *Handler) RegisterRoutes(rps float64, burst int, verifier usecase.FirebaseTokenVerifier, sentryDSN string) http.Handler {
	r := gin.New()

	r.Use(middleware.SentryMiddleware(sentryDSN))

	// Use Gin's colorful logger locally; structured slog logger in staging/production.
	if gin.Mode() == gin.DebugMode {
		r.Use(gin.Recovery(), gin.Logger())
	} else {
		r.Use(gin.Recovery(), middleware.Logger())
	}

	r.Use(middleware.RateLimit(rps, burst))

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.GET("/", h.HelloWorldHandler)
	r.GET("/health", h.HealthHandler)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	if verifier != nil {
		api.Use(middleware.FirebaseAuth(verifier))
	}
	api.GET("/me", h.MeHandler)

	return r
}
