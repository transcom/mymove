/* global cy */
describe('TSP User Views Shipment', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('tsp user views shipments in queue new shipments', function() {
    tspUserViewsShipments();
  });
});

function tspUserViewsShipments() {
  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find shipment
  cy.get('div').contains('BACON1');
}
