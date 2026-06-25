package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/internal/infrastructure/streams"
	"backend/internal/transport/middleware"
	"backend/internal/usecase"
)

type registerFCMTokenRequest struct {
	Token    string `json:"token"    binding:"required"`
	Platform string `json:"platform" binding:"required,oneof=android ios web"`
}

type unregisterFCMTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// RegisterFCMToken saves a device FCM registration token for the authenticated user.
//
// @Summary     Register FCM token
// @Description Saves a device FCM registration token for the authenticated user. If the token already exists it is upserted.
// @Tags        fcm
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body body registerFCMTokenRequest true "Token registration payload"
// @Success     200  {object} object{message=string}
// @Failure     400  {object} object{error=string}
// @Failure     401  {object} object{error=string}
// @Failure     500  {object} object{error=string}
// @Router      /api/v1/fcm/register [post]
func (h *Handler) RegisterFCMToken(c *gin.Context) {
	val, _ := c.Get(middleware.FirebaseClaimsKey)
	claims, ok := val.(*usecase.FirebaseToken)
	if !ok {
		JSONError(c, http.StatusUnauthorized, "unauthorized", "missing or invalid token")
		return
	}

	var req registerFCMTokenRequest
	if !bindJSON(c, &req) {
		return
	}

	if err := h.fcmTokenRepo.SaveToken(c.Request.Context(), claims.UID, req.Token, req.Platform); err != nil {
		JSONError(c, http.StatusInternalServerError, "internal_error", "failed to save token")
		return
	}

	if h.streamProducer != nil {
		if err := h.streamProducer.Publish(c.Request.Context(), streams.StreamUserCreated, streams.UserCreatedEvent{
			UserID: claims.UID,
			Email:  claims.Email,
			Name:   claims.Name,
		}); err != nil {
			slog.Warn("fcm: failed to publish user created event", "user_id", claims.UID, "err", err)
		}
	}

	JSON(c, gin.H{"message": "token registered"})
}

// UnregisterFCMToken removes an FCM device token, typically called on logout.
//
// @Summary     Unregister FCM token
// @Description Removes a device FCM registration token. Typically called on logout or device change.
// @Tags        fcm
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body body unregisterFCMTokenRequest true "Token to remove"
// @Success     200  {object} object{message=string}
// @Failure     400  {object} object{error=string}
// @Failure     401  {object} object{error=string}
// @Failure     500  {object} object{error=string}
// @Router      /api/v1/fcm/unregister [delete]
func (h *Handler) UnregisterFCMToken(c *gin.Context) {
	val, _ := c.Get(middleware.FirebaseClaimsKey)
	claims, ok := val.(*usecase.FirebaseToken)
	if !ok {
		JSONError(c, http.StatusUnauthorized, "unauthorized", "missing or invalid token")
		return
	}

	var req unregisterFCMTokenRequest
	if !bindJSON(c, &req) {
		return
	}

	if err := h.fcmTokenRepo.DeleteToken(c.Request.Context(), claims.UID, req.Token); err != nil {
		JSONError(c, http.StatusInternalServerError, "internal_error", "failed to remove token")
		return
	}

	JSON(c, gin.H{"message": "token unregistered"})
}
