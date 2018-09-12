/* global cy, */
describe('office user finds the move', () => {
  beforeEach(() => {
    // We need to reset the test database because we're finding a move by a
    // specific ID. If the test runs twice in a row, the move will be cancelled
    // and unable to be canceled again.
    cy.exec('make db_e2e_reset');
    cy.signIntoOffice();
  });
  afterEach(() => {
    // Reset the database for future tests that need this move to not be
    // canceled.
    cy.exec('make db_e2e_reset');
  });

  it('office user cancels the move', () => {
    // Open the move
    cy.visit('/queues/new/moves/0db80bd6-de75-439e-bf89-deaafa1d0dc8/ppm');

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
      .contains('Cancel entire move')
      .click();
    cy
      .get('.cancel-panel button')
      .contains('Yes, cancel move')
      .click();
    cy
      .get('.usa-alert-success')
      .contains('Move #VGHEIS for Submitted, PPM has been canceled');
  });
});
