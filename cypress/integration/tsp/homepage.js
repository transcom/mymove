/* global cy */
import { tspAppName } from '../../support/constants';

describe('TSP Home Page', function() {
  beforeEach(() => {
    cy.setupBaseUrl(tspAppName);
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
  cy.patientVisit('/');
}

function tspUserIsOnSignInPage() {
  cy.contains('tsp.move.mil');
  cy.contains('Sign In');
}

function tspUserIsWelcomed() {
  cy.signIntoTSP();
  cy.get('strong').contains('Welcome, Leo');
}
