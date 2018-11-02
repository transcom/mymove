import {
  getFixture,
  userClearsServiceAgent,
  userCancelsServiceAgent,
  userInputsServiceAgent,
  userSavesServiceAgent,
} from '../../support/testTspServiceAgents';

/* global cy */
describe('office user can view service agents', function() {
  beforeEach(() => {
    cy.signIntoOffice();
  });
  it('office user views and edits service agent panels', function() {
    officeUserEntersServiceAgents();
    userClearsServiceAgent('Origin');
    userInputsServiceAgent('OriginUpdate');
    userClearsServiceAgent('Destination');
    userInputsServiceAgent('DestinationUpdate');
    userSavesServiceAgent('OriginUpdate');
  });
});

function officeUserEntersServiceAgents() {
  // Open all moves queue
  cy.visit('/queues/all');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/all/);
  });

  // Find move and open it
  cy
    .get('div')
    .contains('LRKREK')
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

  // Verify that the Service Agent Panel contains expected data
  cy.get('span').contains('ACME Movers');

  // Click on edit Service Agent
  cy
    .get('.editable-panel-header')
    .contains('TSP & Servicing Agents')
    .siblings()
    .click();
}
