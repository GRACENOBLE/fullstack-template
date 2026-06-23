package middleware

import (
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"backend/internal/usecase"
)

// GeoLocationKey is the Gin context key under which *domain.GeoLocation is stored.
const GeoLocationKey = "geo_location"

// GeoFromRequest is a best-effort middleware that resolves the request IP to
// geographic metadata and stores it in the Gin context under GeoLocationKey.
// If geolocation fails (private IP, rate-limit, network error), the request
// continues without geo data — handlers must nil-check before reading the key.
func GeoFromRequest(locator usecase.GeoLocator) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := RealIP(c.Request)
		if geo, err := locator.Lookup(c.Request.Context(), ip); err == nil {
			c.Set(GeoLocationKey, geo)
		}
		c.Next()
	}
}

// RealIP extracts the originating IP from the request, respecting
// X-Forwarded-For (Railway/proxy) and X-Real-IP headers.
// Exported so tests and other packages can call it directly.
func RealIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.TrimSpace(strings.SplitN(xff, ",", 2)[0])
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
