import { testPremoveSurvey } from '../../support/testPremoveSurvey';

/* global cy */
describe('TSP User Completes Premove Survey', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('tsp user enters premove survey', function() {
    tspUserEntersPremoveSurvey();
  });
});

function tspUserEntersPremoveSurvey() {
  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find shipment and open it
  cy
    .get('div')
    .contains('BACON1')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/shipments\/[^/]+/);
  });

  testPremoveSurvey();
}
