/* global cy */

describe('office user finds the shipment', function() {
  beforeEach(() => {
    cy.signIntoOffice();
  });
  const queues = {
    new: 'new',
    ppm: 'ppm',
    acceptedHHG: 'hhg_accepted',
    deliveredHHG: 'hhg_delivered',
    completedHHG: 'hhg_completed',
    all: 'all',
  };
  it('office user sees accepted HHG in Accepted HHGs queue', function() {
    officeUserVisitsQueue(queues.acceptedHHG);
    officeUserViewsMove('COMBO2');
  });
  it('office user sees delivered HHG in Delivered HHG queue', function() {
    officeUserVisitsQueue(queues.deliveredHHG);
    officeUserViewsMove('COMBO3');
  });
  it('office user sees completed HHG in Completed HHGs queue', function() {
    officeUserVisitsQueue(queues.completedHHG);
    officeUserViewsMove('COMBO4');
  });
  it('office user approves shipment', function() {
    officeUserVisitsQueue(queues.all);
    officeUserViewsMove('COMBO2');
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

function officeUserVisitsQueue(queue) {
  const queueName = queue.toLowerCase();
  // eslint-disable-next-line security/detect-non-literal-regexp
  const routePattern = new RegExp(`^/queues/${queueName}`);
  cy.patientVisit(`/queues/${queueName}`);

  cy.location().should(loc => {
    expect(loc.pathname).to.match(routePattern);
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

function officeUserViewsMove(locator) {
  cy
    .get('div')
    .contains(locator)
    .should('exist')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });
}
