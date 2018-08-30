/* global cy */
describe('service agents', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('tsp user enters origin service agent', function() {
    tspUserEntersServiceAgent('Origin');
  });
  it('tsp user enters destination service agent', function() {
    tspUserEntersServiceAgent('Destination');
  });
});

function tspUserEntersServiceAgent(role) {
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
    .type('Jenny at ACME Movers')
    .blur();
  cy
    .get('input[name="email"]')
    .first()
    .type('jenny_acme@example.com')
    .blur();
  cy
    .get('input[name="phone_number"]')
    .first()
    .type('303-867-5309')
    .blur();

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
    .contains('Jenny at ACME Movers');
  cy
    .get('div.email')
    .get('span')
    .contains('jenny_acme@example.com');
  cy
    .get('div.phone_number')
    .get('span')
    .contains('303-867-5309');

  // Refresh browser and make sure changes persist
  cy.reload();

  cy
    .get('div.point_of_contact')
    .get('span')
    .contains('Jenny at ACME Movers');
  cy
    .get('div.email')
    .get('span')
    .contains('jenny_acme@example.com');
  cy
    .get('div.phone_number')
    .get('span')
    .contains('303-867-5309');
}
