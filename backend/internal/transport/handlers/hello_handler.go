package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HelloWorldHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello World"})
}
