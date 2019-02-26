import {
  fillAndSavePreApprovalRequest,
  editPreApprovalRequest,
  deletePreApprovalRequest,
} from '../../support/preapprovals/testCreateRequest';
import { addOrigional105b, add105b } from '../../support/preapprovals/test105be';

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
  it('TSP user creates origional 105B request', function() {
    test105beOrigional();
  });
  it('TSP user creates 105B request', function() {
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

function test105beOrigional() {
  cy.setFeatureFlag('robustAccessorial=false');

  cy.selectQueueItemMoveLocator('DATESP');
  addOrigional105b();
  cy.get('td').contains('12.0000');
}

function test105be() {
  cy.selectQueueItemMoveLocator('DATESP');

  add105b();
  cy.get('td').contains(`12.0000 notes notes`);
  cy
    .get('td')
    .contains(`description description Crate: 25" x 25" x 25" (9.04 cu ft) Item: 30" x 30" x 30" notes notes`);
}
