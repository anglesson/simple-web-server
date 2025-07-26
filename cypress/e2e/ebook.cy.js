describe('Ebook Management', () => {
  it('should redirect to login when accessing ebook list without authentication', () => {
    cy.visit('/ebook')
    cy.url().should('include', '/login')
  })

  it('should redirect to login when accessing create ebook without authentication', () => {
    cy.visit('/ebook/create')
    cy.url().should('include', '/login')
  })

  it('should redirect to login when accessing edit ebook without authentication', () => {
    cy.visit('/ebook/edit/1')
    cy.url().should('include', '/login')
  })

  it('should redirect to login when accessing view ebook without authentication', () => {
    cy.visit('/ebook/view/1')
    cy.url().should('include', '/login')
  })
}) 