import { adminBaseURL } from '../../support/constants';

describe('Admin Customer Users Page', function () {
  before(() => {
    cy.prepareAdminApp();
  });

  it('successfully navigates to customer users page', function () {
    cy.signInAsNewAdminUser();
    cy.get('a[href*="customer_users"]').click();
    cy.url().should('eq', adminBaseURL + '/system/customer_users');
    cy.get('header').contains('Customer users');

    const columnLabels = ['Id', 'Email', 'Active', 'Created at'];
    columnLabels.forEach((label) => {
      cy.get('table').contains(label);
    });
  });
});
