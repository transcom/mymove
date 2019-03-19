/* global cy */
import { milmoveAppName } from '../../support/constants';

describe('The Home Page', function() {
  beforeEach(() => {
    cy.setupBaseUrl(milmoveAppName);
  });
  it('successfully loads when not logged in', function() {
    cy.logout();
    milmoveUserIsOnSignInPage();
  });
  it('milmove user is properly welcomed', function() {
    milmoveUserIsWelcomed();
  });
  it('contains the correct phone number contact information', function() {
    cy.visit('/');
    cy.contains('(833) 645-6683');
  });
});

function milmoveUserIsOnSignInPage() {
  cy.contains('Welcome');
  cy.contains('Sign In');
}

function milmoveUserIsWelcomed() {
  cy.signIntoMyMoveAsUser('e10d5964-c070-49cb-9bd1-eaf9f7348eb7');
  cy.get('strong').contains('Welcome, PPM');
}
