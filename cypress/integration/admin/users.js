import { adminBaseURL } from '../../support/constants';

describe('Users Page', function () {
  before(() => {
    cy.prepareAdminApp();
  });

  it('successfully navigates to users page', function () {
    cy.signInAsNewAdminUser();
    cy.get('a[href*="system/users"]').click();
    cy.url().should('eq', adminBaseURL + '/system/users');
    cy.get('header').contains('Users');

    const columnLabels = ['Id', 'Email', 'Active', 'Created at'];
    columnLabels.forEach((label) => {
      cy.get('table').contains(label);
    });
  });
});

describe('Users Details Show Page', function () {
  before(() => {
    cy.prepareAdminApp();
  });

  it('pulls up details page for a user', function () {
    cy.signInAsNewAdminUser();
    failOnStatusCode: false;
    cy.get('a[href*="system/users"]').click();
    cy.url().should('eq', adminBaseURL + '/system/users');
    cy.get('[data-testid="user-id"]').first().click();

    const pageContent = ['User ID', 'User email', 'Active', 'Created at', 'Updated at'];
    pageContent.forEach((label) => {
      cy.get('label').contains(label);
    });
  });
});

describe('Users Details Edit Page', function () {
  before(() => {
    cy.prepareAdminApp();
  });

  it('pulls up edit page for a user', function () {
    cy.signInAsNewAdminUser();
    failOnStatusCode: false;
    cy.get('a[href*="system/users"]').click();
    cy.url().should('eq', adminBaseURL + '/system/users');
    cy.get('[data-testid="user-id"]').first().click();
    cy.get('a').contains('Edit').click();

    const pageContent = [
      'Id',
      'Login gov email',
      'Active',
      'Revoke admin session',
      'Revoke office session',
      'Revoke mil session',
    ];
    pageContent.forEach((label) => {
      cy.get('label').contains(label);
    });
  });
});
