/* global cy */
describe('Office ppm queue', () => {
  beforeEach(() => {
    cy.signIntoOffice();
    cy.get('[data-cy=ppm-queue]').click();
  });

  it('does not have a GBL column', checkForGBLColumn);
  it('shows a clock icon for records with PAYMENT_REQUESTED or SUBMITTED status', checkIcons);
});

function checkForGBLColumn() {
  cy.contains('GBL').should('not.exist');
}

function checkIcons() {
  cy.wait(1000);
  cy.get('[data-cy=ppm-queue-icon]').each((el, index, list) => {
    cy
      .wrap(el)
      .closest('.rt-td')
      .next('.rt-td')
      .find('[data-cy=status]')
      .should(status => {
        expect(status.text()).to.match(/Payment requested|Submitted/);
      });
  });
}
