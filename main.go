package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"auth-service/config"
	"auth-service/database"
	"auth-service/handlers"
	"auth-service/middleware"
	"auth-service/routes"
	"auth-service/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Carregar configurações
	cfg := config.Load()

	// Configurar modo do Gin
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Conectar ao banco de dados
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Inicializar tabelas
	if err := db.InitTables(); err != nil {
		log.Fatalf("Erro ao inicializar tabelas: %v", err)
	}

	// Inicializar serviços
	authService := services.NewAuthService(db, cfg)

	// Inicializar handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Inicializar middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Configurar rotas
	router := routes.SetupRoutes(authHandler, authMiddleware)

	// Configurar servidor HTTP
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Canal para receber sinais de interrupção
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor em uma goroutine
	go func() {
		log.Printf("Servidor iniciado na porta %s", cfg.Server.Port)
		log.Printf("Ambiente: %s", cfg.Server.Env)
		log.Printf("Documentação da API disponível em: http://localhost:%s/api/v1", cfg.Server.Port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	// Aguardar sinal de interrupção
	<-quit
	log.Println("Desligando servidor...")

	// Contexto com timeout para shutdown graceful
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Tentar desligar o servidor gracefulmente
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Erro ao desligar servidor: %v", err)
	}

	log.Println("Servidor desligado com sucesso")
}
