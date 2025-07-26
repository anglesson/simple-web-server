// ***********************************************
// This example commands.js shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************

// Comandos básicos customizados
Cypress.Commands.add('getByData', (selector) => {
  return cy.get(`[data-test="${selector}"]`)
})

Cypress.Commands.add('shouldBeVisible', (selector) => {
  return cy.get(selector).should('be.visible')
})

Cypress.Commands.add('shouldNotBeVisible', (selector) => {
  return cy.get(selector).should('not.be.visible')
})

Cypress.Commands.add('fillForm', (formData) => {
  Object.keys(formData).forEach(field => {
    cy.get(`[name="${field}"]`).type(formData[field])
  })
})

Cypress.Commands.add('submitForm', () => {
  cy.get('form').submit()
})

Cypress.Commands.add('waitForPageLoad', () => {
  cy.get('body').should('be.visible')
}) 