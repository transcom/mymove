/* global cy */
import { officeAppName } from '../../support/constants';

describe('Office Home Page', function() {
  beforeEach(() => {
    cy.setupBaseUrl(officeAppName);
  });
  it('creates new devlocal user', function() {
    cy.signInAsNewOfficeUser();
  });
  it('successfully loads when not logged in', function() {
    cy.logout();
    officeUserIsOnSignInPage();
  });
  it('office user is properly welcomed', function() {
    officeUserIsWelcomed();
  });
  it('open accepted shipments queue and see moves', function() {
    cy.signIntoOffice();
    officeAllMoves();
  });
  it('office user can use a single click to view move info', function() {
    cy.waitForReactTableLoad();

    cy.get('[data-cy=queueTableRow]:first').click();
    cy.url().should('include', '/moves/');
  });
});

function officeUserIsOnSignInPage() {
  cy.contains('office.move.mil');
  cy.contains('Sign In');
}

function officeUserIsWelcomed() {
  cy.signIntoOffice();
  cy.get('strong').contains('Welcome, Leo');
}

function officeAllMoves() {
  cy.patientVisit('/queues/all');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/all/);
  });

  cy
    .get('[data-cy=locator]')
    .contains('NOSHOW')
    .should('not.exist');
}
