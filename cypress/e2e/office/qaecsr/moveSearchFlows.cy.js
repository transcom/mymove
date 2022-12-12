import { QAECSROfficeUserType } from '../../../support/constants';

const moveCode = 'QAEHLP';
describe('QAE/CSR Move Search', () => {
  before(() => {
    cy.prepareOfficeApp();
  });

  beforeEach(() => {
    cy.intercept('**/ghc/v1/swagger.yaml').as('getGHCClient');
    cy.intercept('**/ghc/v1/moves/search').as('getSearchResults');
    cy.intercept('**/ghc/v1/move/**').as('getMoves');
    cy.intercept('**/ghc/v1/orders/**').as('getOrders');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_shipments').as('getMTOShipments');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_service_items').as('getMTOServiceItems');

    const userId = '2419b1d6-097f-4dc4-8171-8f858967b4db';
    cy.apiSignInAsUser(userId, QAECSROfficeUserType);
  });

  it('is able to search by move code', () => {
    // Type move code into search bar (move code is default search type)
    cy.get('input[name="searchText"]').as('moveCodeSearchInput');
    cy.get('@moveCodeSearchInput').type(moveCode).blur();

    // Search for moves
    cy.get('button').contains('Search').click();
    cy.wait(['@getSearchResults']);

    // Verify results table contents
    cy.get('tbody > tr').as('results');
    cy.get('@results').should('have.length', 1);
    cy.get('@results').first().contains(moveCode);

    // Click result to navigate to move details page
    cy.get('@results').first().click();
    cy.url().should('include', `/moves/${moveCode}/details`);
    cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);
    cy.get('h1').contains('Move details');
  });

  it('is able to search by DOD ID', () => {
    const dodID = '1000000000';

    // Type dodID into search bar and select DOD ID as search type
    cy.get('[type="radio"]').check('dodID', { force: true });
    cy.get('input[name="searchText"]').type(dodID).blur();

    // Search for moves
    cy.get('button').contains('Search').click();
    cy.wait(['@getSearchResults']);

    // Verify results table contents
    cy.get('tbody > tr').as('results');
    cy.get('@results').should('have.length', 1);
    cy.get('@results').first().contains(moveCode);
    cy.get('@results').first().contains(dodID);

    // Click result to navigate to move details page
    cy.get('@results').first().click();
    cy.url().should('include', `/moves/${moveCode}/details`);
    cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);
    cy.get('h1').contains('Move details');
  });

  it('is able to search by customer name', () => {
    const name = 'QAECSRTestFirst';

    // Type dodID into search bar and select DOD ID as search type
    cy.get('[type="radio"]').check('customerName', { force: true });
    cy.get('input[name="searchText"]').type(name).blur();

    // Search for moves
    cy.get('button').contains('Search').click();
    cy.wait(['@getSearchResults']);

    // Verify results table contents
    cy.get('tbody > tr').as('results');
    cy.get('@results').should('have.length', 1); // Page size of 20 in table
    cy.get('@results').first().contains(name); // Should have QAECSRTestFirst in name

    // Click result to navigate to move details page
    cy.get('@results').first().click();
    cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);
    cy.get('h1').contains('Move details');
  });

  it('handles searches that do not return results', () => {
    // Search for a bad move code
    cy.get('input[name="searchText"]').type('BAD_ID').blur();
    cy.get('button').contains('Search').click();
    cy.wait(['@getSearchResults']);

    // Verify no results
    cy.get('[data-testid=table-queue] > h2').contains('Results (0)');
    cy.get('[data-testid=table-queue] > p').contains('No results found.');
  });
});
