/* global cy */
describe('TSP User On Queues Page', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('tsp user navigates to an invalid shipment', function() {
    tspUserViewsInvalidShipment();
  });
});

function tspUserViewsInvalidShipment() {
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
