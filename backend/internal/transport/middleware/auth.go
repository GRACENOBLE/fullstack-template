package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"backend/internal/usecase"
)

// FirebaseClaimsKey is the Gin context key under which the verified FirebaseToken is stored.
const FirebaseClaimsKey = "firebase_claims"

// FirebaseAuth returns a Gin middleware that validates a Firebase ID token from the
// Authorization: Bearer <token> header. Sets FirebaseClaimsKey in the Gin context on success.
// Returns 401 when the header is missing, malformed, or the token fails verification.
func FirebaseAuth(verifier usecase.FirebaseTokenVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "missing or invalid Authorization header"}})
			return
		}
		idToken := strings.TrimPrefix(header, "Bearer ")

		claims, err := verifier.VerifyIDToken(c.Request.Context(), idToken)
		if err != nil || claims == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "invalid or expired token"}})
			return
		}

		c.Set(FirebaseClaimsKey, claims)
		c.Next()
	}
}
