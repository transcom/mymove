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
  it('office user with approved move completes delivered HHG shipment', function() {
    officeUserCompletesHHG();
  });
  // Commenting this out for now since unbilled invoice line items are still in development
  // it('office user with completed move approve payment for invoice (sends invoice)', function() {
  //   officeUserApprovePaymentInvoice();
  // });
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

function officeUserViewsDeliveredShipment() {
  // Open new moves queue
  cy.patientVisit('/queues/hhg_delivered');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/hhg_delivered/);
  });

  // Find move (generated in e2ebasic.go) and open it
  cy
    .get('div')
    .contains('SCHNOO')
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

function officeUserViewsCompletedShipment() {
  // Open new moves queue
  cy.patientVisit('/queues/hhg_completed');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/hhg_completed/);
  });

  // Find move (generated in e2ebasic.go) and open it
  cy
    .get('div')
    .contains('NOCHKA')
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
  cy.patientVisit('/queues/hhg_accepted');
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

function officeUserApprovesOnlyBasicsHHG() {
  // Open accepted hhg queue
  cy.patientVisit('/queues/new');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find move and open it
  cy
    .get('div')
    .contains('BACON6')
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
  cy
    .get('button')
    .contains('Complete Shipments')
    .should('be.disabled');

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
  cy
    .get('button')
    .contains('Approve Shipment')
    .should('be.disabled');

  // Disabled because already approved and not delivered
  cy
    .get('button')
    .contains('Approve Shipment')
    .should('be.disabled');
  cy
    .get('button')
    .contains('Complete Shipments')
    .should('be.disabled');

  cy.get('.status').contains('Awarded');
}

function officeUserApprovesHHG() {
  // Open accepted hhg queue
  cy.patientVisit('/queues/hhg_accepted');
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
  cy
    .get('button')
    .contains('Complete Shipments')
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

  // Disabled because already approved and not delivered
  cy
    .get('button')
    .contains('Approve Shipment')
    .should('be.disabled');
  cy
    .get('button')
    .contains('Complete Shipments')
    .should('be.disabled');

  cy.get('.status').contains('Approved');
}

function officeUserCompletesHHG() {
  // Open delivered hhg queue
  cy.patientVisit('/queues/hhg_delivered');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/hhg_delivered/);
  });

  // Find move and open it
  cy
    .get('div')
    .contains('SSETZN')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  // Basics Approved
  cy.get('.status').contains('Approved');

  // Click on HHG tab
  cy
    .get('span')
    .contains('HHG')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/hhg/);
  });

  // Complete HHG
  cy.get('.status').contains('Delivered');

  cy
    .get('button')
    .contains('Complete Shipments')
    .click();

  cy
    .get('button')
    .contains('Approve Shipment')
    .should('be.disabled');

  cy.get('.status').contains('Completed');
}

function officeUserApprovePaymentInvoice() {
  // Open completed hhg queue
  cy.patientVisit('/queues/hhg_delivered');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/hhg_delivered/);
  });

  // Find move and open it
  cy
    .get('div')
    .contains('DOOB')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  // Basics Approved
  cy.get('.status').contains('Approved');

  // Click on HHG tab
  cy
    .get('span')
    .contains('HHG')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/hhg/);
  });

  // Submit Invoice for HHG
  cy.get('.status').contains('Delivered');

  cy
    .get('.invoice-panel button')
    .contains('Approve Payment')
    .should('be.disabled');
  // .click();  TODO: figure out how not to make call to GEX
}
