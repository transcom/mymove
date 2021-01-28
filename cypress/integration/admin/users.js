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
  // after(() => {

  // });

  // Is there a way to pass in any id instead of a specific ID in the URL
  // pullout columnLabels var
  //  Make a note that this test is dependent on the previous test
  it('pulls up details page for a user', function () {
    cy.signInAsNewAdminUser();
    cy.get('a[href*="system/users"]').click();
    cy.url().should('eq', adminBaseURL + '/system/users');
    // Fix the line below to get any id for any user in the list of users
    // cy.get('a[href*="system/users/${id.first()}/show"]').click();
    cy.get('a[href*="/show"]').first().click();

    const pageContent = ['user ID', 'user email', 'Active', 'Created at', 'Updated at'];
    pageContent.forEach((string) => {
      cy.get('input').contains(string);
    });
  });
});
