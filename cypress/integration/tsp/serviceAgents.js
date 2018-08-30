/* global cy */
describe('service agents', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('tsp user enters and cancels origin service agent', function() {
    tspUserEntersServiceAgent('Origin');
    tspUserCancelsServiceAgent('Origin');
  });
  it('tsp user enters and cancels destination service agent', function() {
    tspUserEntersServiceAgent('Destination');
    tspUserCancelsServiceAgent('Destination');
  });
  it('tsp user enters origin service agent', function() {
    tspUserEntersServiceAgent('Origin');
    tspUserSavesServiceAgent('Origin');
  });
  it('tsp user enters destination service agent', function() {
    tspUserEntersServiceAgent('Destination');
    tspUserSavesServiceAgent('Destination');
  });
});

function getFixture(role) {
  return {
    Origin: {
      PointOfContact: 'Jenny at ACME Movers',
      Email: 'jenny_acme@example.com',
      Phone: '303-867-5309',
    },
    Destination: {
      PointOfContact: 'Alice at ACME Movers',
      Email: 'alice_acme@example.com',
      Phone: '303-867-5310',
    },
  }[role];
}

function tspUserEntersServiceAgent(role) {
  const fixture = getFixture(role);

  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find shipment and open it
  cy
    .get('div')
    .contains('KBACON')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/shipments\/[^/]+/);
  });

  // Click on edit Service Agent
  cy
    .get('.editable-panel-header')
    .contains(role)
    .siblings()
    .click();

  // Enter details in form and save orders
  cy
    .get('input[name="point_of_contact"]')
    .first()
    .type(fixture.PointOfContact)
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
    .get('div.point_of_contact')
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
    .get('div.point_of_contact')
    .get('span')
    .contains(fixture.PointOfContact);
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
    .get('div.point_of_contact')
    .get('span')
    .contains(fixture.PointOfContact);
  cy
    .get('div.email')
    .get('span')
    .contains(fixture.Email);
  cy
    .get('div.phone_number')
    .get('span')
    .contains(fixture.Phone);
}
