package middleware

import (
	"encoding/json"
	"log"
	"time"

	"study-manager-service/internal/config"

	"github.com/gin-gonic/gin"
)

// AuditLogger representa o logger de auditoria
type AuditLogger struct {
	cfg *config.Config
}

// NewAuditLogger cria um novo logger de auditoria
func NewAuditLogger(cfg *config.Config) *AuditLogger {
	return &AuditLogger{cfg: cfg}
}

// AuditSensitiveOperations registra operações sensíveis
func (al *AuditLogger) AuditSensitiveOperations() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !al.cfg.Security.AuditEnabled {
			c.Next()
			return
		}

		start := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path

		// Verificar se é uma operação sensível
		if al.isSensitiveOperation(method) {
			// Executar a requisição
			c.Next()

			// Registrar após a execução
			duration := time.Since(start)
			al.logSensitiveOperation(c, method, path, duration)
		} else {
			c.Next()
		}
	}
}

// isSensitiveOperation verifica se é uma operação sensível
func (al *AuditLogger) isSensitiveOperation(method string) bool {
	sensitiveMethods := []string{"POST", "PUT", "PATCH", "DELETE"}
	for _, m := range sensitiveMethods {
		if method == m {
			return true
		}
	}
	return false
}

// logSensitiveOperation registra uma operação sensível
func (al *AuditLogger) logSensitiveOperation(c *gin.Context, method, path string, duration time.Duration) {
	userID := c.GetString("user_id")
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	statusCode := c.Writer.Status()

	auditLog := map[string]interface{}{
		"timestamp":    time.Now().UTC().Format(time.RFC3339),
		"event_type":   "SENSITIVE_OPERATION",
		"method":       method,
		"path":         path,
		"user_id":      userID,
		"client_ip":    clientIP,
		"user_agent":   userAgent,
		"status_code":  statusCode,
		"duration_ms":  duration.Milliseconds(),
		"service":      "study-manager-service",
	}

	logJSON, _ := json.Marshal(auditLog)
	log.Printf("AUDIT: %s", string(logJSON))
}
