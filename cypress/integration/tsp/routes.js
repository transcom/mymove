/* global cy, Cypress */
describe('TSP User Navigating the App', function() {
  it('unauthorized user tries to access an authorized route', function() {
    unauthorizedTspUserGoesToAuthorizedRoute();
  });
  it('tsp user navigates to an invalid shipment', function() {
    tspUserViewsInvalidShipment();
  });
});

function unauthorizedTspUserGoesToAuthorizedRoute() {
  Cypress.config('baseUrl', 'http://tsplocal:4000');
  cy.logout();
  cy.visit('/queues/new');
  cy.contains('Welcome to tsp.move.mil');
  cy.contains('Sign In');
}

function tspUserViewsInvalidShipment() {
  cy.signIntoTSP();
  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // visit an invalid url
  cy.visit('/shipments/some-invalid-uuid');

  // redirected to the queues page due to invalid shipment
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });
}
