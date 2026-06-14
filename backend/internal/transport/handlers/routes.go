package handlers

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "backend/docs/swagger"
	"backend/internal/transport/middleware"
)

// RegisterRoutes creates the Gin engine, applies middleware, and registers all routes.
func (h *Handler) RegisterRoutes() http.Handler {
	r := gin.New()

	// Use Gin's colorful logger locally; structured slog logger in staging/production.
	if gin.Mode() == gin.DebugMode {
		r.Use(gin.Recovery(), gin.Logger())
	} else {
		r.Use(gin.Recovery(), middleware.Logger())
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.GET("/", h.HelloWorldHandler)
	r.GET("/health", h.HealthHandler)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
