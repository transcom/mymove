import {
  fillAndSavePreApprovalRequest,
  fillInvalidPreApprovalRequest,
  editPreApprovalRequest,
  approvePreApprovalRequest,
  deletePreApprovalRequest,
} from '../../support/testCreatePreApprovalRequest';

/* global cy */
describe('office user interacts with pre approval request panel', function() {
  beforeEach(() => {
    cy.signIntoOffice();
  });
  it('office user creates pre approval request', function() {
    officeUserCannotAddInvalidPreApprovalRequest();
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
});

function officeUserCreatesPreApprovalRequest() {
  cy.server();
  cy.route('GET', '/internal/queues/all').as('queuesAll');
  cy.route('POST', '/api/v1/shipments/*/accessorials').as('accessorialsCheck');
  // Open new moves queue
  cy.patientVisit('/queues/all');
  cy.wait('@queuesAll');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/all/);
  });

  // Find move and open it
  cy
    .get('div')
    .contains('RLKBEM')
    .dblclick();

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

  fillAndSavePreApprovalRequest('130B');
  cy.wait('@accessorialsCheck');
  // Verify data has been saved in the UI
  cy.get('td').contains('Bulky Article: Motorcycle/Rec vehicle');
}

function officeUserCannotAddInvalidPreApprovalRequest() {
  cy.server();
  cy.route('GET', '/internal/queues/all').as('queuesAll');
  // Open new moves queue
  cy.patientVisit('/queues/all');
  cy.wait('@queuesAll');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/all/);
  });

  // Find move and open it
  cy
    .get('div')
    .contains('RLKBEM')
    .dblclick();

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

  fillInvalidPreApprovalRequest();
}
function officeUserEditsPreApprovalRequest() {
  cy.server();
  cy.route('GET', '/internal/queues/all').as('queuesAll');
  cy.route('PUT', '/api/v1/shipments/accessorials/*').as('updateClick');
  // Open new moves queue
  cy.patientVisit('/queues/all');
  cy.wait('@queuesAll');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/all/);
  });

  // Find move and open it
  cy
    .get('div')
    .contains('RLKBEM')
    .dblclick();

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
  cy.wait('@updateClick');
  cy.get('td').contains('notes notes edited');
}

function officeUserApprovesPreApprovalRequest() {
  cy.server();
  cy.route('GET', '/internal/queues/all').as('queuesAll');
  cy.route('POST', '/api/v1/shipments/accessorials/*/approve').as('approveClick');
  // cy.route('PUT', '/api/v1/shipments/accessorials/*').as('updateClick');
  // Open new moves queue
  cy.patientVisit('/queues/all');
  cy.wait('@queuesAll');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/all/);
  });

  // Find move and open it
  cy
    .get('div')
    .contains('RLKBEM')
    .dblclick();

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
  cy.wait('@approveClick');
  cy.get('.pre-approval-panel td').contains('Approved');
}

function officeUserDeletesPreApprovalRequest() {
  cy.server();
  cy.route('GET', '/internal/queues/all').as('queuesAll');
  cy.route('DELETE', '/api/v1/shipments/accessorials/*').as('deleteClick');
  // Open new moves queue
  cy.patientVisit('/queues/all');
  cy.wait('@queuesAll');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/all/);
  });

  // Find move and open it
  cy
    .get('div')
    .contains('RLKBEM')
    .dblclick();

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
  cy.wait('@deleteClick');
  cy
    .get('.pre-approval-panel td')
    .first()
    .should('not.contain', 'Bulky Article: Motorcycle/Rec vehicle');
}
