# DashUI Integration Guide

## Visão Geral

Este projeto utiliza o **DashUI** como template base para a interface administrativa. O DashUI é um template moderno baseado em Bootstrap 5 que oferece componentes pré-estilizados e uma experiência de usuário consistente.

## Estrutura de Arquivos

```
web/
├── assets/
│   ├── css/
│   │   └── theme.min.css          # CSS principal do DashUI
│   ├── js/
│   │   └── theme.min.js           # JavaScript principal do DashUI
│   ├── images/
│   │   ├── svg/                   # Ícones SVG
│   │   └── placeholder/           # Imagens placeholder
│   └── libs/                      # Bibliotecas externas
├── layouts/
│   ├── admin.html                 # Layout administrativo
│   └── guest.html                 # Layout público
├── pages/
│   ├── dashboard.html             # Dashboard principal
│   ├── client/
│   │   ├── create.html            # Formulário de criação
│   │   └── update.html            # Formulário de edição
│   └── ...
├── partials/
│   ├── table_footer.html          # Paginação
│   ├── empty_table.html           # Estado vazio
│   └── ...
└── mails/
    └── ...                        # Templates de email
```

## Componentes Principais

### 1. Layouts

#### Admin Layout (`admin.html`)
- Sidebar de navegação
- Navbar superior
- Área de conteúdo principal
- Sistema de notificações toast

#### Guest Layout (`guest.html`)
- Layout limpo para páginas públicas
- Sem sidebar
- Foco no conteúdo principal

### 2. Estrutura de Página Padrão

```html
{{ define "title" }} Título da Página {{ end }}
{{ define "content" }}
<div class="container-fluid p-6">
  <!-- Page Header -->
  <div class="row">
    <div class="col-lg-12 col-md-12 col-12">
      <div class="border-bottom pb-4 mb-4">
        <div class="row align-items-center">
          <div class="col">
            <h3 class="mb-0 fw-bold">Título</h3>
            <p class="mb-0 text-muted">Descrição</p>
          </div>
          <div class="col-auto">
            <!-- Ações principais -->
          </div>
        </div>
      </div>
    </div>
  </div>
  
  <!-- Content Area -->
  <div class="py-6">
    <div class="row">
      <div class="col-xl-12 col-lg-12 col-md-12 col-12">
        <div class="card h-100">
          <!-- Conteúdo do card -->
        </div>
      </div>
    </div>
  </div>
</div>
{{end}}
```

### 3. Componentes Reutilizáveis

#### Cards
```html
<div class="card h-100">
  <div class="card-header bg-white py-4">
    <!-- Header do card -->
  </div>
  <div class="card-body">
    <!-- Conteúdo do card -->
  </div>
  <div class="card-footer">
    <!-- Footer do card -->
  </div>
</div>
```

#### Tabelas
```html
<div class="table-responsive">
  <table class="table table-hover text-nowrap">
    <thead class="table-light">
      <tr>
        <th scope="col" class="border-0">Coluna</th>
      </tr>
    </thead>
    <tbody>
      <!-- Linhas da tabela -->
    </tbody>
  </table>
</div>
```

#### Avatares
```html
<div class="avatar avatar-sm avatar-indicators avatar-online">
  <div class="avatar-content bg-primary">
    <span class="text-white fw-bold">{{ .GetInitials }}</span>
  </div>
</div>
```

#### Badges
```html
<span class="badge bg-success-subtle text-success">
  <i data-feather="check-circle" class="icon-xs me-1"></i>
  Ativo
</span>
```

#### Botões
```html
<a href="#" class="btn btn-primary">
  <i data-feather="plus" class="icon-xs me-2"></i>
  Adicionar
</a>
```

#### Dropdowns
```html
<div class="dropdown dropstart">
  <a class="btn btn-sm btn-icon btn-ghost-secondary rounded-circle" href="#" data-bs-toggle="dropdown">
    <i data-feather="more-vertical" class="icon-xs"></i>
  </a>
  <div class="dropdown-menu">
    <a class="dropdown-item" href="#">
      <i data-feather="edit-2" class="icon-xs me-2"></i>
      Editar
    </a>
  </div>
</div>
```

## Sistema de Cores

### Cores Primárias
- **Primary**: `btn-primary`, `bg-primary`, `text-primary`
- **Secondary**: `btn-outline-secondary`, `bg-light`
- **Success**: `badge bg-success-subtle text-success`
- **Danger**: `badge bg-danger-subtle text-danger`
- **Warning**: `badge bg-warning-subtle text-warning`
- **Info**: `badge bg-info-subtle text-info`

### Classes Utilitárias
- **Spacing**: `p-6`, `py-4`, `mb-4`, `me-2`
- **Typography**: `fw-bold`, `fw-semi-bold`, `text-muted`, `fs-6`
- **Layout**: `d-flex`, `align-items-center`, `justify-content-between`

## Ícones (Feather Icons)

### Uso Básico
```html
<i data-feather="plus" class="icon-xs me-2"></i>
```

### Tamanhos Disponíveis
- `icon-xs`: 12px
- `icon-sm`: 16px
- `icon-md`: 20px
- `icon-lg`: 24px

### Ícones Comuns
- `plus`: Adicionar
- `edit-2`: Editar
- `trash-2`: Excluir
- `eye`: Visualizar
- `search`: Buscar
- `upload`: Upload
- `download`: Download
- `check-circle`: Sucesso
- `x-circle`: Erro
- `more-vertical`: Menu de ações

## Integração com Backend

### Renderização de Templates
```go
template.View(w, r, "page_name", map[string]any{
    "Data": data,
    "Pagination": pagination,
}, "admin")
```

### Estrutura de Dados
```go
// No handler
data := map[string]any{
    "Clients": clients,
    "Pagination": pagination,
    "SearchTerm": term,
}

// No template
{{ range .Clients }}
  {{ .Name }}
{{ end }}
```

### Paginação
```go
pagination := models.NewPagination(page, perPage)
pagination.SetTotal(totalCount)
```

### Flash Messages
```go
flashMessage := ch.flashMessageFactory(w, r)
flashMessage.Success("Operação realizada com sucesso!")
flashMessage.Error("Erro na operação")
```

## Padrões de UX

### Estados da Interface

#### Estado Vazio
```html
{{ template "empty-table" . }}
```

#### Loading State
```html
<div class="d-flex justify-content-center">
  <div class="spinner-border" role="status">
    <span class="visually-hidden">Carregando...</span>
  </div>
</div>
```

#### Error State
```html
<div class="alert alert-danger" role="alert">
  <i data-feather="alert-circle" class="icon-xs me-2"></i>
  Erro ao carregar dados
</div>
```

### Interatividade

#### Busca com Debounce
```javascript
let searchTimeout;
searchInput.addEventListener("input", function () {
  clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    if (searchInput.value === "") {
      history.replaceState(null, "", window.location.pathname);
      document.location.reload();
    } else {
      searchInput.form.submit();
    }
  }, 500);
});
```

#### Confirmação de Ações
```javascript
function confirmDelete(id) {
  if (confirm('Tem certeza que deseja excluir este item?')) {
    // Executar exclusão
  }
}
```

## Responsividade

### Breakpoints Bootstrap 5
- **xs**: < 576px
- **sm**: ≥ 576px
- **md**: ≥ 768px
- **lg**: ≥ 992px
- **xl**: ≥ 1200px
- **xxl**: ≥ 1400px

### Classes Responsivas
```html
<div class="col-lg-8 col-md-6 col-12">
  <!-- Conteúdo responsivo -->
</div>
```

## Performance

### Otimizações Recomendadas
1. **Lazy Loading**: Carregar imagens sob demanda
2. **Minificação**: Usar arquivos CSS/JS minificados
3. **Caching**: Implementar cache para assets estáticos
4. **Pagination**: Usar paginação para grandes listas
5. **Debounce**: Implementar debounce em buscas

### Boas Práticas
- Usar `table-responsive` para tabelas grandes
- Implementar loading states para operações assíncronas
- Otimizar queries do banco de dados
- Usar índices adequados para buscas

## Testes

### Cypress E2E
```javascript
describe('Client List', () => {
  it('should display clients correctly', () => {
    cy.visit('/client')
    cy.get('.table tbody tr').should('have.length.greaterThan', 0)
  })
})
```

### Testes de Responsividade
- Testar em diferentes tamanhos de tela
- Verificar navegação mobile
- Validar acessibilidade

## Troubleshooting

### Problemas Comuns

#### Ícones não aparecem
```javascript
// Inicializar Feather Icons
feather.replace()
```

#### Tabela não responsiva
```html
<div class="table-responsive">
  <table class="table">
    <!-- Conteúdo da tabela -->
  </table>
</div>
```

#### Dropdown não funciona
```html
<!-- Verificar se Bootstrap JS está carregado -->
<script src="/assets/libs/bootstrap/dist/js/bootstrap.bundle.min.js"></script>
```

### Debug
```javascript
// Verificar se DashUI está carregado
console.log('DashUI loaded:', typeof DashUI !== 'undefined')
```

## Recursos Adicionais

### Documentação
- [Bootstrap 5 Documentation](https://getbootstrap.com/docs/5.3/)
- [Feather Icons](https://feathericons.com/)
- [DashUI Template](https://dashui.gethugothemes.com/)

### Ferramentas
- Bootstrap Icons
- Feather Icons
- Prism.js (syntax highlighting)
- ApexCharts (gráficos)

### Customização
- Modificar `theme.min.css` para customizações
- Adicionar classes utilitárias personalizadas
- Criar componentes customizados seguindo o padrão 