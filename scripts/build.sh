#!/bin/bash

# Script de build para o projeto SimpleWebServer
set -e

echo "🚀 Iniciando build do SimpleWebServer..."

# Criar diretório bin se não existir
mkdir -p bin

# Compilar o projeto
echo "📦 Compilando o projeto..."
GOOS=linux GOARCH=amd64 go build -o bin/simple-web-server cmd/web/main.go

# Tornar o binário executável
chmod +x bin/simple-web-server

echo "✅ Build concluído! Binário criado em: bin/simple-web-server"
echo "📊 Tamanho do binário: $(du -h bin/simple-web-server | cut -f1)" 