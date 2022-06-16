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
    cy.intercept('**/ghc/v1/queues/counseling?page=1&perPage=20&sort=submittedAt&order=asc&locator=PPMADD').as(
      'getFilterSortedMoves',
    );
    cy.intercept('**/ghc/v1/move/**').as('getMoves');
    cy.intercept('**/ghc/v1/orders/**').as('getOrders');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_shipments').as('getMTOShipments');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_service_items').as('getMTOServiceItems');

    cy.intercept('POST', '**/ghc/v1/mto-shipments').as('createShipment');
    cy.intercept('DELETE', '**/ghc/v1/shipments/**').as('deleteShipment');

    const userId = 'a6c8663f-998f-4626-a978-ad60da2476ec';
    cy.apiSignInAsUser(userId, ServicesCounselorOfficeUserType);
  });

  it('is able to add a new PPM shipment', () => {
    const moveLocator = 'PPMADD';

    /**
     * Move Details page
     */
    navigateToShipmentDetails(moveLocator);

    // Delete exisitng shipment
    cy.get('[data-testid="expectedDepartureDate"]').contains('15 Mar 2020');
    cy.get('[data-testid="ShipmentContainer"] .usa-button').click();
    cy.wait(['@getMTOShipments', '@getMoves']);
    cy.get('[data-testid="grid"] button').contains('Delete shipment').click();
    cy.get('[data-testid="modal"]').should('be.visible');

    cy.get('[data-testid="modal"] button').contains('Delete shipment').click({ force: true });
    cy.wait(['@deleteShipment', '@getMoves']);

    cy.get('[data-testid="ShipmentContainer"] .usa-button').should('have.length', 0);

    // Click add shipment button and select PPM
    cy.get('[data-testid="dropdown"]').first().select('PPM');

    // Fill out page one
    cy.get('input[name="pickupPostalCode"]').clear().type('90210').blur();
    cy.get('input[name="secondPickupPostalCode"]').clear().type('07003').blur();

    cy.get('input[name="destinationPostalCode"]').clear().type('76127');
    cy.get('input[name="secondDestinationPostalCode"]').clear().type('08540');

    cy.get('input[name="expectedDepartureDate"]').clear().type('09 Jun 2022').blur();

    cy.get('input[name="estimatedWeight"]').clear().type('4000').blur();
    cy.get('input[name="hasProGear"][value="yes"]').check({ force: true });
    cy.get('input[name="proGearWeight"]').type(1000).blur();
    cy.get('input[name="spouseProGearWeight"]').type(500).blur();

    cy.get('button[data-testid="submitForm"]').click();

    // Fill out page two
    cy.get('input[name="advance"]').clear().type('2100').blur();
    cy.get('textarea[data-testid="counselor-remarks"]')
      .clear()
      .type('The requested advance amount has been added.')
      .blur();

    cy.get('button[data-testid="submitForm"]').click();

    // TODO User should be automatically redirected to the Move Details page. This needs to be updated when that work is completed
    cy.get('a[data-testid="MoveDetails-Tab"]').click();
    cy.url().should('include', `/counseling/moves/${moveLocator}/details`);
    cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);

    // Confirm new shipment is visible

    cy.get('[data-testid="expectedDepartureDate"]').contains('09 Jun 2022');
    cy.get('[data-testid="originZIP"]').contains('90210');
    cy.get('[data-testid="destinationZIP"]').contains('76127');
    cy.get('[data-testid="sitPlanned"]').contains('no');
    cy.get('[data-testid="estimatedWeight"]').contains('4,000 lbs');
    // TODO uncomment these assertions after display bug has been fixed
    // cy.get('[data-testid="hasRequestedAdvance"]').contains('Yes, $5,987');
    // cy.get('[data-testid="counselorRemarks"]').contains('The requested advance amount has been added.');
  });
});
