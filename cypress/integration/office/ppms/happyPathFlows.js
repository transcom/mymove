import { ServicesCounselorOfficeUserType } from '../../../support/constants';
import { navigateToShipmentDetails } from '../../../support/ppmOfficeShared';

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
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_service_items').as('getMTOServiceItems');

    const userId = 'a6c8663f-998f-4626-a978-ad60da2476ec'; // services_counselor_role@office.mil
    cy.apiSignInAsUser(userId, ServicesCounselorOfficeUserType);
  });

  it('is able to edit a PPM shipment', () => {
    const moveLocator = 'PPMSCF';

    navigateToShipmentDetails(moveLocator);

    // view existing shipment
    cy.get('[data-testid="ShipmentContainer"] .usa-button').click();
    cy.wait(['@getMTOShipments', '@getMoves', '@getOrders']);

    // add SIT
    cy.get('input[name="sitExpected"][value="yes"]').check({ force: true });
    cy.get('input[name="sitEstimatedWeight"]').clear().type(1000).blur();
    cy.get('input[name="sitEstimatedEntryDate"]').clear().type('01 Mar 2020').blur();
    cy.get('input[name="sitEstimatedDepartureDate"]').clear().type('31 Mar 2020').blur();

    // submit page 1 of form
    cy.get('[data-testid="submitForm"]').should('be.enabled').click();
    cy.wait(['@updateMTOShipments', '@getMTOShipments', '@getMoves', '@getOrders', '@getMTOServiceItems']);

    // update page 2
    cy.get('input[name="advance"]').clear().type('6000').blur();
    cy.get('[data-testid="errorMessage"]').contains('Required');
    cy.get('[data-testid="counselor-remarks"]').clear().type('Increased incentive to max').blur();

    // submit page 2 of form
    cy.get('[data-testid="submitForm"]').should('be.enabled').click();
    cy.wait(['@updateMTOShipments', '@getMTOShipments', '@getMoves', '@getOrders', '@getMTOServiceItems']);

    // expand details and verify information
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
    cy.get('input[name="expectedDepartureDate"]').clear().type('09 Jun 2022').blur();
    cy.get('input[name="pickupPostalCode"]').clear().type('90210').blur();
    cy.get('input[name="secondPickupPostalCode"]').clear().type('07003').blur();

    cy.get('input[name="destinationPostalCode"]').clear().type('76127').blur();
    cy.get('input[name="secondDestinationPostalCode"]').clear().type('08540').blur();

    cy.get('input[name="estimatedWeight"]').clear().type('4000').blur();
    cy.get('input[name="hasProGear"][value="yes"]').check({ force: true });
    cy.get('input[name="proGearWeight"]').type(1000).blur();
    cy.get('input[name="spouseProGearWeight"]').type(500).blur();

    cy.get('input[name="sitExpected"][value="yes"]').check({ force: true });
    cy.get('input[name="sitEstimatedWeight"]').clear().type(1000).blur();
    cy.get('input[name="sitEstimatedEntryDate"]').clear().type('01 Mar 2020').blur();
    cy.get('input[name="sitEstimatedDepartureDate"]').clear().type('31 Mar 2020').blur();

    cy.get('[data-testid="submitForm"]').should('be.enabled').click();
    cy.wait(['@getMTOShipments', '@getMoves', '@getOrders']);

    // Fill out page two
    cy.get('input[name="advanceRequested"][value="Yes"]').check({ force: true });
    cy.get('input[name="advance"]').clear().type('10000').blur();
    cy.get('[data-testid="errorMessage"]').contains('Required');
    cy.get('[data-testid="counselor-remarks"]').clear().type('Increased incentive to max').blur();

    cy.get('[data-testid="submitForm"]').should('be.enabled').click();
    cy.wait(['@updateMTOShipments', '@getMTOShipments', '@getMoves', '@getOrders', '@getMTOServiceItems']);
  });
});
