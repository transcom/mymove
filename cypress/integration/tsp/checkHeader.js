/* global cy */
describe('TSP User Checks Shipment Info Header', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('tsp user sees header info', function() {
    tspUserViewsHeaderInfo();
  });
});

function tspUserViewsHeaderInfo() {
  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find a shipment and open it
  cy
    .get('div')
    .contains('BACON1')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  // Check the move type and code
  cy.contains('MOVE INFO — HHG CODE D');

  // Check the name is correct
  cy.get('h1').contains('Submitted, HHG');

  // Check the info bar
  cy
    .get('ul')
    .contains('li', 'GBL# LKBM7000002')
    .parentsUntil('div')
    .contains('li', 'Locator# BACON1')
    .parentsUntil('div')
    .contains('li', 'LKBM to LKBM')
    .parentsUntil('div')
    .contains('li', 'DoD ID# 4444567890')
    .parentsUntil('div')
    .contains('li', 'Status: Awarded')
    .parentsUntil('div')
    .contains('li', '555-555-5555');
}
