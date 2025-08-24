package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	Type     string // "sqlite" ou "mysql"
	Path     string // Para SQLite
	Host     string // Para MySQL
	Port     string // Para MySQL
	User     string // Para MySQL
	Password string // Para MySQL
	Name     string // Para MySQL
	SSLMode  string // Para MySQL
}

type JWTConfig struct {
	Secret                 string
	ExpirationHours        int
	RefreshExpirationHours int
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Type:     getEnv("DB_TYPE", "sqlite"),
			Path:     getEnv("DB_PATH", "auth_service.db"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "auth_service"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:                 getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
			ExpirationHours:        getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
			RefreshExpirationHours: getEnvAsInt("JWT_REFRESH_EXPIRATION_HOURS", 168),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
