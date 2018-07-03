/* global cy */
describe('office user finds the move', function() {
  beforeEach(() => {
    Cypress.config('baseUrl', 'http://officelocal:4000');
    cy.signInAsUser('9bfa91d2-7a0c-4de0-ae02-b8cf8b4b858b');
  });
  it('office user views moves in queue new moves', function() {
    officeUserViewsMoves();
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

  // Verify basics and edit
  // Click on Orders document link, check that html is rendered
  cy
    .get('.document')
    .find('a')
    .then(function($a) {
      console.log($a);
      const href = $a.prop('href');
      cy
        .request(href)
        .its('body')
        .should('include', '</html>');
    });

  // Click on edit orders
  cy
    .get('.editable-panel-header')
    .contains('Orders')
    .siblings()
    .click();

  // Enter details in form and save orders
  cy.get('input[name="orders.orders_number"]').type('12345');
  cy.get('select[name="orders.orders_type_detail"]').select('DELAYED_APPROVAL');
  cy
    .get('button')
    .contains('Save')
    .click();

  // Enter details in form and save
  cy
    .get('.editable-panel-header')
    .contains('Accounting')
    .siblings()
    .click({ force: true });

  cy.get('input[name="tac"]').type('12345');

  cy
    .get('.editable-panel-header')
    .contains('Accounting')
    .siblings()
    .click({ force: true });

  cy.get('select[name="department_indicator"]').select('AIR_FORCE');
  cy.get('input[name="tac"]').type('12345');

  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save')
    .click();

  // Refresh browser and make sure changes persist
  cy.visit('/');

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
  // cy
  //   .get('button')
  //   .contains('Approve PPM')
  //   .should('be.enabled');
  // cy
  //   .get('button')
  //   .contains('Approve PPM')
  //   .click();

  // // Go back to new moves queue
  // cy
  //   .get('span')
  //   .contains('New Moves Queue')
  //   .click();

  // // Open PPMs
  // cy
  //   .get('span')
  //   .contains('PPMs')
  //   .click();

  // // Find move and open it
  // cy
  //   .get('div')
  //   .contains('VGHEIS')
  //   .dblclick();

  // Verify PPM
  // tbd
}
