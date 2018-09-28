/* global cy */
describe('TSP User generates GBL', function() {
  beforeEach(() => {
    cy.signIntoTSP();
    // cy.resetDb()
  });
  it('tsp user generates GBL from shipment info page', function() {
    tspUserGeneratesGBL();
  });
});

function tspUserGeneratesGBL() {
  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find shipment
  cy
    .get('div')
    .contains('GBLGBL')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  cy
    .get('button')
    .contains('Generate Bill of Lading')
    .should('be.enabled');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  // If clicked too soon, there's a server error
  cy.wait(500);
  cy
    .get('button')
    .contains('Generate Bill of Lading')
    .click();

  cy.get('.usa-alert-success').contains('GBL generated successfully.');

  cy
    .get('button')
    .contains('Generate Bill of Lading')
    .click();

  cy
    .get('.usa-alert-warning')
    .contains('There is already a GBL for this shipment. ');
}
