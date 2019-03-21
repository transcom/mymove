import {
  fillAndSavePreApprovalRequest,
  editPreApprovalRequest,
  deletePreApprovalRequest,
} from '../../support/preapprovals/testCreateRequest';
import { test35A, test35ALegacy } from '../../support/preapprovals/test35a';
import { test105be, test105beLegacy } from '../../support/preapprovals/test105be';

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
  it('Add legacy 35A to verify it displays correctly', function() {
    test35ALegacy();
  });
  it('TSP user creates 35A request', function() {
    test35A();
  });
  it('Add legacy 105B/E to verify they display correctly', function() {
    test105beLegacy();
  });
  it('TSP user creates 105B/E request', function() {
    test105be();
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
  cy.get('tr[data-cy="130B"]').should(td => {
    const text = td.text();
    expect(text).to.include('Bulky Article: Motorcycle/Rec vehicle');
  });
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
  cy.get('tr[data-cy="130B"]').should(td => {
    const text = td.text();
    expect(text).to.include('edited');
  });
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
