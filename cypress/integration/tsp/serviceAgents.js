/* global cy */
describe('TSP User enters and updates Service Agents', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('tsp user enters and cancels origin service agent', function() {
    tspUserEntersServiceAgent('Origin');
    tspUserSeesNoServiceAgent();
    tspUserInputsServiceAgent('Origin');
    tspUserCancelsServiceAgent('Origin');
  });
  it('tsp user enters and cancels destination service agent', function() {
    tspUserEntersServiceAgent('Destination');
    tspUserSeesNoServiceAgent();
    tspUserInputsServiceAgent('Destination');
    tspUserCancelsServiceAgent('Destination');
  });
  it('tsp user enters origin service agent', function() {
    tspUserEntersServiceAgent('Origin');
    tspUserSeesNoServiceAgent();
    tspUserInputsServiceAgent('Origin');
    tspUserSavesServiceAgent('Origin');
  });
  it('tsp user enters destination service agent', function() {
    tspUserEntersServiceAgent('Destination');
    tspUserSeesNoServiceAgent();
    tspUserInputsServiceAgent('Destination');
    tspUserSavesServiceAgent('Destination');
  });
  it('tsp user updates origin service agent', function() {
    tspUserEntersServiceAgent('Origin');
    tspUserClearsServiceAgent('Origin');
    tspUserInputsServiceAgent('OriginUpdate');
    tspUserSavesServiceAgent('OriginUpdate');
  });
  it('tsp user updates destination service agent', function() {
    tspUserEntersServiceAgent('Destination');
    tspUserClearsServiceAgent('Destination');
    tspUserInputsServiceAgent('DestinationUpdate');
    tspUserSavesServiceAgent('DestinationUpdate');
  });
  it('tsp user accepts a shipment', function() {
    tspUserAcceptsShipment();
  });

  it('tsp user assigns a service agent', function() {
    tspUserClicksAssignServiceAgent('ASSIGN');
    tspUserInputsServiceAgent('Origin');
    tspUserSavesServiceAgent('Origin');
    tspUserVerifiesServiceAgentAssigned();
  });
});

function getFixture(role) {
  return {
    Origin: {
      Company: 'ACME Movers',
      Email: 'acme@example.com',
      Phone: '303-867-5309',
    },
    OriginUpdate: {
      Company: 'ACME Movers',
      Email: 'acmemovers@example.com',
      Phone: '303-867-5308',
    },
    Destination: {
      Company: 'ACE Movers',
      Email: 'acmemoving@example.com',
      Phone: '303-867-5310',
    },
    DestinationUpdate: {
      Company: 'ACE Moving Company',
      Email: 'moveme@example.com',
      Phone: '303-867-5311',
    },
  }[role];
}

function tspUserSeesNoServiceAgent() {
  // Make sure the fields are empty to begin with
  // This helps make sure the test data hasn't changed elsewhere accidentally
  cy.get('input[name="company"]').should('have.value', '');
  cy.get('input[name="email"]').should('have.value', '');
  cy.get('input[name="phone_number"]').should('have.value', '');
}

function tspUserEntersServiceAgent(role) {
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
    .contains(role)
    .siblings()
    .click();
}

function tspUserInputsServiceAgent(role) {
  const fixture = getFixture(role);

  // Enter details in form
  cy
    .get('input[name="company"]')
    .first()
    .type(fixture.Company)
    .blur();
  cy
    .get('input[name="email"]')
    .first()
    .type(fixture.Email)
    .blur();
  cy
    .get('input[name="phone_number"]')
    .first()
    .type(fixture.Phone)
    .blur();
}

function tspUserClearsServiceAgent(role) {
  const fixture = getFixture(role);

  // Clear details in form
  cy
    .get('input[name="company"]')
    .clear()
    .blur();
  cy
    .get('input[name="email"]')
    .clear()
    .blur();
  cy
    .get('input[name="phone_number"]')
    .clear()
    .blur();
}

function tspUserCancelsServiceAgent(role) {
  cy
    .get('button')
    .contains('Cancel')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Cancel')
    .click();

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

function tspUserSavesServiceAgent(role) {
  const fixture = getFixture(role);

  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save')
    .click();

  // Verify data has been saved in the UI
  cy
    .get('div.company')
    .get('span')
    .contains(fixture.Company);
  cy
    .get('div.email')
    .get('span')
    .contains(fixture.Email);
  cy
    .get('div.phone_number')
    .get('span')
    .contains(fixture.Phone);

  // Refresh browser and make sure changes persist
  cy.reload();

  cy
    .get('div.company')
    .get('span')
    .contains(fixture.Company);
  cy
    .get('div.email')
    .get('span')
    .contains(fixture.Email);
  cy
    .get('div.phone_number')
    .get('span')
    .contains(fixture.Phone);
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
  cy
    .get('li')
    .get('b')
    .contains('Awarded');

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
  cy
    .get('li')
    .get('b')
    .contains('Accepted');

  cy.get('a').contains('All Shipments Queue');
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
  cy
    .get('li')
    .get('b')
    .contains('Accepted');

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
