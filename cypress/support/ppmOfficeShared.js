export function navigateToShipmentDetails(moveLocator) {
  /**
   * SC Moves queue
   */
  cy.wait(['@getSortedMoves']);
  cy.get('input[name="locator"]').as('moveCodeFilterInput');

  // type in move code/locator to filter
  cy.get('@moveCodeFilterInput').type(moveLocator).blur();
  cy.wait(['@getFilterSortedMoves']);

  // check if results appear, should be 1
  // and see if result have move code
  cy.get('tbody > tr').as('results');
  cy.get('@results').should('have.length', 1);
  cy.get('@results').first().contains(moveLocator);

  // click result to navigate to move details page
  cy.get('@results').first().click();
  cy.url().should('include', `/counseling/moves/${moveLocator}/details`);
  cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);
}
