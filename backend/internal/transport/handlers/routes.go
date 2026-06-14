package handlers

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"backend/internal/transport/middleware"
)

// RegisterRoutes creates the Gin engine, applies middleware, and registers all routes.
func (h *Handler) RegisterRoutes() http.Handler {
	r := gin.New()
	r.Use(gin.Recovery(), middleware.Logger())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.GET("/", h.HelloWorldHandler)
	r.GET("/health", h.HealthHandler)

	return r
}
