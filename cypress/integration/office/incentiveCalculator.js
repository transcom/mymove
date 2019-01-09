/* global cy */
describe('office user uses incentive calculator', () => {
  beforeEach(() => {
    cy.signIntoOffice();
  });
  it('finds calculator and executes it', () => {
    // Open move ppm tab
    cy.patientVisit('/queues/new/moves/0db80bd6-de75-439e-bf89-deaafa1d0dc8/ppm');
    // Click on PPM tab
    cy.get('.incentive-calc').within(() => {
      cy
        .get('input[name="planned_move_date"]')
        .first()
        .clear()
        .type('9/2/2018{enter}')
        .blur();

      cy.get('input[name="weight"]').type('3141');

      cy.get('[data-cy=calc]').click();

      cy.get('.calculated-result').contains('PPM Incentive');

      cy.get('[data-cy=reset]').click();
      cy.get('input[name="weight"]').should('be.empty');
    });
  });
});
