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
    cy.contains('Office users');
    cy.url().should('eq', adminBaseURL + '/system/office_users');
  });
});
