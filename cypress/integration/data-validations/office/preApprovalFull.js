import {
  fillAndSavePreApprovalRequest,
  fillInvalidPreApprovalRequest,
  editPreApprovalRequest,
  approvePreApprovalRequest,
  deletePreApprovalRequest,
} from '../../../support/testCreatePreApprovalRequest';

/* global cy */
describe('office user interacts with pre approval request panel', function() {
  beforeEach(() => {
    cy.signIntoOffice();
  });

  // THIS TEST IS MEANT TO BE RUN AS ONE OFF AND ONLY IN DEV MACHINE
  // PLEASE DO NOT DELETE OR REMOVE THE SKIP
  // THE GOAL OF THIS TEST IS DATA VALIDATIONS FOR PROD DATA
  // TO RUN ONLY THIS TEST, REMOVE `SKIP` AND REPLACE WITH `ONLY`
  it.skip('office user iterates through every pre approval request', () => {
    officeUserIterateThroughAllPARS();
  });
});

function officeUserIterateThroughAllPARS() {
  cy.server();
  cy.route('POST', '/api/v1/shipments/*/accessorials').as('accessorialsCheck');
  cy.route('POST', '/internal/shipments/*/invoice').as('invoiceSubmit');

  // Open new moves queue
  cy.visit('/queues/all');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/all/);
  });

  // Find move and open it
  cy
    .get('div')
    .contains('DOOB')
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

  let PARS = [
    '4A',
    '4B',
    '28A',
    '28B',
    '28C',
    '35A',
    '105B',
    '105D',
    '105E',
    '120A',
    '120B',
    '120C',
    '120D',
    '120E',
    '120F',
    '125A',
    '125B',
    '125C',
    '125D',
    '130A',
    '130B',
    '130C',
    '130D',
    '130E',
    '130F',
    '130G',
    '130H',
    '130I',
    '130J',
    '175A',
    '225A',
    '225B',
    '226A',
  ];
  PARS.forEach(par => {
    //add a pre approval request first
    fillAndSavePreApprovalRequest(par);

    // Verify par is displayed in Pre-Approval Request panel and approve the request
    cy.wait('@accessorialsCheck');
    cy.get('[data-cy=' + par + ']').should('contain', par);
    cy.get('[data-test=approve-request]').click({ multiple: true });
  });

  //invoice the shipment and ensure success
  cy
    .get('button')
    .contains('Approve Payment')
    .click()
    .then(() => {
      cy
        .get('button')
        .contains('Approve')
        .click()
        .wait('@invoiceSubmit');
      cy.get('[data-cy="invoice success message"]');
    });
}
