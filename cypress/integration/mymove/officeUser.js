/* global cy */
describe('office user finds the move', function() {
  beforeEach(() => {
    cy.signInAsOfficeUser();
  });
  it('office user views moves in queue new moves', function() {
    officeUserViewsMoves();
  });
});

function officeUserViewsMoves() {
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });
  cy
    .get('div')
    .contains('H9VSGX')
    .dblclick();
  cy
    .get('button')
    .contains('Approve PPM')
    .should('be.disabled');
  cy
    .get('button')
    .contains('Approve Basics')
    .click();
  // cy.get('button').contains('Approve PPM').should('.not.be.disabled');
}
