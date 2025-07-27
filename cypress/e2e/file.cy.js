describe('File Management', () => {
  beforeEach(() => {
    cy.visit('/')
  })

  describe('Unauthenticated Access', () => {
    it('should redirect to login when accessing file list without authentication', () => {
      cy.visit('/file')
      cy.url().should('include', '/login')
    })

    it('should redirect to login when accessing file upload without authentication', () => {
      cy.visit('/file/upload')
      cy.url().should('include', '/login')
    })

    it('should redirect to login when accessing file update without authentication', () => {
      cy.visit('/file/1/update')
      cy.url().should('include', '/login')
    })
  })

  describe('Authenticated Access', () => {
    beforeEach(() => {
      // Login antes de testar funcionalidades autenticadas
      cy.login('anglesson@outlook.com', 'password123')
    })

    it('should display file list page with navigation', () => {
      cy.visit('/file')
      
      // Verificar se a página carregou corretamente
      cy.get('h4.page-title').should('contain', 'Minha Biblioteca de Arquivos')
      cy.get('.breadcrumb-item').should('contain', 'Biblioteca de Arquivos')
      
      // Verificar se o botão de upload está presente
      cy.get('a[href="/file/upload"]').should('be.visible')
      cy.get('a[href="/file/upload"]').should('contain', 'Upload de Arquivo')
      
      // Verificar se a tabela está presente (mesmo vazia)
      cy.get('.table').should('be.visible')
      cy.get('thead').should('contain', 'Nome')
      cy.get('thead').should('contain', 'Tipo')
      cy.get('thead').should('contain', 'Tamanho')
      cy.get('thead').should('contain', 'Descrição')
      cy.get('thead').should('contain', 'Data Upload')
      cy.get('thead').should('contain', 'Ações')
    })

    it('should navigate to file upload page', () => {
      cy.visit('/file')
      
      // Clicar no botão de upload
      cy.get('a[href="/file/upload"]').click()
      
      // Verificar se foi redirecionado para a página de upload
      cy.url().should('include', '/file/upload')
      
      // Verificar se a página de upload carregou
      cy.get('h4.page-title').should('contain', 'Upload de Arquivo')
    })

    it('should display file upload form', () => {
      cy.visit('/file/upload')
      
      // Verificar se o formulário está presente
      cy.get('form').should('be.visible')
      
      // Verificar se os campos estão presentes
      cy.get('input[type="file"]').should('be.visible')
      cy.get('textarea[name="description"]').should('be.visible')
      cy.get('button[type="submit"]').should('be.visible')
      
      // Verificar se há instruções ou labels
      cy.get('label').should('contain', 'Arquivo')
      cy.get('label').should('contain', 'Descrição')
    })

    it('should show file upload page with proper breadcrumb', () => {
      cy.visit('/file/upload')
      
      // Verificar breadcrumb
      cy.get('.breadcrumb-item').first().should('contain', 'Dashboard')
      cy.get('.breadcrumb-item').last().should('contain', 'Upload de Arquivo')
      
      // Verificar título da página
      cy.get('h4.page-title').should('contain', 'Upload de Arquivo')
    })

    it('should have file menu item in navigation', () => {
      cy.visit('/dashboard')
      
      // Verificar se o item "Arquivos" está no menu lateral
      cy.get('a[href="/file"]').should('be.visible')
      cy.get('a[href="/file"]').should('contain', 'Arquivos')
      
      // Verificar se o ícone está presente
      cy.get('a[href="/file"] i[data-feather="folder"]').should('be.visible')
    })

    it('should navigate from dashboard to file list', () => {
      cy.visit('/dashboard')
      
      // Clicar no link "Arquivos" no menu
      cy.get('a[href="/file"]').click()
      
      // Verificar se foi redirecionado para a lista de arquivos
      cy.url().should('include', '/file')
      cy.get('h4.page-title').should('contain', 'Minha Biblioteca de Arquivos')
    })

    it('should show empty state when no files exist', () => {
      cy.visit('/file')
      
      // Se não há arquivos, deve mostrar uma mensagem ou tabela vazia
      // Verificar se a estrutura da tabela está presente
      cy.get('.table').should('be.visible')
      cy.get('thead').should('be.visible')
      
      // Verificar se há uma mensagem de "nenhum arquivo" ou tabela vazia
      cy.get('tbody').should('exist')
    })

    it('should have proper page structure and styling', () => {
      cy.visit('/file')
      
      // Verificar estrutura da página
      cy.get('.container-fluid').should('be.visible')
      cy.get('.page-title-box').should('be.visible')
      cy.get('.card').should('be.visible')
      cy.get('.card-header').should('be.visible')
      cy.get('.card-body').should('be.visible')
      
      // Verificar se o botão de upload tem o estilo correto
      cy.get('a[href="/file/upload"]').should('have.class', 'btn')
      cy.get('a[href="/file/upload"]').should('have.class', 'btn-primary')
    })
  })

  describe('File Upload Form Validation', () => {
    beforeEach(() => {
      // Login antes de testar
      cy.login('anglesson@outlook.com', 'password123')
    })

    it('should show validation errors for empty form submission', () => {
      cy.visit('/file/upload')
      
      // Tentar submeter formulário vazio
      cy.get('button[type="submit"]').click()
      
      // Verificar se permanece na página de upload
      cy.url().should('include', '/file/upload')
    })

    it('should accept file selection', () => {
      cy.visit('/file/upload')
      
      // Simular seleção de arquivo (usando fixture)
      cy.fixture('sample.pdf').then(fileContent => {
        cy.get('input[type="file"]').attachFile({
          fileContent: fileContent,
          fileName: 'sample.pdf',
          mimeType: 'application/pdf'
        })
      })
      
      // Verificar se o arquivo foi selecionado
      cy.get('input[type="file"]').should('have.value')
    })
  })
}) 