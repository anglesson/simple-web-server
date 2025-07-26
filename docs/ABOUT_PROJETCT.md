# Sistema de Gestão para Infoprodutores

## Visão Geral

Este sistema é uma solução completa desenvolvida para **Infoprodutores** - profissionais que criam apostilas digitais, também conhecidas como **ebooks**. O sistema oferece uma plataforma completa para gerenciamento, distribuição e monetização de conteúdo digital.

## Características Principais

### 🎯 Público-Alvo
- **Infoprodutores**: Criadores de conteúdo digital (apostilas, ebooks, cursos)
- **Clientes**: Compradores dos produtos digitais

### 📚 Gestão de Ebooks
- Criação e edição de ebooks com múltiplos arquivos PDF
- Sistema de preços e valores
- Upload e armazenamento de arquivos
- Controle de status (ativo/inativo)
- Imagens de capa personalizadas

### 👥 Gestão de Clientes
- Cadastro e validação de clientes
- Importação em lote via CSV
- Validação de CPF através da Receita Federal
- Histórico de compras e downloads
- Sistema de relacionamento cliente-criador

### 💳 Sistema de Pagamentos
- Integração com **Stripe** para processamento de pagamentos
- Checkout seguro para compras
- Webhooks para sincronização de status
- Sistema de assinaturas

### 🔒 Proteção de Conteúdo
- Sistema de **watermark** para proteção de PDFs
- Controle de downloads por compra
- Limite de downloads configurável
- Expiração de acesso por tempo
- Logs de download para auditoria

### 📊 Dashboard e Relatórios
- Visão geral de vendas e performance
- Estatísticas de downloads
- Relatórios de clientes
- Análise de receita

## Arquitetura do Sistema

### Tecnologias Utilizadas
- **Backend**: Go (Golang)
- **Framework Web**: Chi Router
- **ORM**: GORM
- **Banco de Dados**: SQLite (desenvolvimento) / PostgreSQL (produção)
- **Frontend**: HTML + Bootstrap 5 + JavaScript
- **Pagamentos**: Stripe
- **Armazenamento**: S3 (AWS)
- **Email**: GoMail

### Estrutura do Projeto

```
SimpleWebServer/
├── cmd/web/                    # Ponto de entrada da aplicação
├── internal/
│   ├── config/                 # Configurações do sistema
│   ├── handler/                # Controladores HTTP
│   ├── models/                 # Modelos de dados
│   ├── repository/             # Camada de acesso a dados
│   └── service/                # Lógica de negócio
├── pkg/                        # Bibliotecas externas
├── web/                        # Frontend (templates, assets)
└── docs/                       # Documentação
```

## Funcionalidades Detalhadas

### 1. Gestão de Criadores (Infoprodutores)

#### Cadastro e Validação
- Registro com dados pessoais completos
- Validação de CPF na Receita Federal
- Verificação de maioridade (18+ anos)
- Sistema de autenticação seguro

#### Perfil e Configurações
- Edição de dados pessoais
- Configurações de conta
- Gestão de assinatura
- Configurações de pagamento

### 2. Gestão de Ebooks

#### Criação de Ebooks
```go
type Ebook struct {
    Title       string  // Título do ebook
    Description string  // Descrição detalhada
    Value       float64 // Preço em reais
    Status      bool    // Ativo/Inativo
    Image       string  // Imagem de capa
    File        string  // Arquivo PDF principal
    FileURL     string  // URL do arquivo
    CreatorID   uint    // ID do criador
}
```

#### Funcionalidades
- Upload de arquivos PDF
- Definição de preços
- Configuração de status
- Upload de imagens de capa
- Edição e atualização de conteúdo

### 3. Gestão de Clientes

#### Cadastro de Clientes
```go
type Client struct {
    Name      string     // Nome completo
    CPF       string     // CPF único
    Birthdate string     // Data de nascimento
    Email     string     // Email de contato
    Phone     string     // Telefone
    Validated bool       // Status de validação
    Creators  []*Creator // Relacionamento com criadores
    Purchases []*Purchase // Histórico de compras
}
```

#### Funcionalidades
- Cadastro individual de clientes
- Importação em lote via CSV
- Validação automática de CPF
- Gestão de relacionamentos
- Histórico completo de compras

### 4. Sistema de Compras e Downloads

#### Modelo de Compra
```go
type Purchase struct {
    EbookID       uint      // ID do ebook comprado
    ClientID      uint      // ID do cliente
    ExpiresAt     time.Time // Data de expiração
    DownloadsUsed int       // Downloads realizados
    DownloadLimit int       // Limite de downloads (-1 = ilimitado)
    Downloads     []DownloadLog // Log de downloads
}
```

#### Controle de Acesso
- Limite configurável de downloads
- Expiração por tempo
- Logs detalhados de acesso
- Sistema de watermark para proteção

### 5. Sistema de Pagamentos

#### Integração Stripe
- Checkout seguro
- Processamento de cartões
- Webhooks para sincronização
- Sistema de assinaturas recorrentes
- Relatórios de transações

### 6. Proteção de Conteúdo

#### Watermark
- Marcação automática de PDFs
- Inclusão de dados do comprador
- Timestamp de download
- Proteção contra redistribuição

#### Controle de Acesso
- URLs temporárias para download
- Verificação de permissões
- Logs de auditoria
- Sistema de expiração

## Fluxos Principais

### 1. Fluxo de Venda
1. Cliente acessa página do ebook
2. Realiza pagamento via Stripe
3. Sistema cria registro de compra
4. Cliente recebe link de download
5. Sistema aplica watermark no PDF
6. Download é registrado no log

### 2. Fluxo de Criação de Ebook
1. Criador faz login no sistema
2. Acessa área de criação de ebooks
3. Preenche informações do produto
4. Faz upload do arquivo PDF
5. Define preço e configurações
6. Ativa o ebook para venda

### 3. Fluxo de Gestão de Clientes
1. Criador acessa área de clientes
2. Cadastra cliente individual ou importa lista
3. Sistema valida CPF automaticamente
4. Cliente é associado ao criador
5. Histórico de compras é mantido

## Segurança e Compliance

### Validação de Dados
- Validação de CPF na Receita Federal
- Verificação de maioridade
- Validação de emails
- Sanitização de dados

### Proteção de Conteúdo
- Watermark automático
- URLs temporárias
- Controle de downloads
- Logs de auditoria

### Segurança de Pagamentos
- Integração PCI-compliant (Stripe)
- Criptografia de dados sensíveis
- Webhooks seguros
- Validação de transações

## Configuração e Deploy

### Requisitos
- Go 1.19+
- SQLite (dev) / PostgreSQL (prod)
- Conta Stripe
- Bucket S3 (opcional)

### Variáveis de Ambiente
```bash
# Banco de Dados
DATABASE_URL=postgres://user:pass@localhost/dbname

# Stripe
STRIPE_SECRET_KEY=sk_test_...
STRIPE_PUBLISHABLE_KEY=pk_test_...

# S3 (opcional)
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_REGION=us-east-1
S3_BUCKET=your_bucket_name

# Email
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your_email@gmail.com
SMTP_PASS=your_password
```

### Comandos Principais
```bash
# Executar servidor
make run

# Executar testes
make test

# Build para produção
make build
```

## Roadmap e Melhorias Futuras

### Funcionalidades Planejadas
- [ ] Sistema de cupons de desconto
- [ ] Relatórios avançados de analytics
- [ ] Integração com mais gateways de pagamento
- [ ] API REST para integrações
- [ ] Sistema de afiliados
- [ ] Notificações push
- [ ] App mobile para criadores

### Melhorias Técnicas
- [ ] Cache Redis para performance
- [ ] CDN para arquivos estáticos
- [ ] Sistema de backup automático
- [ ] Monitoramento e alertas
- [ ] Testes de carga
- [ ] Documentação da API

## Suporte e Contribuição

### Como Contribuir
1. Fork do repositório
2. Crie uma branch para sua feature
3. Implemente as mudanças
4. Adicione testes
5. Submeta um pull request

### Padrões de Código
- Seguir convenções Go
- Usar Conventional Commits
- Manter cobertura de testes alta
- Documentar APIs e funções complexas

### Testes
- Testes unitários para todos os serviços
- Testes de integração para handlers
- Testes E2E com Cypress
- Cobertura mínima de 80%

---

**Desenvolvido com ❤️ para infoprodutores brasileiros**
