import { ServicesCounselorOfficeUserType } from '../../../support/constants';

describe('Services counselor user', () => {
  before(() => {
    cy.prepareOfficeApp();
  });

  beforeEach(() => {
    cy.intercept('**/ghc/v1/swagger.yaml').as('getGHCClient');
    cy.intercept('**/ghc/v1/queues/counseling?page=1&perPage=20&sort=submittedAt&order=asc').as('getSortedMoves');
    cy.intercept('**/ghc/v1/queues/counseling?page=1&perPage=20&sort=submittedAt&order=asc&locator=SCE4ET').as(
      'getFilterSortedMoves',
    );
    cy.intercept('**/ghc/v1/move/**').as('getMoves');
    cy.intercept('**/ghc/v1/orders/**').as('getOrders');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_shipments').as('getMTOShipments');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_service_items').as('getMTOServiceItems');
    cy.intercept('**/ghc/v1/move-task-orders/**/status/service-counseling-completed').as(
      'patchServiceCounselingCompleted',
    );

    const userId = 'a6c8663f-998f-4626-a978-ad60da2476ec';
    cy.apiSignInAsUser(userId, ServicesCounselorOfficeUserType);
  });

  it('is able to click on move and submit after using the move code filter', () => {
    const moveLocator = 'SCE4ET';

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

    /**
     * Move Details page
     */
    // click to trigger confirmation modal
    cy.contains('Submit move details').click();

    // modal should pop up with text
    cy.get('h2').contains('Are you sure?');
    cy.get('p').contains('You can’t make changes after you submit the move.');

    // click submit
    cy.get('button').contains('Yes, submit').click();
    cy.waitFor(['@patchServiceCounselingCompleted', '@getMoves']);

    // verify success alert
    cy.contains('Move submitted.');
  });
});
