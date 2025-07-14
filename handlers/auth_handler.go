package handlers

import (
	"net/http"

	"auth-service/models"
	"auth-service/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Registrar novo usuário
// @Description Cria uma nova conta de usuário
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "Dados do usuário"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos: " + err.Error()})
		return
	}

	user, err := h.authService.Register(&req)
	if err != nil {
		if err.Error() == "email já está em uso" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary Fazer login
// @Description Autentica um usuário e retorna tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Credenciais de login"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos: " + err.Error()})
		return
	}

	tokens, err := h.authService.Login(&req)
	if err != nil {
		if err.Error() == "credenciais inválidas" || err.Error() == "usuário inativo" || err.Error() == "cliente não encontrado" || err.Error() == "cliente inativo" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// RefreshToken godoc
// @Summary Renovar token
// @Description Renova o access token usando um refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh body models.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos: " + err.Error()})
		return
	}

	tokens, err := h.authService.RefreshToken(&req)
	if err != nil {
		if err.Error() == "refresh token inválido" || err.Error() == "refresh token revogado" || err.Error() == "refresh token expirado" || err.Error() == "usuário inativo" || err.Error() == "cliente não encontrado" || err.Error() == "cliente inativo" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// ValidateToken godoc
// @Summary Validar token
// @Description Valida um access token e retorna informações do usuário
// @Tags auth
// @Accept json
// @Produce json
// @Param validation body models.ValidateTokenRequest true "Token para validação"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/validate [post]
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	var req models.ValidateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos: " + err.Error()})
		return
	}

	user, err := h.authService.ValidateToken(req.Token, req.ClientID)
	if err != nil {
		if err.Error() == "token inválido" || err.Error() == "cliente não encontrado" || err.Error() == "cliente inativo" || err.Error() == "usuário inativo" || err.Error() == "tipo de token inválido" || err.Error() == "cliente não autorizado" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateClient godoc
// @Summary Criar cliente
// @Description Cria um novo cliente para autenticação
// @Tags clients
// @Accept json
// @Produce json
// @Param client body map[string]string true "Dados do cliente"
// @Success 201 {object} models.Client
// @Failure 400 {object} map[string]interface{}
// @Router /clients [post]
func (h *AuthHandler) CreateClient(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos: " + err.Error()})
		return
	}

	client, err := h.authService.CreateClient(req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, client)
}

// GetProfile godoc
// @Summary Obter perfil do usuário
// @Description Retorna informações do usuário autenticado
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} map[string]interface{}
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não autenticado"})
		return
	}

	c.JSON(http.StatusOK, user)
}
