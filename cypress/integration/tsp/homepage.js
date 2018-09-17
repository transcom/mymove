/* global cy */
describe('TSP Home Page', function() {
  beforeEach(() => {
    Cypress.config('baseUrl', 'http://tsplocal:8080');
  });
  it('successfully loads when not logged in', function() {
    // Logs out any users
    cy.clearCookies();
    cy.visit('/');
    cy.contains('tsp.move.mil');
    cy.contains('Sign In');
  });
});
