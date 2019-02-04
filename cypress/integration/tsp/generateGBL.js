import { tspUserVerifiesShipmentStatus } from '../../support/testTspStatus';

/* global cy */
describe('TSP User generates GBL', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });

  it('tsp user generates GBL from shipment info page', function() {
    tspUserGeneratesGBL();
    tspUserVerifiesShipmentStatus('Outbound');
  });

  it('tsp user can open a GBL from the shipment info page', function() {
    tspUserViewsGBL();
  });

  it('tsp user cannot generate gbl until office approves shipment', function() {
    tspUserCannotGenerateGBL();
  });
});

function tspUserCannotGenerateGBL() {
  const gblButtonText = 'Generate the GBL';

  cy.patientVisit('/queues/accepted');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/accepted/);
  });

  cy.selectQueueItemMoveLocator('GBLDIS');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  cy
    .get('button')
    .contains(gblButtonText)
    .should('be.disabled');
}

function tspUserGeneratesGBL() {
  const gblButtonText = 'Generate the GBL';

  // Open approved shipments queue
  cy.patientVisit('/queues/approved');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/approved/);
  });

  // Find shipment
  cy.selectQueueItemMoveLocator('GBLGBL');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  cy
    .get('button')
    .contains(gblButtonText)
    .should('be.enabled');

  cy.get('.documents').should('not.contain', 'Government Bill Of Lading');

  cy
    .get('button')
    .contains(gblButtonText)
    .click();

  cy
    .get('button')
    .contains(gblButtonText)
    .should('be.disabled');

  // I have seen this take anywhere from 8s - 18s. Until we optimize it, giving the test a long
  // timeout.
  cy.get('.usa-alert-success', { timeout: 20000 }).contains('GBL has been created');

  cy.get('.usa-alert-success').contains('Click the button to view, print, or download the GBL.');

  cy
    .get('button')
    .contains('View GBL')
    .should('be.enabled');

  cy
    .get('button')
    .contains(gblButtonText)
    .should('not.exist');

  cy.get('.documents').should('contain', 'Government Bill Of Lading');
}

function tspUserViewsGBL() {
  // Open approved shipments queue
  cy.patientVisit('/queues/approved');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/approved/);
  });

  // Find shipment
  cy.selectQueueItemMoveLocator('GBLGBL');

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
