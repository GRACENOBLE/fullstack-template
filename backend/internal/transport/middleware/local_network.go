package middleware

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LocalNetworkOnly rejects requests that originate outside loopback or RFC 1918
// private address space. Use this to restrict internal-only endpoints (e.g.
// /metrics) from being reachable by external clients in production.
func LocalNetworkOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := net.ParseIP(c.ClientIP())
		if ip == nil || (!ip.IsLoopback() && !ip.IsPrivate()) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}
