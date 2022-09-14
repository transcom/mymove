export const searchForAndNavigateToMove = (moveLocator) => {
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
  cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);
};
