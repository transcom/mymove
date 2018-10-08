import { testPremoveSurvey } from '../../support/testPremoveSurvey';

/* global cy */
describe('TSP User Ships a Shipment', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('tsp user Picks Up a shipment', function() {
    tspUserPicksUpShipment();
    tspUserDeliversShipment();
  });
});

function tspUserPicksUpShipment() {
  // Open approved shipments queue
  cy
    .get('div')
    .contains('Approved Shipments')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/approved/);
  });

  // Find shipment and open it
  cy
    .get('div')
    .contains('SHIPME')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  // Click the Transport button
  cy
    .get('div')
    .contains('Enter Pickup')
    .click();

  // Done button should be disabled.
  cy
    .get('button')
    .contains('Done')
    .should('be.disabled');

  // Pick a date!
  cy
    .get('div')
    .contains('Actual Pickup Date')
    .get('input')
    .click();

  cy
    .get('div')
    .contains('11')
    .click();

  // Cancel
  cy
    .get('button')
    .contains('Cancel')
    .click();

  // Wash, Rinse, Repeat
  // Click the Transport button
  cy
    .get('div')
    .contains('Enter Pickup')
    .click();

  // Done button should be disabled.
  cy
    .get('button')
    .contains('Done')
    .should('be.disabled');

  // Pick a date!
  cy
    .get('div')
    .contains('Actual Pickup Date')
    .get('input')
    .click();

  cy
    .get('div')
    .contains('11')
    .click();

  cy
    .get('button')
    .contains('Done')
    .click();

  // Appears in dates panel
  cy.get('div.actual_pickup_date').contains('11');

  // New status
  cy.get('li').contains('In_transit');
}

function tspUserDeliversShipment() {
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  // Click the Transport button
  cy
    .get('div')
    .contains('Enter Delivery')
    .click();

  // Done button should be disabled.
  cy
    .get('button')
    .contains('Done')
    .should('be.disabled');

  // Pick a date!
  cy
    .get('div')
    .contains('Actual Delivery Date')
    .get('input')
    .click();

  cy
    .get('div')
    .contains('13')
    .click();

  // Cancel
  cy
    .get('button')
    .contains('Cancel')
    .click();

  // Wash, Rinse, Repeat
  // Click the Transport button
  cy
    .get('div')
    .contains('Enter Delivery')
    .click();

  // Done button should be disabled.
  cy
    .get('button')
    .contains('Done')
    .should('be.disabled');

  // Pick a date!
  cy
    .get('div')
    .contains('Actual Delivery Date')
    .get('input')
    .click();

  cy
    .get('div')
    .contains('13')
    .click();

  cy
    .get('button')
    .contains('Done')
    .click();

  // New status
  cy.get('li').contains('Delivered');
}
