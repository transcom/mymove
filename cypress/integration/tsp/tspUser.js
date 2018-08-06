/* global cy */
describe('tsp user', function() {
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

  // Find shipment (requires HHG Move)
  /*
  cy
    .get('div')
    .contains('VGHEIS');
    // TODO: (2018_08_01 cgilmer) Open shipment
  */
}
