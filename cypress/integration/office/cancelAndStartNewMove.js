/* global cy, */
describe('office user finds the move', () => {
  it('office user cancels the move', () => {
    cy.signIntoOffice();

    // Open the move
    cy.patientVisit('/queues/new/moves/0db80bd6-de75-439e-bf89-deaafa1d0dc9/ppm');

    // Find the Cancel Move button
    cy
      .get('button')
      .contains('Cancel Move')
      .click();

    // Enter cancel reason
    cy.get('.cancel-panel textarea:first').type('Canceling this move as a test!');
    cy
      .get('.cancel-panel button')
      .contains('Cancel Move')
      .click();
    cy
      .get('.cancel-panel button')
      .contains('Yes, Cancel Move')
      .click();
    cy.get('.usa-alert-success').contains('Move #CANCEL for Submitted, PPM has been canceled');
  });
  it('Service Member starts a new move after previous move is canceled.', () => {
    cy.signIntoMyMoveAsUser('e10d5964-c070-49cb-9bd1-eaf9f7348eb7');

    // Landing page contains move canceled alert message
    cy.contains('Your move was canceled');

    // User clicks on Start a new move and proceeds to orders page
    cy
      .get('button')
      .contains('Start')
      .click();

    cy.get('h1').contains('Review your Profile');

    cy.nextPage();

    cy.get('h1').contains('Tell Us About Your Move Orders');
  });
});
