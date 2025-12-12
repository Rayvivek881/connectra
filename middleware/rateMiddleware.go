package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	tokens   int
	maxLimit int
	mu       sync.Mutex
	once     sync.Once
)

func startTokenRefiller() {
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			mu.Lock()
			if tokens < maxLimit {
				tokens++
			}
			mu.Unlock()
		}
	}()
}

func RateLimiter() gin.HandlerFunc {
	once.Do(func() {
		maxLimit = 60
		tokens = maxLimit
		startTokenRefiller()
	})

	return func(c *gin.Context) {
		mu.Lock()
		if tokens <= 0 {
			mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate limit exceeded",
				"message": "too many requests, please try again later",
			})
			return
		}
		tokens--
		mu.Unlock()
		c.Next()
	}
}
