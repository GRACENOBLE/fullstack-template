package middleware

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LocalNetworkOnly rejects requests that originate outside loopback or RFC 1918
// private address space. Use this to restrict internal-only endpoints (e.g.
// /metrics, /debug/pprof) from being reachable by external clients in production.
//
// Uses RemoteAddr (the TCP peer address) rather than c.ClientIP() to prevent
// X-Forwarded-For spoofing — an external caller cannot forge their RemoteAddr.
func LocalNetworkOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// RemoteAddr is "host:port"; strip the port before parsing.
		host, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		ip := net.ParseIP(host)
		if ip == nil || (!ip.IsLoopback() && !ip.IsPrivate()) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}
