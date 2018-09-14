/* global cy */
describe('TSP User Views Shipment', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('tsp user views shipments in queue new shipments', function() {
    tspUserViewsShipments();
  });
  it('tsp user views in transit hhg moves in queue HHGs In Transit', function() {
    tspUserViewsInTransitShipment();
  });
  it('tsp user views delivered hhg moves in queue HHGs Delivered', function() {
    tspUserViewsDeliveredShipment();
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

function tspUserViewsInTransitShipment() {
  // Open new moves queue
  cy.visit('/queues/in_transit');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/in_transit/);
  });

  // Find move (generated in e2ebasic.go) and open it
  cy
    .get('div')
    .contains('NINOPK')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/in_transit\/shipments\/[^/]+/);
  });
}

function tspUserViewsDeliveredShipment() {
  // Open new moves queue
  cy.visit('/queues/delivered');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/delivered/);
  });

  // Find move (generated in e2ebasic.go) and open it
  cy
    .get('div')
    .contains('SCHNOO')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/delivered\/shipments\/[^/]+/);
  });
}
