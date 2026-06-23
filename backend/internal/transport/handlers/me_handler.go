package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/internal/domain"
	"backend/internal/transport/middleware"
	"backend/internal/usecase"
)

type updateMeRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateMeHandler godoc
//
//	@Summary		Update current user profile
//	@Description	Upserts the authenticated user's profile record. Requires a valid Firebase ID token.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			body	body		updateMeRequest	true	"Profile update"
//	@Success		200		{object}	UserAlias
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/api/v1/me [patch]
//	@Security		BearerAuth
func (h *Handler) UpdateMeHandler(c *gin.Context) {
	var req updateMeRequest
	if !bindJSON(c, &req) {
		return
	}

	raw, exists := c.Get(middleware.FirebaseClaimsKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	claims, ok := raw.(*usecase.FirebaseToken)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	u := &domain.User{
		FirebaseUID: claims.UID,
		Name:        req.Name,
		Email:       claims.Email,
		PhotoURL:    claims.PhotoURL,
	}
	updated, err := h.userRepo.Upsert(c.Request.Context(), u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// DeleteMeHandler godoc
//
//	@Summary		Delete current user account
//	@Description	Removes the authenticated user's profile record from the database. The Firebase account is not deleted.
//	@Tags			users
//	@Produce		json
//	@Success		204
//	@Failure		401	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/api/v1/me [delete]
//	@Security		BearerAuth
func (h *Handler) DeleteMeHandler(c *gin.Context) {
	raw, exists := c.Get(middleware.FirebaseClaimsKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	claims, ok := raw.(*usecase.FirebaseToken)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.userRepo.DeleteByFirebaseUID(c.Request.Context(), claims.UID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete account"})
		return
	}
	c.Status(http.StatusNoContent)
}
