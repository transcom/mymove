/* global cy */
describe('tsp user', function() {
  beforeEach(() => {
    Cypress.config('baseUrl', 'http://tsplocal:4000');
    cy.signInAsUser('6cd03e5b-bee8-4e97-a340-fecb8f3d5465');
  });
});
