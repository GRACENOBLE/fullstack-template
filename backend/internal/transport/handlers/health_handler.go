package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HealthHandler(c *gin.Context) {
	stats, err := h.healthUC.GetHealth(c.Request.Context())
	if err != nil {
		slog.Warn("health check failed", "error", err)
		c.JSON(http.StatusServiceUnavailable, stats)
		return
	}
	c.JSON(http.StatusOK, stats)
}
