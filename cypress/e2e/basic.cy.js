describe('Basic Application Test', () => {
  it('should load the home page', () => {
    cy.visit('/')
    cy.get('body').should('be.visible')
  })

  it('should load login page', () => {
    cy.visit('/login')
    cy.get('input[name="email"]').should('be.visible')
    cy.get('input[name="password"]').should('be.visible')
  })

  it('should load register page', () => {
    cy.visit('/register')
    cy.get('input[name="name"]').should('be.visible')
    cy.get('input[name="email"]').should('be.visible')
  })

  it('should load forget password page', () => {
    cy.visit('/forget-password')
    cy.get('input[name="email"]').should('be.visible')
  })

  it('should redirect to login when accessing dashboard', () => {
    cy.visit('/dashboard')
    cy.url().should('include', '/login')
  })
}) 