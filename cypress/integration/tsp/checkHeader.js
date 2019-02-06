import { tspUserVerifiesShipmentStatus } from '../../support/testTspStatus';

/* global cy */
describe('TSP User Checks Shipment Info Header', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('tsp user sees header info for HHG move', function() {
    tspUserViewsHHGHeaderInfo();
  });
  it('tsp user sees header info for HHG_PPM move', function() {
    tspUserViewsHHGPPMHeaderInfo();
  });
});

function tspUserViewsHHGHeaderInfo() {
  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find a shipment and open it
  cy.selectQueueItemMoveLocator('HHGPPM');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  // Check the move type and code
  cy.contains('MOVE INFO — HHG_PPM CODE D');

  // Check the name is correct
  cy.get('div').contains('Submitted, HHGPPM');

  // Check the status
  tspUserVerifiesShipmentStatus('Shipment awarded');

  // Check the info bar
  cy
    .get('ul')
    .contains('li', 'GBL# LKNQ7123456')
    .parentsUntil('div')
    .contains('li', 'Locator# HHGPPM')
    .parentsUntil('div')
    .contains('li', 'KKFA to HAFC')
    .parentsUntil('div')
    .contains('li', 'DoD ID# 4224567890')
    .parentsUntil('div')
    .contains('li', '555-555-5555');
}

function tspUserViewsHHGPPMHeaderInfo() {
  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find a shipment and open it
  cy.selectQueueItemMoveLocator('HHGPPM');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  // Check the move type and code
  cy.contains('MOVE INFO — HHG_PPM CODE D');

  // Check the name is correct
  cy.get('div').contains('Submitted, HHG');

  // Check the status
  tspUserVerifiesShipmentStatus('Shipment awarded');

  // Check the info bar
  cy
    .get('ul')
    .contains('li', 'GBL# LKNQ7123456')
    .parentsUntil('div')
    .contains('li', 'Locator# HHGPPM')
    .parentsUntil('div')
    .contains('li', 'KKFA to HAFC')
    .parentsUntil('div')
    .contains('li', 'DoD ID# 4224567890')
    .parentsUntil('div')
    .contains('li', '555-555-5555');
}
