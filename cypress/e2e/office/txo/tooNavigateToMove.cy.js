export function navigateToMove(moveLocator) {
  // TOO Moves queue
  cy.wait(['@getSortedOrders']);
  cy.get('input[name="locator"]').as('moveCodeFilterInput');

  // type in move code/locator to filter
  cy.get('@moveCodeFilterInput').type(moveLocator).blur();
  cy.get('tbody > tr').as('results');

  // click result to navigate to move details page
  cy.get('@results').first().click();
  cy.url().should('include', `/moves/${moveLocator}/details`);

  // Move Details page
  cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);
}
