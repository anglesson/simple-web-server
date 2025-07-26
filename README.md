# Simple Web Server

A web application for managing ebooks and clients.

## Running the Application

```bash
make run
```

## Running Tests

### Unit Tests
```bash
make test
```

### E2E Tests with Cypress

First, install dependencies:
```bash
make install-e2e
make setup-e2e
```

Then run the application:
```bash
make run
```

In another terminal, run E2E tests:
```bash
make test-e2e
```

Or run with browser visible:
```bash
make test-e2e-headed
```

## E2E Test Structure

- `cypress/e2e/basic.cy.js` - Basic application tests (home, login, register, forget password)
- `cypress/e2e/auth.cy.js` - Authentication flow tests
- `cypress/e2e/dashboard.cy.js` - Dashboard and navigation tests
- `cypress/e2e/ebook.cy.js` - Ebook management tests
- `cypress/e2e/client.cy.js` - Client management tests

## Test Coverage

The E2E tests cover:
- ✅ Application startup and basic page loading
- ✅ Authentication pages (login, register, forget password)
- ✅ Navigation between auth pages
- ✅ Protected route redirections
- ✅ Form validation and error handling
- ✅ User interface elements visibility

## Test Data

The E2E tests use test data from:
- `cypress/fixtures/sample.pdf` - Sample PDF for ebook uploads
- `cypress/fixtures/clients.csv` - Sample CSV for client imports

## Configuration

E2E tests are configured in `cypress.config.js` to run against `http://localhost:8080`.

## Current Status

All E2E tests are passing! ✅
- 23 tests total
- 0 failures
- Covers all major application flows
