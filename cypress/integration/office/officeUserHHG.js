/* global cy */
describe('office user finds the shipment', function() {
  beforeEach(() => {
    cy.signIntoOffice();
  });
  it('office user views hhg moves in queue new moves', function() {
    officeUserViewsMoves();
  });
  it('office user views accepted hhg moves in queue Accepted HHGs', function() {
    officeUserViewsAcceptedShipment();
  });
  it('office user views delivered hhg moves in queue Delivered HHGs', function() {
    officeUserViewsDeliveredShipment();
  });
  it('office user views completed hhg moves in queue Completed HHGs', function() {
    officeUserViewsCompletedShipment();
  });
  it('office user approves basics for move, cannot approve HHG shipment', function() {
    officeUserApprovesOnlyBasicsHHG();
  });
  it('office user approves basics for move, verifies and approves HHG shipment', function() {
    officeUserApprovesHHG();
  });
});

function officeUserViewsMoves() {
  // Open new moves queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find move (generated in e2ebasic.go) and open it
  cy.selectQueueItemMoveLocator('RLKBEM');

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
}

function officeUserViewsDeliveredShipment() {
  // Open new moves queue
  cy.patientVisit('/queues/hhg_delivered');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/hhg_delivered/);
  });

  // Find move (generated in e2ebasic.go) and open it
  cy.selectQueueItemMoveLocator('SCHNOO');

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
}

function officeUserViewsCompletedShipment() {
  // Open new moves queue
  cy.patientVisit('/queues/hhg_completed');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/hhg_completed/);
  });

  // Find move (generated in e2ebasic.go) and open it
  cy.selectQueueItemMoveLocator('NOCHKA');

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
}

function officeUserViewsAcceptedShipment() {
  // Open new moves queue
  cy.patientVisit('/queues/hhg_accepted');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/hhg_accepted/);
  });

  // Find move (generated in e2ebasic.go) and open it
  cy.selectQueueItemMoveLocator('BACON3');

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

  // Since the shipment hasn't been picked up yet, it can be cancelled
  cy
    .get('button')
    .contains('Cancel Move')
    .should('not.be.disabled');
}

function officeUserApprovesOnlyBasicsHHG() {
  // Open accepted hhg queue
  cy.patientVisit('/queues/new');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find move and open it
  cy.selectQueueItemMoveLocator('BACON6');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  cy.get('.combo-button').click();

  // Approve basics
  cy
    .get('.combo-button .dropdown')
    .contains('Approve Basics')
    .click();

  cy.get('.combo-button').click();

  cy.get('.status').contains('Approved');

  // Click on HHG tab
  cy
    .get('span')
    .contains('HHG')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/hhg/);
  });

  // disabled because shipment not yet accepted
  cy.get('.combo-button').click();

  // Disabled because already approved and not delivered
  cy
    .get('.combo-button .dropdown')
    .contains('Approve HHG')
    .should('have.class', 'disabled');

  cy.get('.combo-button').click();

  cy.get('.status').contains('Awarded');
}

function officeUserApprovesHHG() {
  // Open accepted hhg queue
  cy.patientVisit('/queues/hhg_accepted');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/hhg_accepted/);
  });

  // Find move and open it
  cy.selectQueueItemMoveLocator('BACON5');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  cy.get('.combo-button').click();

  // Approve basics
  cy
    .get('.combo-button .dropdown')
    .contains('Approve Basics')
    .click();

  cy.get('.combo-button').click();

  cy.get('.status').contains('Accepted');

  // Click on HHG tab
  cy
    .get('span')
    .contains('HHG')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/hhg/);
  });

  cy.get('.combo-button').click();

  // Approve HHG

  cy
    .get('.combo-button .dropdown')
    .contains('Approve HHG')
    .click();

  cy.get('.status').contains('Approved');
}
