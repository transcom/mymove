import {
  fillAndSavePreApprovalRequest,
  editPreApprovalRequest,
  approvePreApprovalRequest,
  deletePreApprovalRequest,
} from '../../support/preapprovals/testCreateRequest';
import { add35A } from '../../support/preapprovals/test35a';

/* global cy */
describe('office user interacts with pre approval request panel', function() {
  beforeEach(() => {
    cy.signIntoOffice();
  });
  it('office user creates pre approval request', function() {
    officeUserCreatesPreApprovalRequest();
  });
  it('office user edits pre approval request', function() {
    officeUserEditsPreApprovalRequest();
  });
  it('office user approves pre approval request', function() {
    officeUserApprovesPreApprovalRequest();
  });
  it('office user deletes pre approval request', function() {
    officeUserDeletesPreApprovalRequest();
  });
  it('office user approves and later edits 35A pre approval request', officeUserCreates35APreApprovalRequest);
});

function officeUserCreatesPreApprovalRequest() {
  // Open new moves queue
  cy.patientVisit('/queues/all');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/all/);
  });

  // Find move and open it
  cy.selectQueueItemMoveLocator('RLKBEM');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  // Click on HHG tab
  cy
    .get('span')
    .contains('HHG')
    .click();
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/hhg/);
  });

  // Verify that the Estimates section contains expected data
  cy.get('span').contains('2,000');

  fillAndSavePreApprovalRequest();
  // Verify data has been saved in the UI
  cy.get('tr[data-cy="130B"]').should(td => {
    const text = td.text();
    expect(text).to.include('Bulky Article: Motorcycle/Rec vehicle');
  });
}
function officeUserEditsPreApprovalRequest() {
  // Open new moves queue
  cy.patientVisit('/queues/all');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/all/);
  });

  // Find move and open it
  cy.selectQueueItemMoveLocator('RLKBEM');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  // Click on HHG tab
  cy
    .get('span')
    .contains('HHG')
    .click();
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/hhg/);
  });

  // Verify that the Estimates section contains expected data
  cy.get('span').contains('2,000');

  editPreApprovalRequest();
  // Verify data has been saved in the UI
  cy.get('tr[data-cy="130B"]').should(td => {
    const text = td.text();
    expect(text).to.include('edited');
  });
}

function officeUserApprovesPreApprovalRequest() {
  // Open new moves queue
  cy.patientVisit('/queues/all');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/all/);
  });

  // Find move and open it
  cy.selectQueueItemMoveLocator('RLKBEM');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  // Click on HHG tab
  cy
    .get('span')
    .contains('HHG')
    .click();
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/hhg/);
  });

  // Verify that the Estimates section contains expected data
  cy.get('span').contains('2,000');

  approvePreApprovalRequest();
  cy.get('.pre-approval-panel td').contains('Approved');
}

function officeUserDeletesPreApprovalRequest() {
  // Open new moves queue
  cy.patientVisit('/queues/all');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/all/);
  });

  // Find move and open it
  cy.selectQueueItemMoveLocator('RLKBEM');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  // Click on HHG tab
  cy
    .get('span')
    .contains('HHG')
    .click();
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/hhg/);
  });

  // Verify that the Estimates section contains expected data
  cy.get('span').contains('2,000');

  deletePreApprovalRequest();
  cy
    .get('.pre-approval-panel td')
    .first()
    .should('not.contain', 'Bulky Article: Motorcycle/Rec vehicle');
}

function officeUserCreates35APreApprovalRequest() {
  // Open new moves queue
  cy.patientVisit('/queues/all');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/all/);
  });

  // Find move and open it
  cy.selectQueueItemMoveLocator('DOOB');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  // Click on HHG tab
  cy
    .get('span')
    .contains('HHG')
    .click();
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/hhg/);
  });

  // Verify that the Estimates section contains expected data
  cy.get('span').contains('2,000');

  // Enter 35A info (without Actual Amount) and save
  add35A({ estimate_amount_dollars: 250 });

  // Approve the request
  approvePreApprovalRequest();

  // Edit the request
  cy
    .get('[data-test=edit-request]')
    .first()
    .click();

  cy.typeInInput({ name: 'actual_amount_cents', value: `235` });

  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save')
    .click();

  cy
    .get('td[details-cy="35A-details"]')
    .should(
      'contain',
      'description description 35A reason reason 35A Est. not to exceed: $250.00 Actual amount: $235.00',
    );

  // Edit the request again
  cy
    .get('[data-test=edit-request]')
    .first()
    .click();

  cy.typeInInput({ name: 'actual_amount_cents', value: `220` });

  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save')
    .click();

  cy
    .get('td[details-cy="35A-details"]')
    .should(
      'contain',
      'description description 35A reason reason 35A Est. not to exceed: $250.00 Actual amount: $220.00',
    );

  // The edit should propagate to the Invoice panel
  cy
    .get('[data-cy=unbilled-table] tbody')
    .contains('$220.00')
    .should('exist');

  // Unset the actual amount
  cy
    .get('[data-test=edit-request]')
    .first()
    .click();

  cy.clearInput({ name: 'actual_amount_cents' });

  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save')
    .click();

  cy
    .get('td[details-cy="35A-details"]')
    .should('contain', 'description description 35A reason reason 35A Est. not to exceed: $250.00 Actual amount: --');

  // The edit should propagate to the Invoice panel
  cy
    .get('[data-cy=unbilled-table] tbody')
    .contains('Missing actual amount')
    .parent();
}
