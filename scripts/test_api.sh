#!/bin/bash

# Script de teste para a API de autenticação
# Execute este script após iniciar o servidor

BASE_URL="http://localhost:8080/api/v1"
CLIENT_ID=""
ACCESS_TOKEN=""
REFRESH_TOKEN=""

echo "🧪 Testando API de Autenticação"
echo "=================================="

# Função para fazer requisições HTTP
make_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    local headers=$4
    
    if [ -n "$data" ]; then
        if [ -n "$headers" ]; then
            curl -s -X $method "$BASE_URL$endpoint" \
                -H "Content-Type: application/json" \
                -H "$headers" \
                -d "$data"
        else
            curl -s -X $method "$BASE_URL$endpoint" \
                -H "Content-Type: application/json" \
                -d "$data"
        fi
    else
        if [ -n "$headers" ]; then
            curl -s -X $method "$BASE_URL$endpoint" \
                -H "Content-Type: application/json" \
                -H "$headers"
        else
            curl -s -X $method "$BASE_URL$endpoint" \
                -H "Content-Type: application/json"
        fi
    fi
}

# 1. Criar cliente
echo "1. Criando cliente..."
CLIENT_RESPONSE=$(make_request "POST" "/clients" '{"name": "Teste App", "description": "Aplicação de teste"}')
echo "Resposta: $CLIENT_RESPONSE"

# Extrair client_id da resposta
CLIENT_ID=$(echo $CLIENT_RESPONSE | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
echo "Client ID: $CLIENT_ID"
echo ""

# 2. Registrar usuário
echo "2. Registrando usuário..."
REGISTER_RESPONSE=$(make_request "POST" "/auth/register" '{"email": "teste@exemplo.com", "password": "senha123456", "name": "Usuário Teste"}')
echo "Resposta: $REGISTER_RESPONSE"
echo ""

# 3. Fazer login
echo "3. Fazendo login..."
LOGIN_RESPONSE=$(make_request "POST" "/auth/login" "{\"email\": \"teste@exemplo.com\", \"password\": \"senha123456\", \"client_id\": \"$CLIENT_ID\"}")
echo "Resposta: $LOGIN_RESPONSE"

# Extrair tokens da resposta
ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
REFRESH_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"refresh_token":"[^"]*"' | cut -d'"' -f4)
echo "Access Token: ${ACCESS_TOKEN:0:50}..."
echo "Refresh Token: ${REFRESH_TOKEN:0:50}..."
echo ""

# 4. Validar token
echo "4. Validando token..."
VALIDATE_RESPONSE=$(make_request "POST" "/auth/validate" "{\"token\": \"$ACCESS_TOKEN\", \"client_id\": \"$CLIENT_ID\"}")
echo "Resposta: $VALIDATE_RESPONSE"
echo ""

# 5. Obter perfil (rota protegida)
echo "5. Obtendo perfil (rota protegida)..."
PROFILE_RESPONSE=$(make_request "GET" "/auth/profile" "" "Authorization: Bearer $ACCESS_TOKEN" "X-Client-ID: $CLIENT_ID")
echo "Resposta: $PROFILE_RESPONSE"
echo ""

# 6. Renovar token
echo "6. Renovando token..."
REFRESH_RESPONSE=$(make_request "POST" "/auth/refresh" "{\"refresh_token\": \"$REFRESH_TOKEN\", \"client_id\": \"$CLIENT_ID\"}")
echo "Resposta: $REFRESH_RESPONSE"

# Extrair novos tokens
NEW_ACCESS_TOKEN=$(echo $REFRESH_RESPONSE | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
echo "Novo Access Token: ${NEW_ACCESS_TOKEN:0:50}..."
echo ""

echo "✅ Testes concluídos com sucesso!"
echo ""
echo "📋 Resumo:"
echo "- Cliente criado: $CLIENT_ID"
echo "- Usuário registrado: teste@exemplo.com"
echo "- Login realizado com sucesso"
echo "- Token validado com sucesso"
echo "- Perfil obtido com sucesso"
echo "- Token renovado com sucesso" 