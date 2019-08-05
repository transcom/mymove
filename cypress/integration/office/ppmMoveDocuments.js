/* global cy */
describe('office user finds the shipment', function() {
  beforeEach(() => {
    cy.signIntoOffice();
    cy.get('[data-cy=ppm-queue]').click();
  });
  it('office user views ppm panel and goes to a ppm with a move document for awaiting review', function() {
    officeUserViewsPpmPanel('EXCLDE');
    officeUserChecksExpensePanelForAlert(false);
    officeUserEditsDocumentStatus('Expense Document', 'OK', 'EXCLUDE_FROM_CALCULATION');
    officeUserChecksExpensePanelForAlert(true);
  });
});

function officeUserViewsPpmPanel(locatorId) {
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/ppm/);
  });

  // Find shipment and open it
  cy.selectQueueItemMoveLocator(locatorId);

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  cy
    .get('.nav-tab')
    .contains('PPM')
    .click();
}

function officeUserEditsDocumentStatus(documentTitle, oldDocumentStatus, newDocumentStatus) {
  cy
    .get('.documents')
    .contains(documentTitle)
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+/);
  });

  cy.get('.panel-field.status').contains(oldDocumentStatus);

  cy.get('.editable-panel-edit').click();

  cy
    .get('label[for="moveDocument.status"]')
    .siblings()
    .first()
    .children()
    .select(newDocumentStatus)
    .blur();

  cy.get('.editable-panel-save').click();

  cy.go(-1);
}

function officeUserChecksExpensePanelForAlert(alertShown) {
  const exist = alertShown ? 'be.visible' : 'not.be.visible';
  cy.get('.awaiting-expenses-warning').should(`${exist}`);
}
