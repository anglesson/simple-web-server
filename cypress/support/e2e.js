// ***********************************************************
// This example support/e2e.js is processed and
// loaded automatically before your test files.
//
// This is a great place to put global configuration and
// behavior that modifies Cypress.
//
// You can change the location of this file or turn off
// automatically serving support files with the
// 'supportFile' configuration option.
//
// You can read more here:
// https://on.cypress.io/configuration
// ***********************************************************

// Import commands.js using ES2015 syntax:
import './commands'

// Alternatively you can use CommonJS syntax:
// require('./commands')

// Comandos customizados para autenticação
Cypress.Commands.add('login', (email, password) => {
  cy.visit('/login')
  cy.get('input[name="email"]').type(email)
  cy.get('input[name="password"]').type(password)
  cy.get('form').submit()
  // Aguardar redirecionamento para dashboard
  cy.url().should('include', '/dashboard')
})

Cypress.Commands.add('logout', () => {
  // O logout está em um dropdown, então precisamos clicar nele primeiro
  cy.get('.dropdown-toggle').click()
  cy.get('form[action="/logout"] button[type="submit"]').click()
  cy.url().should('include', '/login')
})

// Comando para limpar dados de teste
Cypress.Commands.add('cleanupTestData', () => {
  // Implementar limpeza de dados de teste se necessário
  cy.log('Cleaning up test data...')
}) 