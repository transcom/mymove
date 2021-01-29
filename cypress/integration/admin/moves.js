import { adminBaseURL } from '../../support/constants';

describe('Moves Page', function () {
  before(() => {
    cy.prepareAdminApp();
  });

  it('successfully navigates to Moves page', function () {
    cy.signInAsNewAdminUser();
    cy.get('a[href*="system/moves"]').click();
    cy.url().should('eq', adminBaseURL + '/system/moves');
    cy.get('header').contains('Moves');

    const columnLabels = [
      'Id',
      'Order Id',
      'Service Member Id',
      'Locator',
      'Status',
      'Show',
      'Created at',
      'Updated at',
    ];
    columnLabels.forEach((label) => {
      cy.get('table').contains(label);
    });
  });
});

describe('Moves Details Show Page', function () {
  before(() => {
    cy.prepareAdminApp();
  });

  it('pulls up details page for a user', function () {
    cy.signInAsNewAdminUser();
    cy.get('a[href*="system/moves"]').click();
    cy.url().should('eq', adminBaseURL + '/system/moves');
    cy.get('span[reference="moves"]').first().click();

    const labels = [
      'Id',
      'Locator',
      'Status',
      'Show',
      'Order Id',
      'Created at',
      'Updated at',
      'Service member Id',
      'Service member first name',
      'Service member middle name',
      'Service member last name',
    ];
    labels.forEach((label) => {
      cy.get('.MuiCardContent-root label').contains(label);
    });
  });
});
