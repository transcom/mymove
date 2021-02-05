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

    // check that the user's login.gov email is shown in the page title
    cy.get('.ra-field-loginGovEmail span.MuiTypography-root')
      .invoke('text')
      .then((loginGovEmail) => {
        cy.get('#react-admin-title').contains(loginGovEmail);
      });

    const pageContent = ['User ID', 'User email', 'Active', 'Created at', 'Updated at'];
    pageContent.forEach((label) => {
      cy.get('label').contains(label);
    });
  });
});
