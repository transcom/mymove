/* global cy */
describe('office user finds the move', function() {
  beforeEach(() => {
    cy.signIntoOffice();
  });
  it('office user views moves in queue new moves', function() {
    officeUserViewsMoves();
  });
  it('office user verifies the orders tab', function() {
    officeUserVerifiesOrders();
  });
  it('office user verifies the accounting tab', function() {
    officeUserVerifiesAccounting();
  });
  it('office user approves move, verifies and approves PPM', function() {
    officeUserApprovesMoveAndVerifiesPPM();
  });
});

function officeUserViewsMoves() {
  // Open new moves queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find move and open it
  cy.selectQueueItemMoveLocator('VGHEIS');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });
}

function officeUserVerifiesOrders() {
  // Open new moves queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find move and open it
  cy.selectQueueItemMoveLocator('VGHEIS');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  // Click on Orders document link, check that link matches
  cy
    .get('.panel-field')
    .contains('Orders (')
    // .find('a')
    .should('have.attr', 'href')
    .and('match', /^\/moves\/[^/]+\/orders/);

  // Click on edit orders
  cy
    .get('.editable-panel-header')
    .contains('Orders')
    .siblings()
    .click();

  // Enter details in form and save orders
  cy.get('input[name="orders.orders_number"]').type('666666');
  cy.get('select[name="orders.orders_type_detail"]').select('DELAYED_APPROVAL');

  cy.get('input[name="orders.orders_issuing_agency"]').type('ISSUING AGENCY');
  cy.get('input[name="orders.paragraph_number"]').type('FP-TP');

  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save')
    .click();

  // Verify data has been saved in the UI
  cy.get('span').contains('666666');

  cy.get('span').contains('ISSUING AGENCY');
  cy.get('span').contains('FP-TP');

  // Refresh browser and make sure changes persist
  cy.patientReload();
  cy
    .get('button')
    .contains('Approve PPM')
    .should('be.disabled');

  cy.get('span').contains('666666');
  cy.get('span').contains('Delayed Approval 20 Weeks or More');

  cy.get('span').contains('ISSUING AGENCY');
  cy.get('span').contains('FP-TP');
}

function officeUserVerifiesAccounting() {
  // Open new moves queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find move and open it
  cy.selectQueueItemMoveLocator('VGHEIS');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  // Enter details in form and save
  cy
    .get('.editable-panel-header')
    .contains('Accounting')
    .siblings()
    .click();

  cy.get('input[name="tac"]').type('6789');
  cy.get('select[name="department_indicator"]').select('AIR_FORCE');
  cy.get('input[name="sac"]').type('N002214CSW32Y9');

  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save')
    .click();

  // Verify data has been saved in the UI
  // added '.tac' as an additional selector since
  // '6789' is part of the DoD ID that was a false positive
  cy.get('.tac > span').contains('6789');
  cy.get('span').contains('N002214CSW32Y9');

  // Refresh browser and make sure changes persist
  cy.patientReload();
  cy
    .get('button')
    .contains('Approve PPM')
    .should('be.disabled');

  cy.get('span').contains('6789');
  cy.get('span').contains('57 - United States Air Force');
  cy.get('span').contains('N002214CSW32Y9');
}

function officeUserApprovesMoveAndVerifiesPPM() {
  // Open new moves queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find move and open it
  cy.selectQueueItemMoveLocator('VGHEIS');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  // Approve the move
  cy
    .get('button')
    .contains('Approve PPM')
    .should('be.disabled');
  cy
    .get('button')
    .contains('Approve Basics')
    .click();

  cy
    .get('button')
    .contains('Approve PPM')
    .should('be.enabled');

  // Open new moves queue
  cy.patientVisit('/');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Open PPMs Queue
  cy
    .get('span')
    .contains('PPMs')
    .click();

  // Find move and open it
  cy.selectQueueItemMoveLocator('VGHEIS');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  // Click on PPM tab
  cy
    .get('span')
    .contains('PPM')
    .click();
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/ppm/);
  });

  // Verify that the Estimates section contains expected data
  cy.get('span').contains('8,000');

  // Approve PPM
  cy
    .get('button')
    .contains('Approve PPM')
    .click();

  // Approve advance
  cy.get('.payment-table').within(() => {
    // Verify the status icon
    cy
      .get('td svg:first')
      .should('have.attr', 'title')
      .and('eq', 'Awaiting Review');
    // Verify the approve checkmark
    cy
      .get('td svg:last')
      .should('have.attr', 'title')
      .and('eq', 'Approve');

    // Approve advance and verify icon change
    cy.get('td svg:last').click();
    cy
      .get('td svg:first')
      .should('have.attr', 'title')
      .and('eq', 'Approved');
  });
}
