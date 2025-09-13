package middleware

import (
	"strings"

	"study-manager-service/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS configura o middleware de CORS
func CORS(cfg *config.Config) gin.HandlerFunc {
	config := cors.Config{
		AllowOrigins:     cfg.Security.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With", "X-Client-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 horas
	}

	return cors.New(config)
}

// SecurityHeaders adiciona headers de segurança
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	}
}

// ValidateOrigin valida a origem da requisição
func ValidateOrigin(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			c.Next()
			return
		}

		allowed := false
		for _, allowedOrigin := range cfg.Security.AllowedOrigins {
			if strings.EqualFold(origin, allowedOrigin) {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(403, gin.H{
				"error":   "origem não permitida",
				"code":    "ORIGIN_NOT_ALLOWED",
				"message": "Origem da requisição não está na lista de origens permitidas",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
