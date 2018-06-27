/* global cy */
describe('The Home Page', function() {
  it('successfully loads when not logged in', function() {
    cy.visit('/');
    cy.contains('Welcome');
    cy.contains('Sign In');
  });
});
