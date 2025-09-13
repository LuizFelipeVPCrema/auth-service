package middleware

import (
	"net/http"
	"strings"

	"study-manager-service/internal/clients"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware representa o middleware de autenticação
type AuthMiddleware struct {
	authClient *clients.AuthClient
}

// NewAuthMiddleware cria um novo middleware de autenticação
func NewAuthMiddleware(authClient *clients.AuthClient) *AuthMiddleware {
	return &AuthMiddleware{
		authClient: authClient,
	}
}

// RequireAuth middleware que exige autenticação
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extrair token do header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "token de autorização não fornecido",
				"code":    "MISSING_AUTH_TOKEN",
				"message": "Header Authorization é obrigatório",
			})
			c.Abort()
			return
		}

		// Verificar formato "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "formato de token inválido",
				"code":    "INVALID_TOKEN_FORMAT",
				"message": "Token deve estar no formato 'Bearer <token>'",
			})
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
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "client_id não fornecido",
				"code":    "MISSING_CLIENT_ID",
				"message": "Header X-Client-ID ou parâmetro client_id é obrigatório",
			})
			c.Abort()
			return
		}

		// Validar token com o auth-service
		user, err := m.authClient.ValidateToken(token, clientID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "token inválido",
				"code":    "INVALID_TOKEN",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		// Verificar se o usuário está ativo
		if !user.Active {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "usuário inativo",
				"code":    "USER_INACTIVE",
				"message": "Usuário não está ativo no sistema",
			})
			c.Abort()
			return
		}

		// Adicionar informações do usuário ao contexto
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("client_id", clientID)

		c.Next()
	}
}

// RequireClientID middleware que exige client_id
func (m *AuthMiddleware) RequireClientID() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := c.GetHeader("X-Client-ID")
		if clientID == "" {
			clientID = c.Query("client_id")
		}

		if clientID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "client_id não fornecido",
				"code":    "MISSING_CLIENT_ID",
				"message": "Header X-Client-ID ou parâmetro client_id é obrigatório",
			})
			c.Abort()
			return
		}

		c.Set("client_id", clientID)
		c.Next()
	}
}
