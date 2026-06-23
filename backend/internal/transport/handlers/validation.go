package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// bindJSON binds and validates the request body. Writes 400 on failure and returns false.
func bindJSON(c *gin.Context, dst any) bool {
	if err := c.ShouldBindJSON(dst); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}
	return true
}

// bindQuery binds and validates query params. Writes 400 on failure and returns false.
func bindQuery(c *gin.Context, dst any) bool {
	if err := c.ShouldBindQuery(dst); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}
	return true
}
