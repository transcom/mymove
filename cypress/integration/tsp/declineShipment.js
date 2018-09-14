import { testPremoveSurvey } from '../../support/testPremoveSurvey';

/* global cy */
describe('TSP User Rejects a Shipment', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('tsp user rejects a shipment', function() {
    tspUserRejectsShipment();
  });
});

function tspUserRejectsShipment() {
  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find shipment and open it
  cy
    .get('div')
    .contains('REJECT')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/shipments\/[^/]+/);
  });

  // Click the Reject button
  cy
    .get('.usa-button-secondary')
    .contains('Reject Shipment')
    .click();

  // New button should be disabled.
  cy
    .get('button')
    .contains('Reject Shipment')
    .should('be.disabled');

  // Give a reason
  cy.get('textarea').type('End to End test.');
  cy
    .get('button')
    .contains('Reject Shipment')
    .click();

  // Cancel
  cy
    .get('div')
    .contains('No, never mind')
    .click();

  // Wash, Rinse, Repeat
  // Click the Reject button
  cy
    .get('button')
    .contains('Reject Shipment')
    .click();

  // New button should be disabled.
  cy
    .get('button')
    .contains('Reject Shipment')
    .should('be.disabled');

  // Give a reason
  cy.get('textarea').type('End to End test.');
  cy
    .get('button')
    .contains('Reject Shipment')
    .click();

  // cy
  //   .get('button')
  //   .contains('Reject Shipment')
  //   .click();

  // cy.location().should(loc => {
  //   expect(loc.pathname).to.match(/^\/queues\/new/);
  // });

  // cy
  //   .get('div')
  //   .contains('REJECT')
  //   .should('not.exist');
}
