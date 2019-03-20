/* global cy */
import { officeAppName } from '../../support/constants';

describe('Office Home Page', function() {
  beforeEach(() => {
    cy.setupBaseUrl(officeAppName);
  });
  it('successfully loads when not logged in', function() {
    cy.logout();
    officeUserIsOnSignInPage();
  });
  it('office user is properly welcomed', function() {
    officeUserIsWelcomed();
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
