package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	mu       sync.Mutex
	attempts map[string][]time.Time
	max      int
	window   time.Duration
}

var limiter = &rateLimiter{
	attempts: make(map[string][]time.Time),
	max:      5,
	window:   time.Minute,
}

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		limiter.mu.Lock()
		now := time.Now()

		// Clean old entries
		valid := make([]time.Time, 0)
		for _, t := range limiter.attempts[ip] {
			if now.Sub(t) < limiter.window {
				valid = append(valid, t)
			}
		}
		limiter.attempts[ip] = valid

		if len(valid) >= limiter.max {
			limiter.mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many attempts, try again later"})
			c.Abort()
			return
		}

		limiter.attempts[ip] = append(limiter.attempts[ip], now)
		limiter.mu.Unlock()

		c.Next()
	}
}
