import { ServicesCounselorOfficeUserType } from '../../../support/constants';
import { navigateToShipmentDetails } from '../../../support/ppmOfficeShared';

describe('Services counselor user', () => {
  before(() => {
    cy.prepareOfficeApp();
  });

  beforeEach(() => {
    cy.intercept('**/ghc/v1/swagger.yaml').as('getGHCClient');
    cy.intercept('**/ghc/v1/queues/counseling?page=1&perPage=20&sort=submittedAt&order=asc').as('getSortedMoves');

    // Note this intercept is specific to a particular move locator
    cy.intercept('**/ghc/v1/queues/counseling?page=1&perPage=20&sort=submittedAt&order=asc&locator=PPMSC1').as(
      'getFilterSortedMoves',
    );
    cy.intercept('**/ghc/v1/move/**').as('getMoves');
    cy.intercept('**/ghc/v1/orders/**').as('getOrders');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_shipments').as('getMTOShipments');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_service_items').as('getMTOServiceItems');

    const userId = 'a6c8663f-998f-4626-a978-ad60da2476ec';
    cy.apiSignInAsUser(userId, ServicesCounselorOfficeUserType);
  });

  it('is able to click on move and submit after using the move code filter', () => {
    const moveLocator = 'PPMSC1';

    /**
     * Move Details page
     */
    navigateToShipmentDetails(moveLocator);

    // Shipment card
    cy.get('[data-testid="ShipmentContainer"] ').trigger('click');

    cy.get('[data-testid="expectedDepartureDate"]').contains('15 Mar 2020');
    cy.get('[data-testid="originZIP"]').contains('90210');
    cy.get('[data-testid="destinationZIP"]').contains('30813');
    cy.get('[data-testid="sitPlanned"]').contains('no');
    cy.get('[data-testid="estimatedWeight"]').contains('4,000 lbs');
    cy.get('[data-testid="hasRequestedAdvance"]').contains('Yes, $5,987');
    cy.get('[data-testid="counselorRemarks"]').contains('—');
  });
});
