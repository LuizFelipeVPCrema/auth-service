# Build stage
FROM golang:1.25.0-alpine AS builder

# Instalar dependências necessárias
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Definir diretório de trabalho
WORKDIR /app

# Copiar arquivos de dependências
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Runtime stage
FROM alpine:latest

# Instalar certificados CA, SQLite e curl
RUN apk --no-cache add ca-certificates sqlite curl

# Criar usuário não-root
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Definir diretório de trabalho
WORKDIR /app

# Criar diretório data e ajustar permissões
RUN mkdir -p /app/data && \
    chown -R appuser:appgroup /app/data

# Copiar binário do builder
COPY --from=builder /app/main .

# Mudar propriedade do arquivo
RUN chown appuser:appgroup main

# Mudar para usuário não-root
USER appuser

# Expor porta
EXPOSE 8080

# Comando para executar a aplicação
CMD ["./main"] 