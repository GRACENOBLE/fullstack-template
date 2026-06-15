package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

const (
	limiterTTL      = 10 * time.Minute
	cleanupInterval = 5 * time.Minute
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimit returns a Gin middleware that limits each client IP to rps requests per
// second with the given burst size. When rps <= 0 the middleware is a no-op.
// Burst is clamped to a minimum of 1 so fractional rps values never block all traffic.
// Stale per-IP entries are evicted every 5 minutes.
func RateLimit(rps float64, burst int) gin.HandlerFunc {
	if rps <= 0 {
		return func(c *gin.Context) { c.Next() }
	}
	if burst < 1 {
		burst = 1
	}

	var mu sync.Mutex
	limiters := make(map[string]*ipLimiter)

	go func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()
		for range ticker.C {
			cutoff := time.Now().Add(-limiterTTL)
			mu.Lock()
			for ip, l := range limiters {
				if l.lastSeen.Before(cutoff) {
					delete(limiters, ip)
				}
			}
			mu.Unlock()
		}
	}()

	getLimiter := func(ip string) *rate.Limiter {
		mu.Lock()
		defer mu.Unlock()
		l, ok := limiters[ip]
		if !ok {
			l = &ipLimiter{limiter: rate.NewLimiter(rate.Limit(rps), burst)}
			limiters[ip] = l
		}
		l.lastSeen = time.Now()
		return l.limiter
	}

	return func(c *gin.Context) {
		if !getLimiter(c.ClientIP()).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}
