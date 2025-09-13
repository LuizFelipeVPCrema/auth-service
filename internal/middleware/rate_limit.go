package middleware

import (
	"net/http"
	"sync"
	"time"

	"study-manager-service/internal/config"

	"github.com/gin-gonic/gin"
)

// RateLimiter representa o rate limiter
type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
}

// NewRateLimiter cria um novo rate limiter
func NewRateLimiter(cfg *config.Config) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    cfg.Security.RateLimitRequests,
		window:   time.Duration(cfg.Security.RateLimitWindowMinutes) * time.Minute,
	}
}

// RateLimit middleware que implementa rate limiting
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := time.Now()

		rl.mutex.Lock()
		defer rl.mutex.Unlock()

		// Limpar requisições antigas
		cutoff := now.Add(-rl.window)
		validRequests := []time.Time{}
		for _, reqTime := range rl.requests[clientIP] {
			if reqTime.After(cutoff) {
				validRequests = append(validRequests, reqTime)
			}
		}
		rl.requests[clientIP] = validRequests

		// Verificar se excedeu o limite
		if len(rl.requests[clientIP]) >= rl.limit {
			resetTime := rl.requests[clientIP][0].Add(rl.window)
			c.Header("X-RateLimit-Limit", string(rune(rl.limit)))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", resetTime.Format(time.RFC3339))
			c.Header("Retry-After", string(rune(int(rl.window.Seconds()))))

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "limite de requisições excedido",
				"code":    "RATE_LIMIT_EXCEEDED",
				"message": "Muitas requisições. Tente novamente mais tarde.",
				"retry_after": int(rl.window.Seconds()),
			})
			c.Abort()
			return
		}

		// Adicionar requisição atual
		rl.requests[clientIP] = append(rl.requests[clientIP], now)

		// Adicionar headers de rate limit
		remaining := rl.limit - len(rl.requests[clientIP])
		c.Header("X-RateLimit-Limit", string(rune(rl.limit)))
		c.Header("X-RateLimit-Remaining", string(rune(remaining)))
		c.Header("X-RateLimit-Reset", rl.requests[clientIP][0].Add(rl.window).Format(time.RFC3339))

		c.Next()
	}
}
