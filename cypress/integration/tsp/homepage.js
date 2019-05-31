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
  it('office user can use a single click to view shipment info', function() {
    cy.signIntoTSP();
    cy.waitForReactTableLoad();

    cy.get('[data-cy=queueTableRow]:first').click();
    cy.url().should('include', '/shipments/');
  });
});

describe('Queue staleness indicator', () => {
  it('displays the correct time ago text', () => {
    cy.clock();
    cy.setupBaseUrl(tspAppName);
    cy.signIntoOffice();
    cy.patientVisit('/queues/all');

    cy.get('[data-cy=staleness-indicator]').should('have.text', 'Last updated a few seconds ago');

    cy.tick(120000);

    cy.get('[data-cy=staleness-indicator]').should('have.text', 'Last updated 2 mins ago');
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
