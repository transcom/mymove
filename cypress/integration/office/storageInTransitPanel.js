/* global cy */
describe('office user finds the shipment', function() {
  beforeEach(() => {
    cy.signIntoOffice();
  });
  it('office user views storage in transit panel', function() {
    officeUserViewsSITPanel();
  });
  it('office user starts and cancels sit approval', function() {
    officeUserStartsAndCancelsSitApproval();
  });
});

function officeUserViewsSITPanel() {
  // Open new moves queue
  cy.patientVisit('/queues/hhg_accepted');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/hhg_accepted/);
  });

  // Find move (generated in e2ebasic.go) and open it
  cy.selectQueueItemMoveLocator('SITREQ');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  cy
    .get('a')
    .contains('HHG')
    .click(); // navtab

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/hhg/);
  });
  cy.get('.storage-in-transit-panel').contains('Storage in Transit');
  cy.get('.storage-in-transit').within(() => {
    cy.contains('Destination SIT');
    cy
      .get('.sit-status-text')
      .contains('Status')
      .parent()
      .siblings()
      .contains('SIT Requested');
    cy
      .contains('Dates')
      .siblings()
      .contains('Est. start date')
      .siblings()
      .contains('22-Mar-2019');
    cy
      .contains('Note')
      .siblings()
      .contains('Shipper phoned to let us know he is delayed until next week.');
    cy
      .contains('Warehouse')
      .siblings()
      .contains('Warehouse ID')
      .siblings()
      .contains('000383');
    cy
      .contains('Warehouse')
      .siblings()
      .contains('Contact info')
      .siblings()
      .contains('(713) 868-3497');
  });
}

function officeUserStartsAndCancelsSitApproval() {
  cy.patientVisit('/queues/hhg_accepted');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/hhg_accepted/);
  });

  cy.selectQueueItemMoveLocator('SITREQ');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  cy
    .get('a')
    .contains('HHG')
    .click(); // navtab

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/hhg/);
  });

  cy
    .get('a')
    .contains('Approve')
    .click()
    .get('.storage-in-transit')
    .should($div => {
      const text = $div.text();
      expect(text).to.include('Approve SIT Request');
      expect(text).to.not.include('Deny');
      expect(text).to.not.include('Edit');
    });

  cy.get('input[name="authorized_start_date"]').should('have.value', '3/22/2019');

  cy
    .get('.usa-button-secondary')
    .contains('Cancel')
    .click()
    .get('.storage-in-transit-panel .add-request')
    .should($div => {
      const text = $div.text();
      expect(text).to.not.include('Approve SIT Request');
    });
}
