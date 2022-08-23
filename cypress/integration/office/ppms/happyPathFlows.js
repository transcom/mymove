import { ServicesCounselorOfficeUserType } from '../../../support/constants';
import {
  fillOutDestinationInfo,
  fillOutIncentiveAndAdvance,
  fillOutOriginInfo,
  fillOutSitExpected,
  fillOutWeight,
  navigateToShipmentDetails,
} from '../../../support/ppmOfficeShared';

describe('Services counselor user', () => {
  before(() => {
    cy.prepareOfficeApp();
  });

  beforeEach(() => {
    cy.intercept('**/ghc/v1/swagger.yaml').as('getGHCClient');
    cy.intercept('**/ghc/v1/queues/counseling?page=1&perPage=20&sort=submittedAt&order=asc').as('getSortedMoves');

    cy.intercept('**/ghc/v1/queues/counseling?page=1&perPage=20&sort=submittedAt&order=asc&locator=**').as(
      'getFilterSortedMoves',
    );
    cy.intercept('**/ghc/v1/move/**').as('getMoves');
    cy.intercept('**/ghc/v1/orders/**').as('getOrders');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_shipments').as('getMTOShipments');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_shipments/**').as('updateMTOShipments');
    cy.intercept('**/ghc/v1/mto-shipments').as('createShipment');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_service_items').as('getMTOServiceItems');

    const userId = 'a6c8663f-998f-4626-a978-ad60da2476ec'; // services_counselor_role@office.mil
    cy.apiSignInAsUser(userId, ServicesCounselorOfficeUserType);
  });

  it('is able to edit a PPM shipment', () => {
    const moveLocator = 'PPMSCF';

    navigateToShipmentDetails(moveLocator);

    // View existing shipment
    cy.get('[data-testid="ShipmentContainer"] .usa-button').click();
    cy.wait(['@getMTOShipments', '@getMoves', '@getOrders']);

    fillOutSitExpected();

    // Submit page 1 of form
    cy.get('[data-testid="submitForm"]').should('be.enabled').click();
    cy.wait(['@updateMTOShipments', '@getMTOShipments', '@getMoves', '@getOrders', '@getMTOServiceItems']);

    // Verify SIT info
    cy.contains('Government constructed cost: $326');
    cy.contains('1,000 lbs of destination SIT at 30813 for 31 days.');

    // Update page 2
    fillOutIncentiveAndAdvance();
    cy.get('[data-testid="errorMessage"]').contains('Required');
    cy.get('[data-testid="counselor-remarks"]').clear().type('Increased incentive to max').blur();

    // Submit page 2 of form
    cy.get('[data-testid="submitForm"]').should('be.enabled').click();
    cy.wait(['@updateMTOShipments', '@getMTOShipments', '@getMoves', '@getOrders', '@getMTOServiceItems']);

    // Expand details and verify information
    cy.contains('Your changes were saved.');
    cy.get('[data-prefix="fas"][data-icon="chevron-down"]').click();
    cy.get('[data-testid="expectedDepartureDate"]').contains('15 Mar 2020');
    cy.get('[data-testid="originZIP"]').contains('90210');
    cy.get('[data-testid="secondOriginZIP"]').contains('90211');
    cy.get('[data-testid="destinationZIP"]').contains('30813');
    cy.get('[data-testid="secondDestinationZIP"]').contains('30814');
    cy.get('[data-testid="sitPlanned"]').contains('yes');
    cy.get('[data-testid="estimatedWeight"]').contains('4,000 lbs');
    cy.get('[data-testid="proGearWeight"]').contains('Yes, 1,987 lbs');
    cy.get('[data-testid="spouseProGear"]').contains('Yes, 498 lbs');
    cy.get('[data-testid="estimatedIncentive"]').contains('$10,000');
    cy.get('[data-testid="hasRequestedAdvance"]').contains('Yes, $6,000');
    cy.get('[data-testid="counselorRemarks"]').contains('Increased incentive to max');
  });

  it('is able to add a second PPM shipment', () => {
    const moveLocator = 'PPMSCF';

    navigateToShipmentDetails(moveLocator);

    cy.get('[data-testid="dropdown"]').select('PPM');
    cy.wait(['@getMTOShipments', '@getMoves', '@getOrders']);

    // Fill out page one
    fillOutOriginInfo();
    fillOutDestinationInfo();
    fillOutSitExpected();
    fillOutWeight({ hasProGear: true });

    cy.get('[data-testid="submitForm"]').should('be.enabled').click();
    cy.wait(['@createShipment', '@getMTOShipments', '@getMoves', '@getOrders']);

    // Verify SIT info
    cy.contains('Government constructed cost: $379');
    cy.contains('1,000 lbs of destination SIT at 76127 for 31 days.');

    // Fill out page two
    fillOutIncentiveAndAdvance({ advance: '10000' });
    cy.get('[data-testid="errorMessage"]').contains('Required');
    cy.get('[data-testid="counselor-remarks"]').clear().type('Added correct incentive').blur();

    // Submit page two
    cy.get('[data-testid="submitForm"]').should('be.enabled').click();
    cy.wait(['@updateMTOShipments', '@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);

    // Expand details and verify information
    cy.contains('Your changes were saved.');
    cy.get('[data-prefix="fas"][data-icon="chevron-down"]').last().click();
    cy.get('[data-testid="expectedDepartureDate"]').contains('09 Jun 2022');
    cy.get('[data-testid="originZIP"]').contains('90210');
    cy.get('[data-testid="secondOriginZIP"]').contains('07003');
    cy.get('[data-testid="destinationZIP"]').contains('76127');
    cy.get('[data-testid="secondDestinationZIP"]').contains('08540');
    cy.get('[data-testid="sitPlanned"]').contains('yes');
    cy.get('[data-testid="estimatedWeight"]').contains('4,000 lbs');
    cy.get('[data-testid="proGearWeight"]').contains('Yes, 1,000 lbs');
    cy.get('[data-testid="spouseProGear"]').contains('Yes, 500 lbs');
    cy.get('[data-testid="estimatedIncentive"]').contains('$201,506');
    // Need to add back when bug is merged in
    // cy.get('[data-testid="hasRequestedAdvance"]').contains('Yes, $10,000');
    cy.get('[data-testid="counselorRemarks"]').contains('Added correct incentive');
  });
});
