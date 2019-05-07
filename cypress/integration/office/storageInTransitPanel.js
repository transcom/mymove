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
  it('office approves sit', function() {
    officeUserApprovesSITRequest();
  });
  it('office user starts and cancels sit edit', function() {
    officeUserStartsAndCancelsSitEdit();
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

  cy.get('[data-cy="hhg-tab"]').click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/hhg/);
  });
  cy.get('[data-cy=storage-in-transit-panel]').contains('Storage in Transit');
  cy.get('[data-cy=storage-in-transit]').within(() => {
    cy.contains('Destination SIT');
    cy
      .get('[data-cy=sit-status-text]')
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

  cy.get('[data-cy="hhg-tab"]').click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/hhg/);
  });

  cy
    .get('a')
    .contains('Approve')
    .click()
    .get('[data-cy=storage-in-transit]')
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
    .get('[data-cy=storage-in-transit-panel] [data-cy=add-request]')
    .should($div => {
      const text = $div.text();
      expect(text).to.not.include('Approve SIT Request');
    });
}

function officeUserStartsAndCancelsSitEdit() {
  cy.patientVisit('/queues/new');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  cy.selectQueueItemMoveLocator('SITAPR');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  cy.get('[data-cy="hhg-tab"]').click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/hhg/);
  });

  cy
    .get('[data-cy=storage-in-transit]')
    .contains('Edit')
    .click()
    .get('.sit-authorization')
    .should($div => {
      const text = $div.text();
      expect(text).to.include('Edit SIT authorization');
    });

  cy
    .get('.usa-button-secondary')
    .contains('Cancel')
    .click()
    .get('[data-cy=storage-in-transit]')
    .should($div => {
      const text = $div.text();
      expect(text).to.include('Approved');
      expect(text).to.not.include('Edit SIT authorization');
    });
}

function officeUserApprovesSITRequest() {
  cy.patientVisit('/queues/hhg_accepted');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/hhg_accepted/);
  });

  cy.selectQueueItemMoveLocator('SITREQ');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  cy.get('[data-cy="hhg-tab"]').click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/hhg/);
  });

  cy
    .get('a')
    .contains('Approve')
    .click()
    .get('[data-cy="storage-in-transit"]')
    .should($div => {
      const text = $div.text();
      expect(text).to.include('Approve SIT Request');
      expect(text).to.not.include('Deny');
      expect(text).to.not.include('Edit');
    });

  cy.get('input[name="authorized_start_date"]').should('have.value', '3/22/2019');

  cy.get('textarea[name="authorization_notes"]').type('this is a note', { force: true, delay: 150 });

  cy
    .get('[data-cy="storage-in-transit-approve-button"]')
    .contains('Approve')
    .click();

  // Refresh browser and make sure changes persist
  cy.patientReload();

  cy.get('[data-cy="storage-in-transit-status"]').contains('Approved');
  cy.get('[data-cy="sit-authorized-start-date"]').contains('22-Mar-2019');
  cy.get('[data-cy="sit-authorization-notes"]').contains('this is a note');
}
