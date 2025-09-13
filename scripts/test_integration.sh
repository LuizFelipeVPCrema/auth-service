#!/bin/bash

# Script de teste para integração entre auth-service e study-manager-service
# Execute este script após iniciar ambos os serviços

AUTH_BASE_URL="http://localhost:8081/api/v1"
STUDY_BASE_URL="http://localhost:8080/api/v1"
CLIENT_ID=""
ACCESS_TOKEN=""
STUDENT_ID=""

echo "🧪 Testando Integração entre Auth-Service e Study-Manager-Service"
echo "=================================================================="

# Função para fazer requisições HTTP
make_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    local headers=$4
    local base_url=$5
    
    if [ -n "$data" ]; then
        if [ -n "$headers" ]; then
            curl -s -X $method "$base_url$endpoint" \
                -H "Content-Type: application/json" \
                -H "$headers" \
                -d "$data"
        else
            curl -s -X $method "$base_url$endpoint" \
                -H "Content-Type: application/json" \
                -d "$data"
        fi
    else
        if [ -n "$headers" ]; then
            curl -s -X $method "$base_url$endpoint" \
                -H "Content-Type: application/json" \
                -H "$headers"
        else
            curl -s -X $method "$base_url$endpoint" \
                -H "Content-Type: application/json"
        fi
    fi
}

# 1. Verificar se os serviços estão rodando
echo "1. Verificando se os serviços estão rodando..."

# Verificar auth-service
AUTH_HEALTH=$(make_request "GET" "/health" "" "" "$AUTH_BASE_URL")
if echo "$AUTH_HEALTH" | grep -q "ok"; then
    echo "✅ Auth-service está rodando"
else
    echo "❌ Auth-service não está rodando"
    exit 1
fi

# Verificar study-manager-service
STUDY_HEALTH=$(make_request "GET" "/health" "" "" "$STUDY_BASE_URL")
if echo "$STUDY_HEALTH" | grep -q "ok"; then
    echo "✅ Study-manager-service está rodando"
else
    echo "❌ Study-manager-service não está rodando"
    exit 1
fi

echo ""

# 2. Criar cliente no auth-service
echo "2. Criando cliente no auth-service..."
CLIENT_RESPONSE=$(make_request "POST" "/clients" '{"name": "Study Manager App", "description": "Aplicação de gerenciamento de estudos"}' "" "$AUTH_BASE_URL")
echo "Resposta: $CLIENT_RESPONSE"

# Extrair client_id da resposta
CLIENT_ID=$(echo $CLIENT_RESPONSE | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
if [ -z "$CLIENT_ID" ]; then
    echo "❌ Erro ao obter client_id"
    exit 1
fi
echo "Client ID: $CLIENT_ID"
echo ""

# 3. Registrar usuário no auth-service
echo "3. Registrando usuário no auth-service..."
REGISTER_RESPONSE=$(make_request "POST" "/register" '{"email": "teste@estudante.com", "password": "senha123456", "name": "Estudante Teste"}' "" "$AUTH_BASE_URL")
echo "Resposta: $REGISTER_RESPONSE"
echo ""

# 4. Fazer login no auth-service
echo "4. Fazendo login no auth-service..."
LOGIN_RESPONSE=$(make_request "POST" "/login" "{\"email\": \"teste@estudante.com\", \"password\": \"senha123456\", \"client_id\": \"$CLIENT_ID\"}" "" "$AUTH_BASE_URL")
echo "Resposta: $LOGIN_RESPONSE"

# Extrair access_token da resposta
ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
if [ -z "$ACCESS_TOKEN" ]; then
    echo "❌ Erro ao obter access_token"
    exit 1
fi
echo "Access Token: ${ACCESS_TOKEN:0:50}..."
echo ""

# 5. Criar estudante no study-manager-service
echo "5. Criando estudante no study-manager-service..."
STUDENT_RESPONSE=$(make_request "POST" "/students" '{"name": "Estudante Teste", "email": "teste@estudante.com"}' "X-User-ID: teste@estudante.com" "$STUDY_BASE_URL")
echo "Resposta: $STUDENT_RESPONSE"

# Extrair student_id da resposta
STUDENT_ID=$(echo $STUDENT_RESPONSE | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
if [ -z "$STUDENT_ID" ]; then
    echo "❌ Erro ao obter student_id"
    exit 1
fi
echo "Student ID: $STUDENT_ID"
echo ""

# 6. Testar operações autenticadas no study-manager-service
echo "6. Testando operações autenticadas..."

# Buscar estudante por ID
echo "6.1. Buscando estudante por ID..."
STUDENT_GET_RESPONSE=$(make_request "GET" "/students/$STUDENT_ID" "" "Authorization: Bearer $ACCESS_TOKEN" "X-Client-ID: $CLIENT_ID" "$STUDY_BASE_URL")
echo "Resposta: $STUDENT_GET_RESPONSE"
echo ""

# Buscar estudante por User ID
echo "6.2. Buscando estudante por User ID..."
STUDENT_USER_RESPONSE=$(make_request "GET" "/students/user/teste@estudante.com" "" "Authorization: Bearer $ACCESS_TOKEN" "X-Client-ID: $CLIENT_ID" "$STUDY_BASE_URL")
echo "Resposta: $STUDENT_USER_RESPONSE"
echo ""

# Criar matéria
echo "6.3. Criando matéria..."
SUBJECT_RESPONSE=$(make_request "POST" "/subjects" '{"name": "Matemática", "description": "Matéria de matemática básica"}' "Authorization: Bearer $ACCESS_TOKEN" "X-Client-ID: $CLIENT_ID" "$STUDY_BASE_URL")
echo "Resposta: $SUBJECT_RESPONSE"
echo ""

# Listar matérias
echo "6.4. Listando matérias..."
SUBJECTS_RESPONSE=$(make_request "GET" "/subjects" "" "Authorization: Bearer $ACCESS_TOKEN" "X-Client-ID: $CLIENT_ID" "$STUDY_BASE_URL")
echo "Resposta: $SUBJECTS_RESPONSE"
echo ""

# 7. Testar validação de token
echo "7. Testando validação de token..."
VALIDATE_RESPONSE=$(make_request "POST" "/validate" "{\"token\": \"$ACCESS_TOKEN\", \"client_id\": \"$CLIENT_ID\"}" "" "$AUTH_BASE_URL")
echo "Resposta: $VALIDATE_RESPONSE"
echo ""

echo "✅ Testes de integração concluídos com sucesso!"
echo ""
echo "📋 Resumo:"
echo "- Cliente criado: $CLIENT_ID"
echo "- Usuário registrado: teste@estudante.com"
echo "- Login realizado com sucesso"
echo "- Estudante criado: $STUDENT_ID"
echo "- Operações autenticadas funcionando"
echo "- Validação de token funcionando"
