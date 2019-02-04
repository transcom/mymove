import {
  fillAndSavePreApprovalRequest,
  editPreApprovalRequest,
  deletePreApprovalRequest,
} from '../../support/testCreatePreApprovalRequest';

/* global cy */
describe('TSP user interacts with pre approval request panel', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('TSP user creates pre approval request', function() {
    tspUserCreatesPreApprovalRequest();
  });
  it('TSP user edits pre approval request', function() {
    tspUserEditsPreApprovalRequest();
  });
  it('TSP user deletes pre approval request', function() {
    tspUserDeletesPreApprovalRequest();
  });
});

function tspUserCreatesPreApprovalRequest() {
  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find shipment and open it
  cy.selectQueueItemMoveLocator('DATESP');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  // TSP users should not able to approve pre approval requests
  cy.get('[data-test=approve-request]').should('not.exist');

  fillAndSavePreApprovalRequest();
  // Verify data has been saved in the UI
  cy.get('td').contains('Bulky Article: Motorcycle/Rec vehicle');
}

function tspUserEditsPreApprovalRequest() {
  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find shipment and open it
  cy.selectQueueItemMoveLocator('DATESP');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  editPreApprovalRequest();
  // Verify data has been saved in the UI
  cy.get('td').contains('notes notes edited');
}

function tspUserDeletesPreApprovalRequest() {
  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find shipment and open it
  cy.selectQueueItemMoveLocator('DATESP');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  deletePreApprovalRequest();
  cy
    .get('.pre-approval-panel td')
    .first()
    .should('not.contain', 'Bulky Article: Motorcycle/Rec vehicle');
}
