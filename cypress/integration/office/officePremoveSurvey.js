import { testPremoveSurvey } from '../../support/testPremoveSurvey';

/* global cy */
describe('office user interacts with premove survey', function() {
  beforeEach(() => {
    cy.signIntoOffice();
  });
  it('office user enters premove survey', function() {
    officeUserEntersPreMoveSurvey();
  });
});

function officeUserEntersPreMoveSurvey() {
  // Open new moves queue
  cy.patientVisit('/queues/all');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/all/);
  });

  // Find move and open it
  cy
    .get('div')
    .contains('RLKBEM')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  // Click on HHG tab
  cy
    .get('span')
    .contains('HHG')
    .click();
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/hhg/);
  });

  // Verify that the Estimates section contains expected data
  cy.get('span').contains('2,000');

  testPremoveSurvey();
}
