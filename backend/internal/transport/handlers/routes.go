package handlers

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "backend/docs/swagger"
	"backend/internal/transport/middleware"
)

// RegisterRoutes creates the Gin engine, applies middleware, and registers all routes.
// rps and burst configure IP-based rate limiting; pass rps<=0 to disable.
// sentryDSN enables Sentry error tracking; pass empty string to disable.
// Firebase auth is read from h.verifier; nil disables auth (dev only).
func (h *Handler) RegisterRoutes(rps float64, burst int, sentryDSN string) http.Handler {
	r := gin.New()

	r.Use(middleware.SentryMiddleware(sentryDSN))

	// Use Gin's colorful logger locally; structured slog logger in staging/production.
	if gin.Mode() == gin.DebugMode {
		r.Use(gin.Recovery(), gin.Logger())
	} else {
		r.Use(gin.Recovery(), middleware.Logger())
	}

	r.Use(middleware.PrometheusMiddleware())
	r.Use(middleware.RateLimit(rps, burst))

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.GET("/", h.HelloWorldHandler)
	r.GET("/health", h.HealthHandler)
	r.GET("/ws", h.WsHandler)

	// In staging/production, restrict /metrics to loopback and RFC 1918 addresses
	// so Prometheus can scrape from the internal network but external clients cannot.
	if gin.Mode() == gin.ReleaseMode {
		r.GET("/metrics", middleware.LocalNetworkOnly(), gin.WrapH(promhttp.Handler()))
	} else {
		r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Asynqmon job-monitoring UI — debug/local only.
	if gin.Mode() == gin.DebugMode && h.queueUI != nil {
		r.GET("/admin/queues", gin.WrapH(h.queueUI))
		r.Any("/admin/queues/*path", gin.WrapH(h.queueUI))
	}

	api := r.Group("/api/v1")
	if h.verifier != nil {
		api.Use(middleware.FirebaseAuth(h.verifier))
	}
	api.GET("/me", h.MeHandler)

	if h.fcmTokenRepo != nil {
		api.POST("/fcm/register", h.RegisterFCMToken)
		api.DELETE("/fcm/unregister", h.UnregisterFCMToken)
	}

	return r
}
