package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/internal/transport/middleware"
	"backend/internal/usecase"
)

// @Summary		Get current user
// @Description	Returns the authenticated Firebase user's decoded claims. Requires Authorization: Bearer <firebase-id-token>.
// @Tags		auth
// @Produce		json
// @Security	BearerAuth
// @Success		200	{object}	FirebaseToken
// @Failure		401	{object}	object{error=string}	"Missing or invalid token"
// @Router		/api/v1/me [get]
func (h *Handler) MeHandler(c *gin.Context) {
	val, _ := c.Get(middleware.FirebaseClaimsKey)
	token, ok := val.(*usecase.FirebaseToken)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.JSON(http.StatusOK, token)
}
