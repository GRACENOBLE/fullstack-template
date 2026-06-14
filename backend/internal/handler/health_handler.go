package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HealthHandler(c *gin.Context) {
	stats, err := h.healthUC.GetHealth(c.Request.Context())
	if err != nil {
		log.Printf("health check failed: %v", err)
		c.JSON(http.StatusServiceUnavailable, stats)
		return
	}
	c.JSON(http.StatusOK, stats)
}
