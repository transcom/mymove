import { ServicesCounselorOfficeUserType } from '../../../support/constants';
import {
  fillOutDestinationInfo,
  fillOutIncentiveAndAdvance,
  fillOutOriginInfo,
  fillOutWeight,
  navigateToShipmentDetails,
} from '../../../support/ppmOfficeShared';

describe('Services counselor user', () => {
  before(() => {
    cy.prepareOfficeApp();
  });

  beforeEach(() => {
    cy.intercept('**/ghc/v1/swagger.yaml').as('getGHCClient');
    cy.intercept('**/ghc/v1/queues/counseling?page=1&perPage=20&sort=submittedAt&order=asc&needsPPMCloseout=false').as(
      'getSortedMoves',
    );

    // Note this intercept is specific to a particular move locator
    cy.intercept(
      '**/ghc/v1/queues/counseling?page=1&perPage=20&sort=closeoutInitiated&order=asc&locator=PPMADD&needsPPMCloseout=true',
    ).as('getFilterSortedMoves');
    cy.intercept('GET', '**/ghc/v1/move/**').as('getMoves');
    cy.intercept('GET', '**/ghc/v1/orders/**').as('getOrders');
    cy.intercept('GET', '**/ghc/v1/move_task_orders/**/mto_shipments').as('getMTOShipments');
    cy.intercept('GET', '**/ghc/v1/move_task_orders/**/mto_service_items').as('getMTOServiceItems');
    cy.intercept('PATCH', '**/ghc/v1/move_task_orders/**/mto_shipments/**').as('updateMTOShipments');

    cy.intercept('POST', '**/ghc/v1/mto-shipments').as('createShipment');
    cy.intercept('DELETE', '**/ghc/v1/shipments/**').as('deleteShipment');

    const userId = 'a6c8663f-998f-4626-a978-ad60da2476ec';
    cy.apiSignInAsUser(userId, ServicesCounselorOfficeUserType);
  });

  it('is able to add a new PPM shipment', () => {
    const moveLocator = 'PPMSC1';

    /**
     * Move Details page
     */
    navigateToShipmentDetails(moveLocator);

    // Delete existing shipment
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
    fillOutOriginInfo();
    fillOutDestinationInfo();
    fillOutWeight({ hasProGear: true });

    cy.get('[data-testid="submitForm"]').should('be.enabled').click();
    cy.wait(['@createShipment', '@getMTOShipments', '@getMoves', '@getOrders']);

    // Fill out page two
    cy.contains('Incentive & advance');
    fillOutIncentiveAndAdvance({ advance: '5987' });
    cy.get('[data-testid="counselor-remarks"]').clear().type('The requested advance amount has been added.').blur();

    cy.get('[data-testid="submitForm"]').should('be.enabled').click();
    cy.wait(['@updateMTOShipments', '@getMTOShipments', '@getMoves', '@getOrders', '@getMTOServiceItems']);

    // Confirm new shipment is visible
    cy.get('[data-testid="ShipmentContainer"]').within(($shipmentContainer) => {
      // Verify unexpanded view
      cy.get('[data-testid="expectedDepartureDate"]').contains('09 Jun 2022');
      cy.get('[data-testid="originZIP"]').contains('90210');
      cy.get('[data-testid="destinationZIP"]').contains('76127');
      cy.get('[data-testid="sitPlanned"]').contains('No');
      cy.get('[data-testid="estimatedWeight"]').contains('4,000 lbs');
      cy.get('[data-testid="hasRequestedAdvance"]').contains('Yes, $5,987');
      cy.get('[data-testid="counselorRemarks"]').contains('The requested advance amount has been added.');

      // Verify fields in the expanded view
      cy.get('[data-prefix="fas"][data-icon="chevron-down"]').click();
      cy.get('[data-testid="secondOriginZIP"]').contains('07003');
      cy.get('[data-testid="secondDestinationZIP"]').contains('08540');
      cy.get('[data-testid="proGearWeight"]').contains('Yes, 1,000 lbs');
      cy.get('[data-testid="spouseProGear"]').contains('Yes, 500 lbs');
      cy.get('[data-testid="estimatedIncentive"]').contains('$67,692');
    });
  });
});
