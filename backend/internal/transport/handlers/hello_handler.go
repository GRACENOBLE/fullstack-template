package handlers

import "github.com/gin-gonic/gin"

//	@Summary	Hello World
//	@Tags		general
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Router		/ [get]
func (h *Handler) HelloWorldHandler(c *gin.Context) {
	JSON(c, gin.H{"message": "Hello World"})
}
