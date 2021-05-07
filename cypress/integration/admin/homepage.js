import { adminBaseURL } from '../../support/constants';

describe('Admin Home Page', function () {
  before(() => {
    cy.prepareAdminApp();
  });

  it('creates new devlocal user', function () {
    cy.signInAsNewAdminUser();
  });

  it('successfully loads when not logged in', function () {
    cy.logout();
    cy.contains('admin.move.mil');
    cy.contains('Sign In');
  });

  it('redirects to the office users page after login', function () {
    cy.signInAsNewAdminUser();
    cy.contains('Office Users');
    cy.url().should('eq', adminBaseURL + '/system/office_users');
  });

  it('admin user can logout and see logout success message', function () {
    cy.contains('Logout').click();
    cy.location().should((loc) => {
      expect(loc.pathname).to.match('/^//');
    });
    cy.url().should('eq', '/');
    cy.get('.usa-alert--success').contains('You have signed out of MilMove').should('exist');
  });
});
