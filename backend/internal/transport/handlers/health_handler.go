package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary	Health check
// @Tags		ops
// @Produce	json
// @Success	200	{object}	HealthStats
// @Failure	503	{object}	HealthStats
// @Router		/health [get]
func (h *Handler) HealthHandler(c *gin.Context) {
	stats, err := h.healthUC.GetHealth(c.Request.Context())
	if err != nil {
		slog.Warn("health check failed", "error", err)
		JSONStatus(c, http.StatusServiceUnavailable, stats)
		return
	}
	JSONStatus(c, http.StatusOK, stats)
}
