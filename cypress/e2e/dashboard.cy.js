describe('Dashboard', () => {
  it('should redirect to login when accessing dashboard without authentication', () => {
    cy.visit('/dashboard')
    cy.url().should('include', '/login')
  })

  it('should redirect to login when accessing ebook page without authentication', () => {
    cy.visit('/ebook')
    cy.url().should('include', '/login')
  })

  it('should redirect to login when accessing client page without authentication', () => {
    cy.visit('/client')
    cy.url().should('include', '/login')
  })

  it('should redirect to login when accessing settings page without authentication', () => {
    cy.visit('/settings')
    cy.url().should('include', '/login')
  })

  it('should redirect to login when accessing send page without authentication', () => {
    cy.visit('/send')
    cy.url().should('include', '/login')
  })
}) 