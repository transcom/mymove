/* global cy */
describe('office user finds the move', function() {
  beforeEach(() => {
    Cypress.config('baseUrl', 'http://officelocal:4000');
    cy.signInAsOfficeUser();
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
    .contains('H9VSGX')
    .dblclick();

  // Verify basics and edit
  // tbd

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
  cy
    .get('button')
    .contains('Approve PPM')
    .click();

  // Go back to new moves queue
  cy
    .get('span')
    .contains('New Moves Queue')
    .click();

  // Open PPMs
  cy
    .get('span')
    .contains('PPMs')
    .click();

  // Find move and open it
  cy
    .get('div')
    .contains('H9VSGX')
    .dblclick();

  // Verify PPM
  // tbd
}
