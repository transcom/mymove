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

    // ensure the page pulls up the correct user
    cy.get('.ra-field-id > div > label')
      .first()
      .next()
      .then(($userId) => {
        cy.get('a').contains('Edit').click();
        cy.url().should('eq', adminBaseURL + '/system/users/' + $userId.text());
      });

    // check page content
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

    // deactivate the user
    cy.get('div[id="active"]').click();
    cy.get('#menu-active ul > li[data-value=false]').click();
    cy.get('button').contains('Save').click();

    // check that user was deactivated
    // note since we picked the first user before, we assume it's the first user now.
    // this should be true if sort order is not changed.
    cy.url().should('eq', adminBaseURL + '/system/users');
    cy.get('td.column-active span > span').first().should('have.attr', 'title', 'No');
  });
});
