/* global cy */
describe('TSP User Checks Shipment Basics', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('tsp user sees service member customer info', function() {
    tspUserViewsCustomerInfo();
  });
});

function tspUserViewsCustomerInfo() {
  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find a shipment and open it
  cy.selectQueueItemMoveLocator('BACON4');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  // Expect Customer Info to be loaded
  cy
    .get('.customer-info')
    .contains('Customer Info')
    .get('.extras.content')
    .should($div => {
      const text = $div.text();
      expect(text).to.include('Preferred contact method:');
      expect(text).to.include('Backup Contacts');
    });
}
