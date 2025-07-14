package middleware

import (
	"net/http"
	"strings"

	"auth-service/services"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService *services.AuthService
}

func NewAuthMiddleware(authService *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extrair token do header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token de autorização não fornecido"})
			c.Abort()
			return
		}

		// Verificar formato "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "formato de token inválido"})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// Extrair client_id do header ou query parameter
		clientID := c.GetHeader("X-Client-ID")
		if clientID == "" {
			clientID = c.Query("client_id")
		}

		if clientID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "client_id não fornecido"})
			c.Abort()
			return
		}

		// Validar token
		user, err := m.authService.ValidateToken(token, clientID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Adicionar informações do usuário ao contexto
		c.Set("user", user)
		c.Set("user_id", user.ID.String())
		c.Set("client_id", clientID)

		c.Next()
	}
}

func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extrair token do header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Verificar formato "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.Next()
			return
		}

		token := tokenParts[1]

		// Extrair client_id do header ou query parameter
		clientID := c.GetHeader("X-Client-ID")
		if clientID == "" {
			clientID = c.Query("client_id")
		}

		if clientID == "" {
			c.Next()
			return
		}

		// Tentar validar token (opcional)
		user, err := m.authService.ValidateToken(token, clientID)
		if err != nil {
			c.Next()
			return
		}

		// Adicionar informações do usuário ao contexto se válido
		c.Set("user", user)
		c.Set("user_id", user.ID.String())
		c.Set("client_id", clientID)

		c.Next()
	}
}
