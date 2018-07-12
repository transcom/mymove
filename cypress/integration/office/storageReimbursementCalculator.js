/* global cy, */
describe('office user finds the move', () => {
  beforeEach(() => {
    cy.signIntoOffice();
  });
  it('office user uses calculator', () => {
    // Open move ppm tab
    cy.visit('/queues/new/moves/0db80bd6-de75-439e-bf89-deaafa1d0dc8/ppm');
    // Click on PPM tab
    cy.get('.calculator-panel input[name="days_in_storage"]').type('30');
    cy.get('.calculator-panel input[name="weight"]').type('3141');

    cy.get('.calculator-panel [data-cy=calc]').click();

    cy
      .get('.calculator-panel .calculated-result')
      .contains('Maximum Obligation');

    cy.get('.calculator-panel [data-cy=reset]').click();
    cy
      .get('.calculator-panel input[name="days_in_storage"]')
      .should('be.empty');
    cy.get('.calculator-panel input[name="weight"]').should('be.empty');
  });
});
