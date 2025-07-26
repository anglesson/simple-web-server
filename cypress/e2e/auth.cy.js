describe('Authentication', () => {
  beforeEach(() => {
    cy.visit('/')
  })

  it('should display login page', () => {
    cy.visit('/login')
    cy.get('input[name="email"]').should('be.visible')
    cy.get('input[name="password"]').should('be.visible')
    cy.get('button[type="submit"]').should('be.visible')
    cy.get('button[type="submit"]').should('contain', 'Entrar')
  })

  it('should display register page', () => {
    cy.visit('/register')
    cy.get('input[name="name"]').should('be.visible')
    cy.get('input[name="email"]').should('be.visible')
    cy.get('input[name="password"]').should('be.visible')
    cy.get('button[type="submit"]').should('contain', 'Criar conta')
  })

  it('should display forget password page', () => {
    cy.visit('/forget-password')
    cy.get('input[name="email"]').should('be.visible')
    cy.get('button[type="submit"]').should('contain', 'Redefinir senha')
  })

  it('should redirect to login when accessing protected route', () => {
    cy.visit('/dashboard')
    cy.url().should('include', '/login')
  })

  it('should show error with invalid credentials', () => {
    cy.visit('/login')
    cy.get('input[name="email"]').type('invalid@example.com')
    cy.get('input[name="password"]').type('wrongpassword')
    cy.get('form').submit()
    
    // Verificar se há alguma mensagem de erro ou permanece na página de login
    cy.url().should('include', '/login')
  })

  it('should navigate between auth pages', () => {
    cy.visit('/login')
    cy.get('a[href="/register"]').click()
    cy.url().should('include', '/register')
    
    cy.get('a[href="/login"]').click()
    cy.url().should('include', '/login')
    
    cy.get('a[href="/forget-password"]').click()
    cy.url().should('include', '/forget-password')
  })
}) 