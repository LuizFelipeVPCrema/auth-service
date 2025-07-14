package services

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"auth-service/config"
	"auth-service/database"
	"auth-service/models"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService struct {
	db   *database.Database
	cfg  *config.Config
	hash *argon2id.Hash
}

func NewAuthService(db *database.Database, cfg *config.Config) *AuthService {
	hash := &argon2id.Hash{
		Memory:      cfg.Hash.Memory,
		Iterations:  cfg.Hash.Iterations,
		Parallelism: cfg.Hash.Parallelism,
		SaltLength:  cfg.Hash.SaltLength,
		KeyLength:   cfg.Hash.KeyLength,
	}

	return &AuthService{
		db:   db,
		cfg:  cfg,
		hash: hash,
	}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.UserResponse, error) {
	// Verificar se o email já existe
	var existingUser models.User
	err := s.db.DB.QueryRow("SELECT id FROM users WHERE email = $1", req.Email).Scan(&existingUser.ID)
	if err != sql.ErrNoRows {
		if err == nil {
			return nil, fmt.Errorf("email já está em uso")
		}
		return nil, fmt.Errorf("erro ao verificar email: %w", err)
	}

	// Hash da senha
	hashedPassword, err := argon2id.CreateHash(req.Password, s.hash)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar hash da senha: %w", err)
	}

	// Criar usuário
	userID := uuid.New()
	now := time.Now()

	_, err = s.db.DB.Exec(`
		INSERT INTO users (id, email, password, name, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, userID, req.Email, hashedPassword, req.Name, true, now, now)

	if err != nil {
		return nil, fmt.Errorf("erro ao criar usuário: %w", err)
	}

	return &models.UserResponse{
		ID:        userID,
		Email:     req.Email,
		Name:      req.Name,
		Active:    true,
		CreatedAt: now,
	}, nil
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.TokenResponse, error) {
	// Verificar se o cliente existe
	var client models.Client
	err := s.db.DB.QueryRow("SELECT id, secret, active FROM clients WHERE id = $1", req.ClientID).Scan(&client.ID, &client.Secret, &client.Active)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cliente não encontrado")
		}
		return nil, fmt.Errorf("erro ao verificar cliente: %w", err)
	}

	if !client.Active {
		return nil, fmt.Errorf("cliente inativo")
	}

	// Buscar usuário
	var user models.User
	err = s.db.DB.QueryRow("SELECT id, email, password, name, active FROM users WHERE email = $1", req.Email).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name, &user.Active)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("credenciais inválidas")
		}
		return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
	}

	if !user.Active {
		return nil, fmt.Errorf("usuário inativo")
	}

	// Verificar senha
	match, err := argon2id.ComparePasswordAndHash(req.Password, user.Password)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar senha: %w", err)
	}

	if !match {
		return nil, fmt.Errorf("credenciais inválidas")
	}

	// Gerar tokens
	accessToken, err := s.generateAccessToken(user, client.ID.String())
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(user.ID, client.ID)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar refresh token: %w", err)
	}

	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.cfg.JWT.ExpirationHours * 3600), // segundos
	}, nil
}

func (s *AuthService) RefreshToken(req *models.RefreshTokenRequest) (*models.TokenResponse, error) {
	// Verificar se o cliente existe
	var client models.Client
	err := s.db.DB.QueryRow("SELECT id, active FROM clients WHERE id = $1", req.ClientID).Scan(&client.ID, &client.Active)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cliente não encontrado")
		}
		return nil, fmt.Errorf("erro ao verificar cliente: %w", err)
	}

	if !client.Active {
		return nil, fmt.Errorf("cliente inativo")
	}

	// Verificar refresh token
	var refreshToken models.RefreshToken
	err = s.db.DB.QueryRow(`
		SELECT id, user_id, client_id, token, expires_at, revoked 
		FROM refresh_tokens 
		WHERE token = $1 AND client_id = $2
	`, req.RefreshToken, client.ID).Scan(
		&refreshToken.ID, &refreshToken.UserID, &refreshToken.ClientID,
		&refreshToken.Token, &refreshToken.ExpiresAt, &refreshToken.Revoked)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("refresh token inválido")
		}
		return nil, fmt.Errorf("erro ao verificar refresh token: %w", err)
	}

	if refreshToken.Revoked {
		return nil, fmt.Errorf("refresh token revogado")
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		return nil, fmt.Errorf("refresh token expirado")
	}

	// Buscar usuário
	var user models.User
	err = s.db.DB.QueryRow("SELECT id, email, name, active FROM users WHERE id = $1", refreshToken.UserID).Scan(
		&user.ID, &user.Email, &user.Name, &user.Active)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
	}

	if !user.Active {
		return nil, fmt.Errorf("usuário inativo")
	}

	// Revogar refresh token atual
	_, err = s.db.DB.Exec("UPDATE refresh_tokens SET revoked = true WHERE id = $1", refreshToken.ID)
	if err != nil {
		return nil, fmt.Errorf("erro ao revogar refresh token: %w", err)
	}

	// Gerar novos tokens
	accessToken, err := s.generateAccessToken(user, client.ID.String())
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar access token: %w", err)
	}

	newRefreshToken, err := s.generateRefreshToken(user.ID, client.ID)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar refresh token: %w", err)
	}

	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.cfg.JWT.ExpirationHours * 3600),
	}, nil
}

func (s *AuthService) ValidateToken(tokenString, clientID string) (*models.UserResponse, error) {
	// Verificar se o cliente existe
	var client models.Client
	err := s.db.DB.QueryRow("SELECT id, active FROM clients WHERE id = $1", clientID).Scan(&client.ID, &client.Active)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cliente não encontrado")
		}
		return nil, fmt.Errorf("erro ao verificar cliente: %w", err)
	}

	if !client.Active {
		return nil, fmt.Errorf("cliente inativo")
	}

	// Validar token JWT
	token, err := jwt.ParseWithClaims(tokenString, &models.JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return []byte(s.cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token inválido: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token inválido")
	}

	claims, ok := token.Claims.(*models.JWTCustomClaims)
	if !ok {
		return nil, fmt.Errorf("claims inválidas")
	}

	if claims.Type != "access" {
		return nil, fmt.Errorf("tipo de token inválido")
	}

	if claims.ClientID != clientID {
		return nil, fmt.Errorf("cliente não autorizado")
	}

	// Buscar usuário
	var user models.User
	err = s.db.DB.QueryRow("SELECT id, email, name, active, created_at FROM users WHERE id = $1", claims.UserID).Scan(
		&user.ID, &user.Email, &user.Name, &user.Active, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
	}

	if !user.Active {
		return nil, fmt.Errorf("usuário inativo")
	}

	return &models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Active:    user.Active,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *AuthService) generateAccessToken(user models.User, clientID string) (string, error) {
	claims := models.JWTCustomClaims{
		UserID:   user.ID.String(),
		Email:    user.Email,
		ClientID: clientID,
		Type:     "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   claims.UserID,
		"email":     claims.Email,
		"client_id": claims.ClientID,
		"type":      claims.Type,
		"exp":       time.Now().Add(time.Duration(s.cfg.JWT.ExpirationHours) * time.Hour).Unix(),
		"iat":       time.Now().Unix(),
	})

	return token.SignedString([]byte(s.cfg.JWT.Secret))
}

func (s *AuthService) generateRefreshToken(userID, clientID uuid.UUID) (string, error) {
	// Gerar token aleatório
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("erro ao gerar token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// Salvar no banco
	refreshTokenID := uuid.New()
	expiresAt := time.Now().Add(time.Duration(s.cfg.JWT.RefreshExpirationHours) * time.Hour)
	now := time.Now()

	_, err := s.db.DB.Exec(`
		INSERT INTO refresh_tokens (id, user_id, client_id, token, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, refreshTokenID, userID, clientID, token, expiresAt, now, now)

	if err != nil {
		return "", fmt.Errorf("erro ao salvar refresh token: %w", err)
	}

	return token, nil
}

func (s *AuthService) CreateClient(name, description string) (*models.Client, error) {
	// Gerar secret aleatório
	secretBytes := make([]byte, 32)
	if _, err := rand.Read(secretBytes); err != nil {
		return nil, fmt.Errorf("erro ao gerar secret: %w", err)
	}
	secret := hex.EncodeToString(secretBytes)

	clientID := uuid.New()
	now := time.Now()

	_, err := s.db.DB.Exec(`
		INSERT INTO clients (id, name, description, secret, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, clientID, name, description, secret, true, now, now)

	if err != nil {
		return nil, fmt.Errorf("erro ao criar cliente: %w", err)
	}

	return &models.Client{
		ID:          clientID,
		Name:        name,
		Description: description,
		Secret:      secret,
		Active:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}
