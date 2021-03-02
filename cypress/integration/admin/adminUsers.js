import { adminBaseURL } from '../../support/constants';

describe('Admin Users List Page', function () {
  before(() => {
    cy.prepareAdminApp();
  });

  it('successfully navigates to Admin Users page', function () {
    cy.signInAsNewAdminUser();
    cy.get('a[href*="system/admin_users"]').click();
    cy.url().should('eq', adminBaseURL + '/system/admin_users');
    cy.get('header').contains('Admin users');

    const columnLabels = ['Id', 'Email', 'First name', 'Last name', 'User Id', 'Active'];
    columnLabels.forEach((label) => {
      cy.get('table').contains(label);
    });
  });
});

describe('Admin Users Show Page', function () {
  before(() => {
    cy.prepareAdminApp();
  });

  it('pulls up details page for an admin user', function () {
    cy.signInAsNewAdminUser();
    cy.get('a[href*="system/admin_users"]').click();
    cy.url().should('eq', adminBaseURL + '/system/admin_users');
    cy.get('tr[resource="admin_users"]').first().click();

    // check that the admin user's name is shown in the page title
    cy.get('.ra-field-firstName span.MuiTypography-root')
      .invoke('text')
      .then((firstName) => {
        cy.get('.ra-field-lastName span.MuiTypography-root')
          .invoke('text')
          .then((lastName) => {
            cy.get('#react-admin-title').should('contain', firstName + ' ' + lastName);
          });
      });

    const labels = [
      'Id',
      'Email',
      'First name',
      'Last name',
      'User Id',
      'Organization',
      'Active',
      'Created at',
      'Updated at',
    ];
    labels.forEach((label) => {
      cy.get('.MuiCardContent-root label').contains(label);
    });
  });
});

describe('Admin Users Edit Page', function () {
  before(() => {
    cy.prepareAdminApp();
  });

  it('pulls up edit page for an admin user', function () {
    cy.signInAsNewAdminUser();
    cy.get('a[href*="system/admin_users"]').click();
    cy.url().should('eq', adminBaseURL + '/system/admin_users');
    cy.get('tr[resource="admin_users"]').first().click();

    // grab the admin user's ID to check that the correct value is in the url
    cy.get('.ra-field-id span.MuiTypography-root')
      .invoke('text')
      .as('adminUserID')
      .then((adminUserID) => {
        // continue to the edit page
        cy.get('a').contains('Edit').click();
        cy.url().should('eq', adminBaseURL + '/system/admin_users/' + adminUserID);
      });

    const disabledFields = ['id', 'email', 'userId', 'createdAt', 'updatedAt'];
    disabledFields.forEach((label) => {
      cy.get('[id="' + label + '"]').should('be.disabled');
    });

    // todo: save update
  });
});
