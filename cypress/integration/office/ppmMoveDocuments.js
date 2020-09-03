describe('office user finds the shipment', function () {
  before(() => {
    cy.prepareOfficeApp();
  });

  beforeEach(() => {
    cy.signIntoOffice();
    cy.get('[data-testid=ppm-queue]').click();
  });

  it('office user views ppm panel and goes to a ppm with a move document for awaiting review', function () {
    officeUserViewsPpmPanel('PMTRVW');
    officeUserChecksExpensePanelForAlert(true);
  });
});

function officeUserViewsPpmPanel(locatorId) {
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/queues\/all/);
  });

  // Find shipment and open it
  cy.selectQueueItemMoveLocator(locatorId);

  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  cy.get('.nav-tab').contains('PPM').click();
}

function officeUserChecksExpensePanelForAlert(alertShown) {
  const exist = alertShown ? 'be.visible' : 'not.be.visible';
  cy.get('.awaiting-expenses-warning').should(`${exist}`);
}
