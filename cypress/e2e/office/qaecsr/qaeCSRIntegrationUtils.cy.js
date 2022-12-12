export const searchForAndNavigateToMove = (moveCode = 'QAEHLP') => {
  cy.get('input[name="searchText"]').as('moveCodeSearchInput');
  cy.get('@moveCodeSearchInput').type(moveCode).blur();
  cy.get('button').contains('Search').click();
  cy.wait(['@getSearchResults']);

  cy.get('tbody > tr').as('results');
  cy.get('@results').should('have.length', 1);
  cy.get('@results').first().contains(moveCode);

  // click result to navigate to move details page
  cy.get('@results').first().click();
  cy.url().should('include', `/moves/${moveCode}/details`);
  cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);
};
