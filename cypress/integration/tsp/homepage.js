/* global cy, Cypress */
describe('TSP Home Page', function() {
  beforeEach(() => {
    Cypress.config('baseUrl', 'http://tsplocal:4000');
  });
  it('successfully loads when not logged in', function() {
    tspUserLogsOut();
    tspUserIsOnSignInPage();
  });
  it('tsp user is properly welcomed', function() {
    tspUserIsWelcomed();
  });
});

function tspUserLogsOut() {
  // Logs out any users
  cy.logout();
  cy.visit('/');
}

function tspUserIsOnSignInPage() {
  cy.contains('tsp.move.mil');
  cy.contains('Sign In');
}

function tspUserIsWelcomed() {
  cy.signIntoTSP();
  cy.get('strong').contains('Welcome, Leo');
}
