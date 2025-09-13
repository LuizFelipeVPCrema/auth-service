package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"study-manager-service/internal/clients"
	"study-manager-service/internal/config"
	"study-manager-service/internal/database"
	"study-manager-service/internal/middleware"
	"study-manager-service/internal/repositories"
	"study-manager-service/internal/routes"
	"study-manager-service/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Carregar configurações
	cfg := config.Load()

	// Configurar modo do Gin
	gin.SetMode(cfg.Server.Mode)

	// Conectar ao banco de dados
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Executar migrações
	if err := db.Migrate(); err != nil {
		log.Fatalf("Erro ao executar migrações: %v", err)
	}

	// Criar índices
	if err := db.CreateIndexes(); err != nil {
		log.Fatalf("Erro ao criar índices: %v", err)
	}

	// Inicializar cliente de autenticação
	authClient := clients.NewAuthClient(cfg)

	// Verificar conectividade com auth-service
	if err := authClient.HealthCheck(); err != nil {
		log.Printf("Aviso: Auth-service não está disponível: %v", err)
	} else {
		log.Println("Conexão com auth-service estabelecida com sucesso")
	}

	// Inicializar repositórios
	studentRepo := repositories.NewStudentRepository(db.DB)
	subjectRepo := repositories.NewSubjectRepository(db.DB)
	examRepo := repositories.NewExamRepository(db.DB)

	// Inicializar serviços
	studentService := services.NewStudentService(studentRepo, authClient)
	subjectService := services.NewSubjectService(subjectRepo, authClient)
	examService := services.NewExamService(examRepo, authClient)

	// Inicializar middlewares
	authMiddleware := middleware.NewAuthMiddleware(authClient)
	rateLimiter := middleware.NewRateLimiter(cfg)
	validationMiddleware := middleware.NewValidationMiddleware(cfg)
	auditLogger := middleware.NewAuditLogger(cfg)

	// Configurar rotas
	router := routes.SetupRoutes(
		cfg,
		authMiddleware,
		rateLimiter,
		validationMiddleware,
		auditLogger,
		studentService,
		subjectService,
		examService,
	)

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
		log.Printf("Study Manager Service iniciado na porta %s", cfg.Server.Port)
		log.Printf("Ambiente: %s", cfg.Server.Mode)
		log.Printf("Documentação da API disponível em: http://localhost:%s/api/v1", cfg.Server.Port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	// Aguardar sinal de interrupção
	<-quit
	log.Println("Desligando Study Manager Service...")

	// Contexto com timeout para shutdown graceful
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Tentar desligar o servidor gracefulmente
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Erro ao desligar servidor: %v", err)
	}

	log.Println("Study Manager Service desligado com sucesso")
}