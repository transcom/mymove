import { QAECSROfficeUserType } from '../../../support/constants';

describe('QAE/CSR Move Search', () => {
  before(() => {
    cy.prepareOfficeApp();
  });

  beforeEach(() => {
    cy.intercept('**/ghc/v1/swagger.yaml').as('getGHCClient');
    cy.intercept('**/ghc/v1/moves/search').as('getSearchResults');

    const userId = '2419b1d6-097f-4dc4-8171-8f858967b4db';
    cy.apiSignInAsUser(userId, QAECSROfficeUserType);
  });

  it('is able to search by move code', () => {
    const moveLocator = 'QAEHLP';
    cy.get('input[name="searchText"]').as('moveCodeSearchInput');
    cy.get('@moveCodeSearchInput').type(moveLocator).blur();
    cy.get('button').contains('Search').click();
    cy.wait(['@getSearchResults']);

    cy.get('tbody > tr').as('results');
    cy.get('@results').should('have.length', 1);
    cy.get('@results').first().contains(moveLocator);

    // click result to navigate to move details page
    cy.get('@results').first().click();
    cy.url().should('include', `/moves/${moveLocator}/details`);
  });
  it('is able to search by DOD ID', () => {
    const moveLocator = 'QAEHLP';
    const dodID = '1000000000';

    cy.get('[type="radio"]').check('dodID', { force: true });

    cy.get('input[name="searchText"]').type(dodID).blur();
    cy.get('button').contains('Search').click();
    cy.wait(['@getSearchResults']);

    cy.get('tbody > tr').as('results');
    cy.get('@results').should('have.length', 1);
    cy.get('@results').first().contains(moveLocator);
    cy.get('@results').first().contains(dodID);

    // click result to navigate to move details page
    cy.get('@results').first().click();
    cy.url().should('include', `/moves/${moveLocator}/details`);
  });
});
