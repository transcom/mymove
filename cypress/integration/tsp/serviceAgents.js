import {
  getFixture,
  userClearsServiceAgent,
  userCancelsServiceAgent,
  userInputsServiceAgent,
  userSavesServiceAgent,
} from '../../support/testTspServiceAgents';
import { tspUserVerifiesShipmentStatus } from '../../support/testTspStatus';

/* global cy */
describe('TSP User enters and updates Service Agents', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('tsp user enters and cancels origin service agent', function() {
    tspUserEntersServiceAgent();
    tspUserSeesNoServiceAgent('Origin');
    userInputsServiceAgent('Origin');
    userCancelsServiceAgent('Origin');
    tspUserViewsBlankServiceAgent();
  });
  it('tsp user enters and cancels destination service agent', function() {
    tspUserEntersServiceAgent();
    tspUserSeesNoServiceAgent('Destination');
    userInputsServiceAgent('Destination');
    userCancelsServiceAgent('Destination');
    tspUserViewsBlankServiceAgent();
  });
  it('tsp user enters origin and destination service agents', function() {
    tspUserEntersServiceAgent();
    tspUserSeesNoServiceAgent('Origin');
    userInputsServiceAgent('Origin');
    tspUserSeesNoServiceAgent('Destination');
    userInputsServiceAgent('Destination');
    userSavesServiceAgent('Destination');
  });
  it('tsp user updates origin and destination service agents', function() {
    tspUserEntersServiceAgent();
    userClearsServiceAgent('Origin');
    userInputsServiceAgent('OriginUpdate');
    userClearsServiceAgent('Destination');
    userInputsServiceAgent('DestinationUpdate');
    userSavesServiceAgent('OriginUpdate');
  });
  it('tsp user accepts a shipment', function() {
    tspUserAcceptsShipment();
  });

  it('tsp user assigns origin and destination service agents using action button', function() {
    tspUserClicksAssignServiceAgent('ASSIGN');
    userInputsServiceAgent('Origin');
    userInputsServiceAgent('Destination');
    userSavesServiceAgentsWizard();
    tspUserVerifiesServiceAgentAssigned();
  });
});

function tspUserSeesNoServiceAgent(role) {
  const fixture = getFixture(role);
  // Make sure the fields are empty to begin with
  // This helps make sure the test data hasn't changed elsewhere accidentally
  cy.get('input[name="' + fixture.Role + '.company"]').should('have.value', '');
  cy.get('input[name="' + fixture.Role + '.email"]').should('have.value', '');
  cy.get('input[name="' + fixture.Role + '.phone_number"]').should('have.value', '');
}

function tspUserViewsBlankServiceAgent() {
  // Verify data has been saved in the UI
  cy
    .get('div.company')
    .get('span')
    .contains('missing');
  cy
    .get('div.email')
    .get('span')
    .contains('missing');
  cy
    .get('div.phone_number')
    .get('span')
    .contains('missing');
}

function tspUserEntersServiceAgent() {
  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find shipment and open it
  cy
    .get('div')
    .contains('BACON2')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  // Click on edit Service Agent
  cy
    .get('.editable-panel-header')
    .contains('TSP & Servicing Agents')
    .siblings()
    .click();
}

function tspUserAcceptsShipment() {
  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find shipment and open it
  cy
    .get('div')
    .contains('BACON2')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  // Status should be Awarded
  tspUserVerifiesShipmentStatus('Shipment awarded');

  cy.get('a').contains('New Shipments Queue');

  cy
    .get('button')
    .contains('Accept Shipment')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Accept Shipment')
    .click();

  // Status should be Accepted
  tspUserVerifiesShipmentStatus('Shipment accepted');

  cy.get('a').contains('Accepted Shipments Queue');
}

function tspUserClicksAssignServiceAgent(locator) {
  cy.visit('/queues/all');

  // Find shipment and open it
  cy
    .get('div')
    .contains(locator)
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  // Status should be Accepted or Approved for "Assign servicing agents" button to exist
  tspUserVerifiesShipmentStatus('Shipment accepted');

  cy
    .get('button')
    .contains('Assign servicing agents')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Assign servicing agents')
    .click();
}

function tspUserVerifiesServiceAgentAssigned() {
  cy.get('button').should('not.contain', 'Assign servicing agents');
}

function userSavesServiceAgentsWizard() {
  const origin = getFixture('Origin');
  const destination = getFixture('Destination');

  cy
    .get('button')
    .contains('Done')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Done')
    .click();

  // Verify data has been saved in the UI
  cy
    .get('div.company')
    .get('span')
    .contains(origin.Company);
  cy
    .get('div.email')
    .get('span')
    .contains(origin.Email);
  cy
    .get('div.phone_number')
    .get('span')
    .contains(origin.Phone);

  // Refresh browser and make sure changes persist
  cy.reload();

  cy
    .get('div.company')
    .get('span')
    .contains(origin.Company);
  cy
    .get('div.email')
    .get('span')
    .contains(origin.Email);
  cy
    .get('div.phone_number')
    .get('span')
    .contains(origin.Phone);

  cy
    .get('div.company')
    .get('span')
    .contains(destination.Company);
  cy
    .get('div.email')
    .get('span')
    .contains(destination.Email);
  cy
    .get('div.phone_number')
    .get('span')
    .contains(destination.Phone);
}
