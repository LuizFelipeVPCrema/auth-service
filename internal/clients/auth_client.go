package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"study-manager-service/internal/config"
	"study-manager-service/internal/models"
)

// AuthClient representa o cliente para comunicação com o auth-service
type AuthClient struct {
	baseURL    string
	clientID   string
	httpClient *http.Client
}

// NewAuthClient cria um novo cliente de autenticação
func NewAuthClient(cfg *config.Config) *AuthClient {
	return &AuthClient{
		baseURL:  cfg.Auth.ServiceURL,
		clientID: cfg.Auth.ClientID,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.Auth.Timeout) * time.Second,
		},
	}
}

// ValidateToken valida um token JWT com o auth-service
func (c *AuthClient) ValidateToken(token, clientID string) (*models.AuthUser, error) {
	// Usar o clientID fornecido ou o configurado
	if clientID == "" {
		clientID = c.clientID
	}

	// Preparar requisição
	reqBody := models.ValidateTokenRequest{
		Token:    token,
		ClientID: clientID,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar requisição: %w", err)
	}

	// Fazer requisição para o auth-service
	url := fmt.Sprintf("%s/api/v1/validate", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Executar requisição
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar requisição: %w", err)
	}
	defer resp.Body.Close()

	// Ler resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	// Verificar status da resposta
	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		if err := json.Unmarshal(body, &errorResp); err == nil {
			if errorMsg, ok := errorResp["error"].(string); ok {
				return nil, fmt.Errorf("erro do auth-service: %s", errorMsg)
			}
		}
		return nil, fmt.Errorf("auth-service retornou status %d: %s", resp.StatusCode, string(body))
	}

	// Deserializar resposta
	var authUser models.ValidateTokenResponse
	if err := json.Unmarshal(body, &authUser); err != nil {
		return nil, fmt.Errorf("erro ao deserializar resposta: %w", err)
	}

	// Converter para AuthUser
	user := &models.AuthUser{
		ID:        authUser.ID,
		Email:     authUser.Email,
		Name:      authUser.Name,
		Active:    authUser.Active,
		CreatedAt: authUser.CreatedAt,
	}

	return user, nil
}

// HealthCheck verifica se o auth-service está disponível
func (c *AuthClient) HealthCheck() error {
	url := fmt.Sprintf("%s/api/v1/health", c.baseURL)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("erro ao verificar saúde do auth-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth-service não está saudável (status: %d)", resp.StatusCode)
	}

	return nil
}
