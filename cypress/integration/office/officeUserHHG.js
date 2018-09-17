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
  cy
    .get('div')
    .contains('RLKBEM')
    .dblclick();

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
  cy.visit('/queues/hhg_accepted');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/hhg_accepted/);
  });

  // Find move (generated in e2ebasic.go) and open it
  cy
    .get('div')
    .contains('BACON3')
    .dblclick();

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

function officeUserApprovesHHG() {
  // Open accepted hhg queue
  cy.visit('/queues/hhg_accepted');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/hhg_accepted/);
  });

  // Find move and open it
  cy
    .get('div')
    .contains('BACON5')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  // Approve basics
  cy
    .get('button')
    .contains('Approve Basics')
    .click();

  // disabled because not on hhg tab
  cy
    .get('button')
    .contains('Approve Shipment')
    .should('be.disabled');

  cy.get('.status').contains('Accepted');

  // Click on HHG tab
  cy
    .get('span')
    .contains('HHG')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/hhg/);
  });

  // Approve HHG
  cy
    .get('button')
    .contains('Approve Shipment')
    .click();

  cy
    .get('button')
    .contains('Approve Shipment')
    .should('be.disabled');

  cy.get('.status').contains('Approved');
}
