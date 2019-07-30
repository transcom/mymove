/* global cy */
describe('office user finds the shipment', function() {
  beforeEach(() => {
    cy.signIntoOffice();
  });
  it('office user views ppm panel and sees and alert for missing expense receipts', function() {
    officeUserViewsPpmPanel('EXCDEX');
  });
});

function officeUserViewsPpmPanel(locatorId) {
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find shipment and open it
  cy.selectQueueItemMoveLocator(locatorId);

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  cy
    .get('.nav-tab')
    .contains('PPM')
    .click();
}
