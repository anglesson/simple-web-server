#!/bin/bash

# Script para configurar o arquivo .env
echo "🔧 Configurando arquivo .env..."

# Verificar se o arquivo .env já existe
if [ -f ".env" ]; then
    echo "⚠️  Arquivo .env já existe. Deseja sobrescrever? (y/N)"
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        echo "❌ Operação cancelada."
        exit 1
    fi
fi

# Copiar template para .env
cp env.template .env

echo "✅ Arquivo .env criado com sucesso!"
echo ""
echo "📝 Próximos passos:"
echo "1. Edite o arquivo .env com suas configurações"
echo "2. Configure as credenciais necessárias:"
echo "   - MAIL_USERNAME e MAIL_PASSWORD para email"
echo "   - S3_ACCESS_KEY e S3_SECRET_KEY para AWS S3"
echo "   - STRIPE_SECRET_KEY para pagamentos"
echo "   - HUB_DEVSENVOLVEDOR_TOKEN para validação de CPF"
echo ""
echo "🔒 IMPORTANTE: Nunca commite o arquivo .env no repositório!"
echo "   O arquivo .env está no .gitignore por segurança." 