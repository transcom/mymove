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
  it('tsp user sees all showable moves', function() {
    cy.signIntoTSP();
    tspAllMoves();
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

function tspAllMoves() {
  cy.patientVisit('/queues/all');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/all/);
  });

  cy
    .get('[data-cy=locator]')
    .contains('NOSHOW')
    .should('not.exist');
}
