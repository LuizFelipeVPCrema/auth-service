package middleware

import (
	"net/http"
	"strings"

	"study-manager-service/internal/config"

	"github.com/gin-gonic/gin"
)

// ValidationMiddleware representa o middleware de validação
type ValidationMiddleware struct {
	cfg *config.Config
}

// NewValidationMiddleware cria um novo middleware de validação
func NewValidationMiddleware(cfg *config.Config) *ValidationMiddleware {
	return &ValidationMiddleware{cfg: cfg}
}

// ValidateInput valida a entrada da requisição
func (vm *ValidationMiddleware) ValidateInput() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar tamanho da requisição
		if c.Request.ContentLength > vm.cfg.Security.MaxRequestSize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error":   "requisição muito grande",
				"code":    "REQUEST_TOO_LARGE",
				"message": "Tamanho da requisição excede o limite permitido",
			})
			c.Abort()
			return
		}

		// Verificar headers suspeitos
		if vm.hasSuspiciousHeaders(c) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "headers suspeitos detectados",
				"code":    "SUSPICIOUS_HEADERS",
				"message": "Headers maliciosos detectados na requisição",
			})
			c.Abort()
			return
		}

		// Verificar user agent suspeito
		if vm.hasSuspiciousUserAgent(c) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "user agent suspeito",
				"code":    "SUSPICIOUS_USER_AGENT",
				"message": "User agent malicioso detectado",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// hasSuspiciousHeaders verifica se há headers suspeitos
func (vm *ValidationMiddleware) hasSuspiciousHeaders(c *gin.Context) bool {
	suspiciousHeaders := []string{
		"X-Forwarded-For",
		"X-Real-IP",
		"X-Forwarded-Host",
		"X-Forwarded-Proto",
		"X-Original-URL",
		"X-Rewrite-URL",
	}

	for _, header := range suspiciousHeaders {
		if c.GetHeader(header) != "" {
			return true
		}
	}

	return false
}

// hasSuspiciousUserAgent verifica se o user agent é suspeito
func (vm *ValidationMiddleware) hasSuspiciousUserAgent(c *gin.Context) bool {
	userAgent := strings.ToLower(c.GetHeader("User-Agent"))
	
	suspiciousAgents := []string{
		"sqlmap",
		"nikto",
		"nmap",
		"masscan",
		"zap",
		"burp",
		"w3af",
		"havij",
		"acunetix",
		"nessus",
		"openvas",
		"wpscan",
		"dirb",
		"gobuster",
		"wfuzz",
		"ffuf",
	}

	for _, agent := range suspiciousAgents {
		if strings.Contains(userAgent, agent) {
			return true
		}
	}

	return false
}
