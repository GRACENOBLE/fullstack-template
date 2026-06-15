package handlers

import "github.com/gin-gonic/gin"

// MetricsHandler serves Prometheus metrics.
// Actual scraping logic is delegated to promhttp.Handler() registered in routes.go.
//
//	@Summary		Prometheus metrics
//	@Tags			observability
//	@Produce		plain
//	@Success		200	{string}	string	"Prometheus exposition format"
//	@Router			/metrics [get]
func (h *Handler) MetricsHandler(_ *gin.Context) {}
