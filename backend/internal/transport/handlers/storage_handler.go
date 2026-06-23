package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type presignRequest struct {
	Filename    string `json:"filename"     binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
}

type presignResponse struct {
	UploadURL string `json:"upload_url"`
	PublicURL string `json:"public_url"`
}

// PresignHandler godoc
// @Summary     Request a presigned upload URL
// @Description Returns a presigned PUT URL and the final public URL. The client uploads directly to R2 using the presigned URL.
// @Tags        storage
// @Accept      json
// @Produce     json
// @Param       body body presignRequest true "Upload request"
// @Success     200 {object} presignResponse
// @Failure     400 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /api/v1/storage/presign [post]
// @Security    BearerAuth
func (h *Handler) PresignHandler(c *gin.Context) {
	var req presignRequest
	if !bindJSON(c, &req) {
		return
	}

	uploadURL, err := h.storageService.PresignUpload(c.Request.Context(), req.Filename, req.ContentType, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate upload URL"})
		return
	}

	c.JSON(http.StatusOK, presignResponse{
		UploadURL: uploadURL,
		PublicURL: h.storageService.PublicURL(req.Filename),
	})
}

// DeleteObjectHandler godoc
// @Summary     Delete a stored object
// @Tags        storage
// @Param       key path string true "Object key"
// @Success     204
// @Failure     500 {object} map[string]string
// @Router      /api/v1/storage/{key} [delete]
// @Security    BearerAuth
func (h *Handler) DeleteObjectHandler(c *gin.Context) {
	key := c.Param("key")

	if err := h.storageService.Delete(c.Request.Context(), key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete object"})
		return
	}

	c.Status(http.StatusNoContent)
}
