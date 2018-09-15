/* global cy */
describe('TSP Home Page', function() {
  beforeEach(() => {
    Cypress.config('baseUrl', 'http://tsplocal:8080');
  });
  it('successfully loads when not logged in', function() {
    cy.visit('/');
    cy.contains('tsp.move.mil');
    cy.contains('Sign In');
  });
});
