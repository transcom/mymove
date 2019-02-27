import { userSavesStorageDetails, userCancelsStorageDetails } from '../../support/storagePanel';

/* global cy */
describe('office user interacts with storage panel', function() {
  beforeEach(() => {
    cy.signIntoOffice();
  });

  it('office user completes storage panel', function() {
    officeUserGoesToStoragePanel('VGHEIS');
    userSavesStorageDetails();
  });

  it('office user cancels storage panel', function() {
    officeUserGoesToStoragePanel('NOADVC');
    userCancelsStorageDetails();
  });
});

function officeUserGoesToStoragePanel(locator) {
  // Open new moves queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find shipment and open it
  cy.selectQueueItemMoveLocator(locator);

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  cy
    .get('.title')
    .contains('PPM')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/ppm/);
  });

  cy
    .get('.editable-panel-header')
    .contains('Storage')
    .siblings()
    .click();
}
