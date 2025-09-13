package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config representa todas as configurações da aplicação
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	Security SecurityConfig
	Upload   UploadConfig
}

// ServerConfig configurações do servidor
type ServerConfig struct {
	Port string
	Mode string
}

// DatabaseConfig configurações do banco de dados
type DatabaseConfig struct {
	Path string
}

// AuthConfig configurações do microserviço de autenticação
type AuthConfig struct {
	ServiceURL string
	ClientID   string
	Timeout    int // em segundos
}

// SecurityConfig configurações de segurança
type SecurityConfig struct {
	RateLimitRequests      int
	RateLimitWindowMinutes int
	AuditEnabled           bool
	AllowedOrigins         []string
	MaxRequestSize         int64 // em bytes
}

// UploadConfig configurações de upload de arquivos
type UploadConfig struct {
	Path            string
	MaxSizeMB       int
	AllowedTypes    []string
	AllowedMimeTypes []string
}

// Load carrega as configurações do arquivo .env e variáveis de ambiente
func Load() *Config {
	// Carregar arquivo .env se existir
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Path: getEnv("DB_PATH", "study_manager.db"),
		},
		Auth: AuthConfig{
			ServiceURL: getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
			ClientID:   getEnv("CLIENT_ID", ""),
			Timeout:    getEnvAsInt("AUTH_TIMEOUT_SECONDS", 30),
		},
		Security: SecurityConfig{
			RateLimitRequests:      getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
			RateLimitWindowMinutes: getEnvAsInt("RATE_LIMIT_WINDOW_MINUTES", 1),
			AuditEnabled:           getEnvAsBool("AUDIT_ENABLED", true),
			AllowedOrigins:         getEnvAsSlice("ALLOWED_ORIGINS", []string{"http://localhost:3000", "http://localhost:4000", "http://localhost:4200"}),
			MaxRequestSize:         int64(getEnvAsInt("MAX_REQUEST_SIZE_MB", 10)) * 1024 * 1024, // Converter MB para bytes
		},
		Upload: UploadConfig{
			Path:             getEnv("UPLOAD_PATH", "./uploads"),
			MaxSizeMB:        getEnvAsInt("MAX_UPLOAD_SIZE_MB", 10),
			AllowedTypes:     getEnvAsSlice("ALLOWED_FILE_TYPES", []string{".pdf", ".doc", ".docx", ".txt", ".jpg", ".jpeg", ".png", ".gif"}),
			AllowedMimeTypes: getEnvAsSlice("ALLOWED_MIME_TYPES", []string{"application/pdf", "application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", "text/plain", "image/jpeg", "image/png", "image/gif"}),
		},
	}
}

// getEnv obtém uma variável de ambiente ou retorna um valor padrão
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt obtém uma variável de ambiente como inteiro ou retorna um valor padrão
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool obtém uma variável de ambiente como booleano ou retorna um valor padrão
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getEnvAsSlice obtém uma variável de ambiente como slice de strings ou retorna um valor padrão
func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}
