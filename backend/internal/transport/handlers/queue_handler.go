package handlers

import "github.com/gin-gonic/gin"

// QueueUIHandler serves the Asynqmon job monitoring UI.
// Only registered in debug mode when REDIS_URL is set.
//
//	@Summary		Asynq job monitoring UI
//	@Tags			observability
//	@Produce		html
//	@Success		200	{string}	string	"Asynqmon UI"
//	@Router			/admin/queues [get]
func (h *Handler) QueueUIHandler(_ *gin.Context) {}
