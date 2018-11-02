/* global cy */
describe('TSP User Checks Customer Info Panel', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('tsp user sees customer info', function() {
    tspUserViewsCustomerInfo();
  });
});

function tspUserOpensAShipment() {
  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find a shipment and open it
  cy
    .get('div')
    .contains('BACON1')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });
}

function tspUserViewsCustomerContactInfo() {
  // Check the name of the service member is bolded
  cy.get('.customer-info').contains('b', 'Submitted, HHG');

  // Check for DoD ID
  cy.get('.customer-info').contains('4444567890');

  // Check for DoD Branch and rank
  cy.get('.customer-info').contains('ARMY - E_1');

  // Check for phone number
  cy.get('.customer-info').contains('555-555-5555');

  // Check for email
  cy
    .get('a')
    .contains('hhg@award.ed')
    .and('have.attr', 'href')
    .and('match', /mailto:hhg@award.ed/);
}

function tspUserViewsBackupContactInfo() {
  // Check Backup Contacts is bolded
  cy.get('.customer-info').contains('b', 'Backup Contacts');

  // Check the name
  cy.get('.customer-info').contains('name (EDIT)');

  // Check for phone number
  cy.get('.customer-info').contains('555-555-5555');

  // Check for email
  cy
    .get('a')
    .contains('email@example.com')
    .and('have.attr', 'href')
    .and('match', /mailto:email@example.com/);
}

function tspUserViewsEntitlementInfo() {
  // Check Entitlements is bolded
  cy.get('.customer-info').contains('b', 'Entitlements');

  // Check the hhg entitlements
  cy.get('.customer-info').contains('8,000 lbs');

  // Check for pro-gear and spouse pro-gear
  cy.get('.customer-info').contains('Pro-gear: 2,000 lbs / Spouse: 500 ');
}

function tspUserViewsCustomerInfo() {
  tspUserOpensAShipment();
  tspUserViewsCustomerContactInfo();
  tspUserViewsBackupContactInfo();
  tspUserViewsEntitlementInfo();
}
