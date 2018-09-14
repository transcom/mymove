/* global cy */
describe('TSP User Views Shipment', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('tsp user views shipments in queue new shipments', function() {
    tspUserViewsShipments();
  });
  it('tsp user views shipments in queue new shipments', function() {
    tspUserViewsApprovedShipments();
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

function tspUserViewsApprovedShipments() {
  // Open accepted shipments queue
  cy
    .get('div')
    .contains('Approved Shipments')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/approved/);
  });

  // Find shipment
  cy.get('div').contains('APPRVD');
  cy
    .get('div')
    .contains('BACON1')
    .should('not.exist');
}
