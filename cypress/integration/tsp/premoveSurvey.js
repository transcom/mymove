import { fillAndSavePremoveSurvey, testPremoveSurvey } from '../../support/testPremoveSurvey';

/* global cy */
describe('TSP User Completes Premove Survey', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('tsp user enters premove survey', function() {
    tspUserEntersPremoveSurveyUnprompted();
  });
  it('tsp user uses action button to enter premove survey', function() {
    tspUserClicksEnterPreMoveSurvey();
    tspUserFillsInPreMoveSurvey();
    tspUserVerifiesPreMoveSurveyEntered();
  });
});

function tspUserEntersPremoveSurveyUnprompted() {
  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find shipment and open it
  cy
    .get('div')
    .contains('PREMVE')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  testPremoveSurvey();
}

function tspUserClicksEnterPreMoveSurvey() {
  // Open approved shipments queue
  cy.visit('/queues/approved');

  // Find shipment and open it
  cy
    .get('div')
    .contains('ENTPMS')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  // Status should be Approved for "Enter pre-move survey" button to exist
  cy
    .get('li')
    .get('b')
    .contains('Approved');

  cy.wait(5000);

  cy
    .get('button')
    .contains('Enter pre-move survey')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Enter pre-move survey')
    .click();
}

function tspUserFillsInPreMoveSurvey() {
  fillAndSavePremoveSurvey();
}

function tspUserVerifiesPreMoveSurveyEntered() {
  cy.get('button').should('not.contain', 'Enter pre-move survey');
}
