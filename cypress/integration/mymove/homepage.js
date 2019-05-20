/* global cy */
import { milmoveAppName } from '../../support/constants';

describe('The Home Page', function() {
  beforeEach(() => {
    cy.setupBaseUrl(milmoveAppName);
  });
  it('creates new devlocal user', function() {
    cy.signInAsNewMilMoveUser();
  });
  it('successfully loads when not logged in', function() {
    cy.logout();
    milmoveUserIsOnSignInPage();
  });
  it('milmove user is properly welcomed', function() {
    milmoveUserIsWelcomed();
  });
  it('contains the link to customer service', function() {
    cy.visit('/');
    cy.get('address').within(() => {
      cy.contains('Contact Us');
      cy.get('a').should('have.attr', 'href', 'https://move.mil/customer-service');
    });
  });
});

function milmoveUserIsOnSignInPage() {
  cy.contains('Welcome');
  cy.contains('Sign In');
}

function milmoveUserIsWelcomed() {
  cy.signIntoMyMoveAsUser('e10d5964-c070-49cb-9bd1-eaf9f7348eb7');
  cy.contains('Welcome, PPM');
}
