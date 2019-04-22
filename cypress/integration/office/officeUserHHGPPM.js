/* global cy */

import { fileUploadTimeout } from '../../support/constants';

describe('office user finds the shipment', function() {
  beforeEach(() => {
    cy.signIntoOffice();
  });
  const queues = {
    new: 'new',
    ppm: 'ppm',
    acceptedHHG: 'hhg_accepted',
    deliveredHHG: 'hhg_delivered',
    completedHHG: 'hhg_completed',
    all: 'all',
  };
  it('office user sees accepted HHG in Accepted HHGs queue', function() {
    officeUserVisitsQueue(queues.acceptedHHG);
    officeUserViewsMove('COMBO2');
  });
  it('office user sees delivered HHG in Delivered HHG queue', function() {
    officeUserVisitsQueue(queues.deliveredHHG);
    officeUserViewsMove('COMBO3');
  });
  it('office user sees completed HHG in Completed HHGs queue', function() {
    officeUserVisitsQueue(queues.completedHHG);
    officeUserViewsMove('COMBO4');
  });
  it('office user approves shipment', function() {
    officeUserVisitsQueue(queues.all);
    officeUserViewsMove('COMBO2');
    officeUserVisitsPPMTab();
    officeUserVisitsHHGTab();
    officeUserApprovesShipment();
  });
  it.only('office user uploads document', function() {
    officeUserSubmitsDocument();
  });
});

function officeUserApprovesShipment() {
  cy.get('.combo-button').click();

  // Approve HHG
  cy
    .get('.combo-button .dropdown')
    .contains('Approve HHG')
    .click();

  cy
    .get('.combo-button .dropdown')
    .contains('Approve HHG')
    .should('have.class', 'disabled');

  cy.get('.status').contains('Approved');
}

function officeUserVisitsQueue(queue) {
  const queueName = queue.toLowerCase();
  // eslint-disable-next-line security/detect-non-literal-regexp
  const routePattern = new RegExp(`^/queues/${queueName}`);
  cy.patientVisit(`/queues/${queueName}`);

  cy.location().should(loc => {
    expect(loc.pathname).to.match(routePattern);
  });
}

function officeUserVisitsPPMTab() {
  // navtab
  cy
    .get('a')
    .contains('PPM')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/ppm/);
  });
}

function officeUserVisitsHHGTab() {
  // navtab
  cy
    .get('a')
    .contains('HHG')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/hhg/);
  });
}

function officeUserViewsMove(locator) {
  cy.selectQueueItemMoveLocator(locator);

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });
}

function officeUserSubmitsDocument() {
  cy.patientVisit('moves/27266e89-df79-4469-8843-05b45741a818/documents/new');

  cy.contains('Upload a new document');
  cy.get('button.submit').should('be.disabled');
  cy.get('select[name="move_document_type"]').select('Expense');
  cy.get('input[name="title"]').type('expense document for combo');
  cy.get('select[name="moving_expense_type"]').select('Contracted Expense');
  cy.get('input[name="requested_amount_cents"]').type('4,000.92');
  cy.get('select[name="payment_method"]').select('Other account');

  cy.get('button.submit').should('be.disabled');

  cy.upload_file('.filepond--root', 'top-secret.png');
  cy
    .get('button.submit', { timeout: fileUploadTimeout })
    .should('not.be.disabled')
    .click();

  cy.contains('expense document for combo');
}
