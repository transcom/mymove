/* global cy */
describe('TSP User generates GBL', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });

  it('tsp user generates GBL from shipment info page', function() {
    tspUserGeneratesGBL();
  });

  it('tsp user can open a GBL from the shipment info page', function() {
    tspUserViewsGBL();
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

  cy
    .get('button')
    .contains('Generate Bill of Lading')
    .should('be.disabled');

  // I have seen this take anywhere from 8s - 18s. Until we optimize it, giving the test a long
  // timeout.
  cy.get('.usa-alert-success', { timeout: 20000 }).contains('GBL has been created');

  cy
    .get('button')
    .contains('View GBL')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Generate Bill of Lading')
    .click();

  cy
    .get('button')
    .contains('Generate Bill of Lading')
    .should('be.disabled');

  cy.get('.usa-alert-warning').contains('There is already a GBL for this shipment. ');
}

function tspUserViewsGBL() {
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

  cy.get('.documents').should($div => expect($div.text()).to.contain('Government Bill Of Lading'));

  cy
    .get('.documents')
    .get('a')
    .contains('Government Bill Of Lading')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+\/documents\/[^/]+/);
  });
}
