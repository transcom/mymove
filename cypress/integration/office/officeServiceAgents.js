import {
  userClearsServiceAgent,
  userInputsServiceAgent,
  userSavesServiceAgent,
  userCancelsServiceAgent,
  userVerifiesTspAssigned,
} from '../../support/testTspServiceAgents';

/* global cy */
describe('office user can view service agents', function() {
  beforeEach(() => {
    cy.signIntoOffice();
  });

  it('office user opens and cancels service agent panel', function() {
    officeUserOpensHhgPanelForMove('LRKREK');
    officeUserVerifiesServiceAgent();
    officeUserEditsServiceAgentPanel();
    userCancelsServiceAgent();
  });

  it('office user views and edits service agent panels', function() {
    officeUserOpensHhgPanelForMove('LRKREK');
    officeUserVerifiesServiceAgent();
    officeUserEditsServiceAgentPanel();
    officeUserSeesBlankTspData();
    userClearsServiceAgent('Origin');
    userInputsServiceAgent('OriginUpdate');
    userClearsServiceAgent('Destination');
    userInputsServiceAgent('DestinationUpdate');
    userSavesServiceAgent('OriginUpdate');
    officeUserSeesBlankTspData();
  });

  it('office user views tsp for awarded move', function() {
    officeUserOpensHhgPanelForMove('BACON1');
    userVerifiesTspAssigned();
    officeUserEditsServiceAgentPanel();
    userVerifiesTspAssigned();
  });
});

function officeUserOpensHhgPanelForMove(moveLocator) {
  // Open all moves queue
  cy.patientVisit('/queues/all');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/all/);
  });

  // Find move and open it
  cy.selectQueueItemMoveLocator(moveLocator);

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
}

function officeUserVerifiesServiceAgent() {
  // Verify that the Service Agent Panel contains expected data
  cy.get('span').contains('ACME Movers');
}

function officeUserEditsServiceAgentPanel() {
  // Click on edit Service Agent
  cy
    .get('.editable-panel-header')
    .contains('TSP & Servicing Agents')
    .siblings()
    .click();
}

function officeUserSeesBlankTspData() {
  cy
    .get('.editable-panel-3-column')
    .contains('TSP')
    .parent()
    .within(() => {
      cy
        .get('.panel-field')
        .contains('Name')
        .parent()
        .should('not.contain', 'undefined');
      cy
        .get('.panel-field')
        .contains('Email')
        .parent()
        .should('not.contain', 'undefined');
      cy
        .get('.panel-field')
        .contains('Phone number')
        .parent()
        .should('not.contain', 'undefined');
    });
}
