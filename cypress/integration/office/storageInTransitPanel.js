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
  it('office approves sit and edits the approval', function() {
    officeUserApprovesSITRequest();
  });
  it('office user starts and cancels sit edit', function() {
    officeUserStartsAndCancelsSitEdit();
  });
  it('office user denies and edits denied sit request', function() {
    officeUserDeniesSITRequest();
  });
  it('office user views remaining days and status of shipment in SIT (with frozen clock)', function() {
    officeUserEntitlementRemainingDays();
  });
  it('office user views remaining days and status of shipment expired in SIT (with frozen clock)', function() {
    officeUserEntitlementRemainingDaysExpired();
  });
});

function officeUserViewsSITPanel() {
  // Open hhg_accepted moves queue
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
  // Open hhg_accepted moves queue
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
  // Open new moves queue
  cy.patientVisit('/queues/new');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find move (generated in e2ebasic.go) and open it
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
  // Open hhg_accepted moves queue
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

  cy.get('[data-cy="sit-edit-link"]').click();

  cy.get('input[name="authorized_start_date"]').clear();
  cy.get('input[name="authorized_start_date"]').type('3/23/2019', { force: true, delay: 150 });
  cy.get('input[name="authorized_start_date"]').should('have.value', '3/23/2019');

  cy.get('textarea[name="authorization_notes"]').type('this is also a note', { force: true, delay: 150 });

  cy
    .get('[data-cy="sit-editor-save-button"]')
    .contains('Save')
    .click();

  cy.patientReload();

  cy.get('[data-cy="sit-authorization-notes"]').contains('this is also a note');
  cy.get('[data-cy="sit-authorized-start-date"]').contains('23-Mar-2019');
}

function officeUserDeniesSITRequest() {
  // Open hhg_accepted moves queue
  cy.patientVisit('/queues/hhg_accepted');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/hhg_accepted/);
  });

  // Find move (generated in e2ebasic.go) and open it
  cy.selectQueueItemMoveLocator('SITDEN');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  cy.get('[data-cy="hhg-tab"]').click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/hhg/);
  });

  cy
    .get('a[data-cy="deny-sit-link"]')
    .contains('Deny')
    .click()
    .get('[data-cy="storage-in-transit"]')
    .should($div => {
      const text = $div.text();
      expect(text).to.include('Deny SIT Request');
      expect(text).to.not.include('Approve');
      expect(text).to.not.include('Edit');
    });
  cy.get('textarea[name="authorization_notes"]').type('this is a denial note');
  cy
    .get('[data-cy="storage-in-transit-deny-button"]')
    .contains('Deny')
    .click();

  cy.get('[data-cy="storage-in-transit-status-denied"]').contains('Denied');
  cy.get('[data-cy="sit-authorization-notes"]').contains('this is a denial note');

  cy.get('[data-cy="sit-edit-link"]').click();

  cy.get('textarea[name="authorization_notes"]').type('this is also a note', { force: true, delay: 150 });

  cy
    .get('[data-cy="sit-editor-save-button"]')
    .contains('Save')
    .click();

  cy.patientReload();

  cy.get('[data-cy="sit-authorization-notes"]').contains('this is also a note');
}

function officeUserGoesToPlacedSIT() {
  // Open all moves queue
  cy.patientVisit('/queues/all');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/all/);
  });

  // Find move (generated in e2ebasic.go) and open it
  cy.selectQueueItemMoveLocator('SITIN1');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  cy.get('[data-cy="hhg-tab"]').click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/hhg/);
  });
}

function officeUserEntitlementRemainingDays() {
  // Freeze the clock so we can test a specific remaining days.
  let now = new Date(Date.UTC(2019, 3, 10)).getTime(); // 4/10/2019
  cy.clock(now);

  officeUserGoesToPlacedSIT();

  cy
    .get('[data-cy=storage-in-transit-panel]')
    .should($div => {
      const text = $div.text();
      expect(text).to.include('In SIT');
      expect(text).to.include('Entitlement: 90 days (78 remaining)');
    })
    .get('[data-cy=storage-in-transit] [data-cy=sit-days-used]')
    .contains('12 days')
    .get('[data-cy=storage-in-transit] [data-cy=sit-expires]')
    .contains('28-Jun-2019');
}

function officeUserEntitlementRemainingDaysExpired() {
  // Freeze the clock so we can test a specific remaining days.
  let now = new Date(Date.UTC(2019, 6, 10)).getTime(); // 7/10/2019
  cy.clock(now);

  officeUserGoesToPlacedSIT();

  cy
    .get('[data-cy=storage-in-transit-panel]')
    .should($div => {
      const text = $div.text();
      expect(text).to.include('In SIT - SIT Expired');
      expect(text).to.include('Entitlement: 90 days (-13 remaining)');
    })
    .get('[data-cy=storage-in-transit] [data-cy=sit-days-used]')
    .contains('103 days')
    .get('[data-cy=storage-in-transit] [data-cy=sit-expires]')
    .contains('28-Jun-2019');
}
