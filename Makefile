# Makefile para Study Manager Service

.PHONY: help build run test clean docker-build docker-run docker-compose-up docker-compose-down

# Variáveis
BINARY_NAME=study-manager-service
DOCKER_IMAGE=study-manager-service
DOCKER_TAG=latest

# Cores para output
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

help: ## Mostra esta ajuda
	@echo "$(GREEN)Study Manager Service - Comandos disponíveis:$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

build: ## Compila a aplicação
	@echo "$(GREEN)Compilando $(BINARY_NAME)...$(NC)"
	go build -o $(BINARY_NAME) main.go
	@echo "$(GREEN)Compilação concluída!$(NC)"

run: ## Executa a aplicação
	@echo "$(GREEN)Executando $(BINARY_NAME)...$(NC)"
	go run main.go

test: ## Executa os testes
	@echo "$(GREEN)Executando testes...$(NC)"
	go test -v ./...

test-coverage: ## Executa testes com cobertura
	@echo "$(GREEN)Executando testes com cobertura...$(NC)"
	go test -v -cover ./...

clean: ## Remove arquivos compilados
	@echo "$(GREEN)Limpando arquivos...$(NC)"
	rm -f $(BINARY_NAME)
	rm -f study_manager.db
	rm -rf uploads/
	@echo "$(GREEN)Limpeza concluída!$(NC)"

deps: ## Instala dependências
	@echo "$(GREEN)Instalando dependências...$(NC)"
	go mod tidy
	go mod download

dev: ## Executa em modo desenvolvimento com hot reload
	@echo "$(GREEN)Executando em modo desenvolvimento...$(NC)"
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "$(YELLOW)Air não encontrado. Instalando...$(NC)"; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

docker-build: ## Constrói a imagem Docker
	@echo "$(GREEN)Construindo imagem Docker...$(NC)"
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "$(GREEN)Imagem construída: $(DOCKER_IMAGE):$(DOCKER_TAG)$(NC)"

docker-run: ## Executa o container Docker
	@echo "$(GREEN)Executando container Docker...$(NC)"
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-compose-up: ## Inicia os serviços com Docker Compose
	@echo "$(GREEN)Iniciando serviços com Docker Compose...$(NC)"
	docker-compose up -d
	@echo "$(GREEN)Serviços iniciados!$(NC)"
	@echo "$(YELLOW)Auth Service: http://localhost:8081$(NC)"
	@echo "$(YELLOW)Study Manager Service: http://localhost:8080$(NC)"

docker-compose-down: ## Para os serviços do Docker Compose
	@echo "$(GREEN)Parando serviços...$(NC)"
	docker-compose down
	@echo "$(GREEN)Serviços parados!$(NC)"

docker-compose-logs: ## Mostra os logs dos serviços
	@echo "$(GREEN)Mostrando logs...$(NC)"
	docker-compose logs -f

test-integration: ## Executa testes de integração
	@echo "$(GREEN)Executando testes de integração...$(NC)"
	chmod +x scripts/test_integration.sh
	./scripts/test_integration.sh

setup: deps ## Configura o ambiente de desenvolvimento
	@echo "$(GREEN)Configurando ambiente...$(NC)"
	@if [ ! -f .env ]; then \
		cp env.example .env; \
		echo "$(YELLOW)Arquivo .env criado. Configure as variáveis necessárias.$(NC)"; \
	fi
	@echo "$(GREEN)Ambiente configurado!$(NC)"

check: ## Verifica o código
	@echo "$(GREEN)Verificando código...$(NC)"
	go vet ./...
	go fmt ./...
	@echo "$(GREEN)Verificação concluída!$(NC)"

install-tools: ## Instala ferramentas de desenvolvimento
	@echo "$(GREEN)Instalando ferramentas...$(NC)"
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(GREEN)Ferramentas instaladas!$(NC)"

lint: ## Executa linter
	@echo "$(GREEN)Executando linter...$(NC)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)golangci-lint não encontrado. Instalando...$(NC)"; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi

# Comando padrão
.DEFAULT_GOAL := help