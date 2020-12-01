describe('Office ppm queue', () => {
  before(() => {
    cy.prepareOfficeApp();
  });

  beforeEach(() => {
    cy.signIntoOffice();
    cy.get('[data-testid=ppm-queue]').click();
  });

  it('does not have a GBL column', () => {
    cy.contains('GBL').should('not.exist');
  });
});
