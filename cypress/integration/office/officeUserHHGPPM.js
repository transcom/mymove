/* global cy */
describe('office user finds the shipment', function() {
  beforeEach(() => {
    cy.signIntoOffice();
  });
  it('office user approves shipment', function() {
    officeUserVisitsAllMovesQueue();
    officeUserViewsMove();
    officeUserVisitsPPMTab();
    officeUserVisitsHHGTab();
    officeUserApprovesShipment();
  });
});

function officeUserApprovesShipment() {
  const ApproveShipmentButton = cy.get('button').contains('Approve HHG');

  ApproveShipmentButton.should('be.enabled');

  ApproveShipmentButton.click();

  ApproveShipmentButton.should('be.disabled');

  cy.get('.status').contains('Approved');
}

function officeUserVisitsAllMovesQueue() {
  cy.patientVisit('/queues/all');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/all/);
  });
}

function officeUserVisitsPPMTab() {
  // navtab
  cy
    .get('a')
    .contains('PPM')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/ppm/);
  });
}

function officeUserVisitsHHGTab() {
  // navtab
  cy
    .get('a')
    .contains('HHG')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/hhg/);
  });
}

function officeUserViewsMove() {
  // Find move (generated in e2ebasic.go) and open it
  cy
    .get('div')
    .contains('COMBO2')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });
}
