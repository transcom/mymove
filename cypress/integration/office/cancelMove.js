/* global cy, */
describe('office user finds the move', () => {
  beforeEach(() => {
    cy.signIntoOffice();
  });
  it('office user cancels the move', () => {
    // Open the move
    cy.visit('/queues/new/moves/0db80bd6-de75-439e-bf89-deaafa1d0dc9/ppm');

    // Find the Cancel Move button
    cy
      .get('button')
      .contains('Cancel Move')
      .click();

    // Enter cancel reason
    cy
      .get('.cancel-panel textarea:first')
      .type('Canceling this move as a test!');
    cy
      .get('.cancel-panel button')
      .contains('Cancel Move')
      .click();
    cy
      .get('.cancel-panel button')
      .contains('Yes, Cancel Move')
      .click();
    cy
      .get('.usa-alert-success')
      .contains('Move #CANCEL for Submitted, PPM has been canceled');
  });
});
