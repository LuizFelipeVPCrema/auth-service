package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"auth-service/config"

	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
)

type Database struct {
	DB     *sql.DB
	DBType string
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	var db *sql.DB
	var err error
	var dbType string

	switch strings.ToLower(cfg.Database.Type) {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Name)

		if cfg.Database.SSLMode != "" && cfg.Database.SSLMode != "disable" {
			dsn += "&tls=" + cfg.Database.SSLMode
		}

		db, err = sql.Open("mysql", dsn)
		dbType = "mysql"

	case "sqlite":
		db, err = sql.Open("sqlite", cfg.Database.Path)
		dbType = "sqlite"

	default:
		return nil, fmt.Errorf("tipo de banco de dados não suportado: %s", cfg.Database.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão com banco: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao conectar com banco: %w", err)
	}

	log.Printf("Conexão com banco de dados %s estabelecida com sucesso", strings.ToUpper(dbType))

	return &Database{DB: db, DBType: dbType}, nil
}

func (d *Database) Close() error {
	return d.DB.Close()
}

// Placeholder converte placeholder genérico para o formato específico do banco
func (d *Database) Placeholder(index int) string {
	switch d.DBType {
	case "mysql":
		return "?"
	case "sqlite":
		return "?"
	default:
		return "?"
	}
}

// GetDataType retorna o tipo de dados específico para cada banco
func (d *Database) GetDataType(genericType string) string {
	switch d.DBType {
	case "mysql":
		switch genericType {
		case "TEXT_ID":
			return "VARCHAR(36)"
		case "TEXT":
			return "TEXT"
		case "INTEGER":
			return "INT"
		case "DATETIME":
			return "DATETIME"
		case "BOOLEAN":
			return "BOOLEAN"
		default:
			return genericType
		}
	case "sqlite":
		switch genericType {
		case "TEXT_ID":
			return "TEXT"
		case "TEXT":
			return "TEXT"
		case "INTEGER":
			return "INTEGER"
		case "DATETIME":
			return "DATETIME"
		case "BOOLEAN":
			return "INTEGER"
		default:
			return genericType
		}
	default:
		return genericType
	}
}

func (d *Database) InitTables() error {
	var queries []string

	// Queries específicas para cada tipo de banco
	if d.DBType == "mysql" {
		queries = []string{
			fmt.Sprintf(`CREATE TABLE IF NOT EXISTS users (
				id %s PRIMARY KEY,
				email %s UNIQUE NOT NULL,
				password %s NOT NULL,
				name %s NOT NULL,
				active %s DEFAULT 1,
				created_at %s DEFAULT CURRENT_TIMESTAMP,
				updated_at %s DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
			)`, d.GetDataType("TEXT_ID"), d.GetDataType("TEXT"), d.GetDataType("TEXT"),
				d.GetDataType("TEXT"), d.GetDataType("BOOLEAN"), d.GetDataType("DATETIME"),
				d.GetDataType("DATETIME")),

			fmt.Sprintf(`CREATE TABLE IF NOT EXISTS clients (
				id %s PRIMARY KEY,
				name %s NOT NULL,
				description %s,
				secret %s NOT NULL,
				active %s DEFAULT 1,
				created_at %s DEFAULT CURRENT_TIMESTAMP,
				updated_at %s DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
			)`, d.GetDataType("TEXT_ID"), d.GetDataType("TEXT"), d.GetDataType("TEXT"),
				d.GetDataType("TEXT"), d.GetDataType("BOOLEAN"), d.GetDataType("DATETIME"),
				d.GetDataType("DATETIME")),

			fmt.Sprintf(`CREATE TABLE IF NOT EXISTS refresh_tokens (
				id %s PRIMARY KEY,
				user_id %s NOT NULL,
				client_id %s NOT NULL,
				token %s UNIQUE NOT NULL,
				expires_at %s NOT NULL,
				revoked %s DEFAULT 0,
				created_at %s DEFAULT CURRENT_TIMESTAMP,
				updated_at %s DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
				FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE
			)`, d.GetDataType("TEXT_ID"), d.GetDataType("TEXT_ID"), d.GetDataType("TEXT_ID"),
				d.GetDataType("TEXT"), d.GetDataType("DATETIME"), d.GetDataType("BOOLEAN"),
				d.GetDataType("DATETIME"), d.GetDataType("DATETIME")),
		}
	} else {
		// SQLite
		queries = []string{
			fmt.Sprintf(`CREATE TABLE IF NOT EXISTS users (
				id %s PRIMARY KEY,
				email %s UNIQUE NOT NULL,
				password %s NOT NULL,
				name %s NOT NULL,
				active %s DEFAULT 1,
				created_at %s DEFAULT CURRENT_TIMESTAMP,
				updated_at %s DEFAULT CURRENT_TIMESTAMP
			)`, d.GetDataType("TEXT_ID"), d.GetDataType("TEXT"), d.GetDataType("TEXT"),
				d.GetDataType("TEXT"), d.GetDataType("BOOLEAN"), d.GetDataType("DATETIME"),
				d.GetDataType("DATETIME")),

			fmt.Sprintf(`CREATE TABLE IF NOT EXISTS clients (
				id %s PRIMARY KEY,
				name %s NOT NULL,
				description %s,
				secret %s NOT NULL,
				active %s DEFAULT 1,
				created_at %s DEFAULT CURRENT_TIMESTAMP,
				updated_at %s DEFAULT CURRENT_TIMESTAMP
			)`, d.GetDataType("TEXT_ID"), d.GetDataType("TEXT"), d.GetDataType("TEXT"),
				d.GetDataType("TEXT"), d.GetDataType("BOOLEAN"), d.GetDataType("DATETIME"),
				d.GetDataType("DATETIME")),

			fmt.Sprintf(`CREATE TABLE IF NOT EXISTS refresh_tokens (
				id %s PRIMARY KEY,
				user_id %s NOT NULL,
				client_id %s NOT NULL,
				token %s UNIQUE NOT NULL,
				expires_at %s NOT NULL,
				revoked %s DEFAULT 0,
				created_at %s DEFAULT CURRENT_TIMESTAMP,
				updated_at %s DEFAULT CURRENT_TIMESTAMP,
				FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
				FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE
			)`, d.GetDataType("TEXT_ID"), d.GetDataType("TEXT_ID"), d.GetDataType("TEXT_ID"),
				d.GetDataType("TEXT"), d.GetDataType("DATETIME"), d.GetDataType("BOOLEAN"),
				d.GetDataType("DATETIME"), d.GetDataType("DATETIME")),
		}
	}

	// Índices comuns
	indexQueries := []string{
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`,
		`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token)`,
	}

	queries = append(queries, indexQueries...)

	for _, query := range queries {
		if _, err := d.DB.Exec(query); err != nil {
			return fmt.Errorf("erro ao executar query: %w", err)
		}
	}

	log.Printf("Tabelas %s inicializadas com sucesso", strings.ToUpper(d.DBType))
	return nil
}
