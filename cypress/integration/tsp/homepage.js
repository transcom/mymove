/* global cy, Cypress */
describe('TSP Home Page', function() {
  beforeEach(() => {
    Cypress.config('baseUrl', 'http://tsplocal:4000');
  });
  it('successfully loads when not logged in', function() {
    // Logs out any users
    cy.logout();
    cy.visit('/');
    cy.contains('tsp.move.mil');
    cy.contains('Sign In');
  });
});
