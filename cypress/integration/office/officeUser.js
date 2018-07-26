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
    cy.resetDb();
  });
});

function officeUserViewsMoves() {
  // Open new moves queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find move and open it
  cy
    .get('div')
    .contains('VGHEIS')
    .dblclick();

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
  cy
    .get('div')
    .contains('VGHEIS')
    .dblclick();

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

  // Refresh browser and make sure changes persist
  cy.reload();
  cy
    .get('button')
    .contains('Approve PPM')
    .should('be.disabled');

  cy.get('span').contains('666666');
  cy.get('span').contains('Delayed Approval 20 Weeks or More');
}

function officeUserVerifiesAccounting() {
  // Open new moves queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find move and open it
  cy
    .get('div')
    .contains('VGHEIS')
    .dblclick();

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

  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save')
    .click();

  // Verify data has been saved in the UI
  cy.get('span').contains('6789');

  // Refresh browser and make sure changes persist
  cy.reload();
  cy
    .get('button')
    .contains('Approve PPM')
    .should('be.disabled');

  cy.get('span').contains('6789');
  cy.get('span').contains('57 - United States Air Force');
}

function officeUserApprovesMoveAndVerifiesPPM() {
  // Open new moves queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find move and open it
  cy
    .get('div')
    .contains('VGHEIS')
    .dblclick();

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
  cy.visit('/');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Open PPMs Queue
  cy
    .get('span')
    .contains('PPMs')
    .click();

  // Find move and open it
  cy
    .get('div')
    .contains('VGHEIS')
    .dblclick();

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
  cy.get('span').contains('8000');

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
