# Makefile para o microserviço de autenticação

.PHONY: help build run test clean docker-build docker-run docker-stop deps lint

# Variáveis
BINARY_NAME=auth-service
MAIN_FILE=main.go
DOCKER_IMAGE=auth-service:latest

# Comando padrão
help: ## Mostra esta ajuda
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Desenvolvimento
deps: ## Instala as dependências
	go mod tidy
	go mod download

build: ## Compila a aplicação
	go build -o $(BINARY_NAME) $(MAIN_FILE)

run: ## Executa a aplicação
	go run $(MAIN_FILE)

dev: ## Executa em modo desenvolvimento com hot reload (requer air)
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air não encontrado. Instalando..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

test: ## Executa os testes
	go test -v ./...

test-coverage: ## Executa os testes com cobertura
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Relatório de cobertura gerado em coverage.html"

lint: ## Executa o linter
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint não encontrado. Instalando..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi

# Docker
docker-build: ## Constrói a imagem Docker
	docker build -t $(DOCKER_IMAGE) .

docker-run: ## Executa com Docker Compose
	docker-compose up -d

docker-stop: ## Para os containers Docker
	docker-compose down

docker-logs: ## Mostra os logs dos containers
	docker-compose logs -f

# Banco de dados
db-setup: ## Configura o banco de dados (requer PostgreSQL local)
	@echo "Criando banco de dados..."
	@psql -U postgres -c "CREATE DATABASE auth_service;" 2>/dev/null || echo "Banco já existe ou erro na conexão"

db-migrate: ## Executa as migrações (inicializa tabelas)
	go run $(MAIN_FILE)

# Limpeza
clean: ## Remove arquivos gerados
	rm -f $(BINARY_NAME)
	rm -f coverage.out
	rm -f coverage.html
	go clean

# Scripts de teste
test-api: ## Executa o script de teste da API
	@if [ -f scripts/test_api.sh ]; then \
		chmod +x scripts/test_api.sh; \
		./scripts/test_api.sh; \
	else \
		echo "Script de teste não encontrado"; \
	fi

# Instalação de ferramentas
install-tools: ## Instala ferramentas de desenvolvimento
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Verificação de saúde
health-check: ## Verifica se o serviço está funcionando
	@echo "Verificando saúde do serviço..."
	@curl -f http://localhost:8080/api/v1/auth/register > /dev/null 2>&1 && echo "✅ Serviço está funcionando" || echo "❌ Serviço não está respondendo"

# Desenvolvimento completo
setup: deps install-tools ## Configura o ambiente de desenvolvimento
	@echo "✅ Ambiente configurado com sucesso!"

dev-full: setup run ## Configura e executa o ambiente completo 