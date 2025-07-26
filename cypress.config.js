const { defineConfig } = require('cypress')

module.exports = defineConfig({
  e2e: {
    baseUrl: 'http://localhost:8080',
    viewportWidth: 1280,
    viewportHeight: 720,
    video: false,
    screenshotOnRunFailure: true,
    defaultCommandTimeout: 10000,
    requestTimeout: 10000,
    responseTimeout: 10000,
    setupNodeEvents(on, config) {
      // implement node event listeners here
    },
  },
  env: {
    // Variáveis de ambiente para os testes
    testUser: {
      email: 'test@example.com',
      password: 'testpassword123',
      name: 'Test User',
      cpf: '123.456.789-00',
      phone: '(11) 99999-9999'
    }
  }
}) 