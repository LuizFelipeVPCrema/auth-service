package models

// AuthUser representa os dados do usuário autenticado vindos do auth-service
type AuthUser struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"created_at"`
}

// ValidateTokenRequest representa a requisição para validar token
type ValidateTokenRequest struct {
	Token    string `json:"token" binding:"required"`
	ClientID string `json:"client_id" binding:"required"`
}

// ValidateTokenResponse representa a resposta da validação de token
type ValidateTokenResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"created_at"`
}
