# Microserviço de Autenticação

Um microserviço de autenticação universal desenvolvido em Go, utilizando JWT para autenticação e PostgreSQL como banco de dados.

## 🚀 Funcionalidades

- **Registro de usuários** - Criação de contas com validação de email
- **Login** - Autenticação com geração de tokens JWT
- **Validação de tokens** - Verificação de tokens de acesso
- **Renovação de tokens** - Sistema de refresh tokens
- **Middleware de autenticação** - Proteção de rotas
- **Suporte a múltiplos clientes** - Autenticação para diferentes aplicações
- **Hash seguro de senhas** - Utilizando Argon2id

## 🛠️ Stack Tecnológica

- **Linguagem**: Go (Golang) 1.21+
- **Framework Web**: Gin
- **Autenticação**: JWT (JSON Web Tokens)
- **Banco de Dados**: PostgreSQL
- **Hash de Senhas**: Argon2id
- **Configuração**: Variáveis de ambiente (.env)

## 📋 Pré-requisitos

- Go 1.21 ou superior
- PostgreSQL 12 ou superior
- Git

## 🔧 Instalação

1. **Clone o repositório**
   ```bash
   git clone <url-do-repositorio>
   cd auth-service
   ```

2. **Instale as dependências**
   ```bash
   go mod tidy
   ```

3. **Configure o banco de dados**
   - Crie um banco PostgreSQL
   - Configure as variáveis de ambiente (veja `.env.example`)

4. **Configure as variáveis de ambiente**
   ```bash
   cp env.example .env
   # Edite o arquivo .env com suas configurações
   ```

5. **Execute o projeto**
   ```bash
   go run main.go
   ```

## ⚙️ Configuração

### Variáveis de Ambiente

Copie o arquivo `env.example` para `.env` e configure:

```env
# Configurações do Servidor
PORT=8080
ENV=development

# Configurações do Banco de Dados
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=auth_service
DB_SSL_MODE=disable

# Configurações JWT
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRATION_HOURS=24
JWT_REFRESH_EXPIRATION_HOURS=168

# Configurações de Hash
HASH_MEMORY=64
HASH_ITERATIONS=3
HASH_PARALLELISM=2
HASH_SALT_LENGTH=16
HASH_KEY_LENGTH=32
```

## 📚 API Endpoints

### Autenticação

#### POST `/api/v1/auth/register`
Registra um novo usuário.

**Request Body:**
```json
{
  "email": "usuario@exemplo.com",
  "password": "senha123456",
  "name": "Nome do Usuário"
}
```

**Response:**
```json
{
  "id": "uuid-do-usuario",
  "email": "usuario@exemplo.com",
  "name": "Nome do Usuário",
  "active": true,
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### POST `/api/v1/auth/login`
Faz login do usuário.

**Request Body:**
```json
{
  "email": "usuario@exemplo.com",
  "password": "senha123456",
  "client_id": "uuid-do-cliente"
}
```

**Response:**
```json
{
  "access_token": "jwt-token",
  "refresh_token": "refresh-token",
  "token_type": "Bearer",
  "expires_in": 86400
}
```

#### POST `/api/v1/auth/refresh`
Renova o access token usando refresh token.

**Request Body:**
```json
{
  "refresh_token": "refresh-token",
  "client_id": "uuid-do-cliente"
}
```

#### POST `/api/v1/auth/validate`
Valida um access token.

**Request Body:**
```json
{
  "token": "jwt-token",
  "client_id": "uuid-do-cliente"
}
```

#### GET `/api/v1/auth/profile` (Protegido)
Retorna informações do usuário autenticado.

**Headers:**
```
Authorization: Bearer <jwt-token>
X-Client-ID: <client-id>
```

### Clientes

#### POST `/api/v1/clients`
Cria um novo cliente para autenticação.

**Request Body:**
```json
{
  "name": "Nome do Cliente",
  "description": "Descrição do cliente"
}
```

**Response:**
```json
{
  "id": "uuid-do-cliente",
  "name": "Nome do Cliente",
  "description": "Descrição do cliente",
  "secret": "secret-gerado",
  "active": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

## 🔐 Segurança

### Tokens JWT

- **Access Token**: Expira em 24 horas (configurável)
- **Refresh Token**: Expira em 7 dias (configurável)
- **Claims customizáveis**: Inclui user_id, email, client_id e tipo
- **Assinatura HMAC**: Utiliza chave secreta configurável

### Hash de Senhas

- **Algoritmo**: Argon2id
- **Configurações otimizadas** para segurança vs performance
- **Salt único** para cada senha

### Validações

- **Email único** por usuário
- **Senha mínima** de 8 caracteres
- **Cliente ativo** obrigatório para autenticação
- **Usuário ativo** obrigatório para operações

## 🏗️ Arquitetura

```
auth-service/
├── config/          # Configurações e variáveis de ambiente
├── database/        # Conexão e inicialização do banco
├── handlers/        # Handlers HTTP (controllers)
├── middleware/      # Middleware de autenticação
├── models/          # Modelos de dados
├── routes/          # Definição de rotas
├── services/        # Lógica de negócio
├── main.go          # Ponto de entrada da aplicação
├── go.mod           # Dependências Go
└── README.md        # Documentação
```

## 🧪 Testando a API

### 1. Criar um cliente
```bash
curl -X POST http://localhost:8080/api/v1/clients \
  -H "Content-Type: application/json" \
  -d '{"name": "Teste App", "description": "Aplicação de teste"}'
```

### 2. Registrar um usuário
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "teste@exemplo.com", "password": "senha123456", "name": "Usuário Teste"}'
```

### 3. Fazer login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "teste@exemplo.com", "password": "senha123456", "client_id": "UUID_DO_CLIENTE"}'
```

### 4. Validar token
```bash
curl -X POST http://localhost:8080/api/v1/auth/validate \
  -H "Content-Type: application/json" \
  -d '{"token": "JWT_TOKEN", "client_id": "UUID_DO_CLIENTE"}'
```

### 5. Obter perfil (protegido)
```bash
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer JWT_TOKEN" \
  -H "X-Client-ID: UUID_DO_CLIENTE"
```

## 🚀 Deploy

### Docker

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

### Build
```bash
go build -o auth-service main.go
```

## 📝 Licença

Este projeto está sob a licença MIT. Veja o arquivo LICENSE para mais detalhes.

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📞 Suporte

Para suporte, envie um email para [seu-email@exemplo.com] ou abra uma issue no repositório. 