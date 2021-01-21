import { adminBaseURL } from '../../support/constants';

describe('Admin Customer Users Page', function () {
  before(() => {
    cy.prepareAdminApp();
  });

  it('can successfully navigate to customer users page', function () {
    cy.signInAsNewAdminUser();
    cy.get('a[href*="customer_users"]').click();
    cy.url().should('eq', adminBaseURL + '/system/customer_users');
  });
});
