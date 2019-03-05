import {
  fillAndSavePreApprovalRequest,
  editPreApprovalRequest,
  deletePreApprovalRequest,
} from '../../support/preapprovals/testCreateRequest';
import { addOriginal105, add105 } from '../../support/preapprovals/test105be';

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
    test105beOriginal();
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

function test105beOriginal() {
  cy.setFeatureFlag('robustAccessorial=false');
  cy.selectQueueItemMoveLocator('DATESP');

  addOriginal105({ code: '105B', quantity1: 12 });
  cy.get('td[details-cy="105B-default-details"]').should('contain', '12.0000 notes notes 105B');

  addOriginal105({ code: '105E', quantity1: 90 });
  cy.get('td[details-cy="105E-default-details"]').should('contain', '90.0000 notes notes 105E');
}

function test105be() {
  cy.selectQueueItemMoveLocator('DATESP');

  add105({ code: '105B', itemSize: 30, crateSize: 25 });
  cy
    .get('td[details-cy="105B-details"]')
    .should(
      'contain',
      'description description 105B Crate: 25" x 25" x 25" (9.04 cu ft) Item: 30" x 30" x 30" notes notes',
    );

  add105({ code: '105E', itemSize: 40, crateSize: 50 });
  cy
    .get('td[details-cy="105E-details"]')
    .should(
      'contain',
      'description description 105E Crate: 50" x 50" x 50" (72.33 cu ft) Item: 40" x 40" x 40" notes notes',
    );
}
