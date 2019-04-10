/* global cy */
import { tspAppName } from '../../support/constants';

describe('TSP Home Page', function() {
  beforeEach(() => {
    cy.setupBaseUrl(tspAppName);
  });
  it('creates new devlocal user', function() {
    cy.signInAsNewTSPUser();
  });
  it('successfully loads when not logged in', function() {
    cy.logout();
    tspUserIsOnSignInPage();
  });
  it('tsp user is properly welcomed', function() {
    tspUserIsWelcomed();
  });
});

function tspUserIsOnSignInPage() {
  cy.contains('tsp.move.mil');
  cy.contains('Sign In');
}

function tspUserIsWelcomed() {
  cy.signIntoTSP();
  cy.get('strong').contains('Welcome, Leo');
}
