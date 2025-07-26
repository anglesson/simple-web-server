describe('Client Management', () => {
  it('should redirect to login when accessing client list without authentication', () => {
    cy.visit('/client')
    cy.url().should('include', '/login')
  })

  it('should redirect to login when accessing create client without authentication', () => {
    cy.visit('/client/new')
    cy.url().should('include', '/login')
  })

  it('should redirect to login when accessing edit client without authentication', () => {
    cy.visit('/client/update/1')
    cy.url().should('include', '/login')
  })
}) 