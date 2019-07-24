/* global cy */
describe('Office ppm queue', () => {
  beforeEach(() => {
    cy.signIntoOffice();
    cy.get('[data-cy=ppm-queue]').click();
  });

  it('does not have a GBL column', checkForGBLColumn);
});

function checkForGBLColumn() {
  cy.contains('GBL').should('not.exist');
}
