import { adminBaseURL } from '../../support/constants';

describe('Offices Page', function () {
  before(() => {
    cy.prepareAdminApp();
  });

  it('successfully navigates to offices page', function () {
    cy.signInAsNewAdminUser();
    cy.get('a[href*="system/offices"]').click();
    cy.url().should('eq', adminBaseURL + '/system/offices');
    cy.get('header').contains('Offices');
    cy.get('[data-source=q]').should('be.visible');

    const columnLabels = ['Id', 'Name', 'Latitude', 'Longitude', 'Gbloc'];
    columnLabels.forEach((label) => {
      cy.get('table').contains(label);
    });
  });
});
