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
    cy.get('a[href*="system/users"]').click();
    cy.url().should('eq', adminBaseURL + '/system/users');
    cy.get('[data-testid="user-id"]').first().click();

    const pageContent = ['user ID', 'user email', 'Active', 'Created at', 'Updated at'];
    pageContent.forEach((label) => {
      cy.get('label').contains(label);
    });
  });
});
