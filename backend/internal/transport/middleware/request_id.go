package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

const RequestIDKey = "request_id"
const RequestIDHeader = "X-Request-ID"

// RequestID reads X-Request-ID from the incoming request. If absent or empty,
// it generates a random 16-byte hex ID. The ID is stored on the Gin context
// under RequestIDKey and echoed back in the X-Request-ID response header.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(RequestIDHeader)
		if id == "" {
			b := make([]byte, 16)
			_, _ = rand.Read(b)
			id = hex.EncodeToString(b)
		}
		c.Set(RequestIDKey, id)
		c.Header(RequestIDHeader, id)
		c.Next()
	}
}
