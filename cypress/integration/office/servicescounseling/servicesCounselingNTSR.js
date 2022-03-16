import { ServicesCounselorOfficeUserType } from '../../../support/constants';

describe('Services counselor user', () => {
  before(() => {
    cy.prepareOfficeApp();
  });

  beforeEach(() => {
    cy.intercept('**/ghc/v1/swagger.yaml').as('getGHCClient');
    cy.intercept('**/ghc/v1/queues/counseling?page=1&perPage=20&sort=submittedAt&order=asc').as('getSortedMoves');
    cy.intercept('**/ghc/v1/move/**').as('getMoves');
    cy.intercept('**/ghc/v1/orders/**').as('getOrders');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_shipments').as('getMTOShipments');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_service_items').as('getMTOServiceItems');
    cy.intercept('**/ghc/v1/move-task-orders/**/status/service-counseling-completed').as(
      'patchServiceCounselingCompleted',
    );
    cy.intercept('POST', '**/ghc/v1/mto-shipments').as('createShipment');
    cy.intercept('PATCH', '**/ghc/v1/move_task_orders/**/mto_shipments/**').as('patchShipment');

    const userId = 'a6c8663f-998f-4626-a978-ad60da2476ec';
    cy.apiSignInAsUser(userId, ServicesCounselorOfficeUserType);
  });

  it('Services Counselor can add an NTS-release shipment to the customer move', () => {
    navigateToMove('NTSRHG');
    addNTSReleaseShipment();
  });

  it('Services Counselor can delete/remove an NTS-release shipment request', () => {
    navigateToMove('NTSRHG');
    addNTSReleaseShipment();

    cy.get('[data-testid="ShipmentContainer"] .usa-button').should('have.length', 3);

    cy.get('[data-testid="ShipmentContainer"] .usa-button').last().click();
    cy.waitFor(['@getMTOServiceItems']);
    // click to trigger confirmation modal
    cy.get('[data-testid="grid"] button').contains('Delete shipment').click();

    cy.get('[data-testid="modal"]').should('be.visible');

    cy.get('[data-testid="modal"] button').contains('Delete shipment').click();
    cy.waitFor(['@patchServiceCounselingCompleted', '@getMoves']);

    cy.get('[data-testid="ShipmentContainer"] .usa-button').should('have.length', 2);
  });

  it('Services Counselor can enter accounting codes on the Orders Page', () => {
    navigateToMove('NTSRHG');

    cy.get('[data-testid="ShipmentContainer"] .usa-button').last().click();
    cy.waitFor(['@getMTOServiceItems']);

    cy.get('[data-testid="grid"] button').contains('Add or edit codes').click();

    cy.get('form').within(($form) => {
      cy.get('input[name="tac"]').click().clear().type('E15A');
      cy.get('input[name="sac"]').click().clear().type('4K988AS098F');
      cy.get('input[name="ntsTac"]').click().clear().type('F123');
      cy.get('input[name="ntsSac"]').click().clear().type('3L988AS098F');
      // Edit orders page | Save
      cy.get('button').contains('Save').click();
    });

    cy.get('[data-testid="tacMDC"]').contains('E15A');
    cy.get('[data-testid="sacSDN"]').contains('4K988AS098F');
    cy.get('[data-testid="NTStac"]').contains('F123');
    cy.get('[data-testid="NTSsac"]').contains('3L988AS098F');
  });

  it('Services Counselor can assign accounting code(s) to a shipment', () => {
    navigateToMove('NTSRHG');

    cy.get('[data-testid="ShipmentContainer"] .usa-button').last().click();
    cy.waitFor(['@getMTOServiceItems, @getMoves']);

    cy.get('[data-testid="radio"] [for="tacType-NTS"]').click();
    cy.get('[data-testid="radio"] [for="sacType-HHG"]').click();

    cy.get('[data-testid="submitForm"]').click();

    cy.wait('@patchShipment');
    cy.get('.usa-alert__text').contains('Your changes were saved.');

    cy.get('[data-testid="ShipmentContainer"]')
      .last()
      .within(($form) => {
        cy.get('[data-icon="chevron-down"]').click();
        cy.get('[data-testid="tacType"]').contains('F123 (NTS)');
        cy.get('[data-testid="sacType"]').contains('4K988AS098F (HHG)');
      });
  });

  it('Services Counselor can submit a move with an NTS-release shipment', () => {
    navigateToMove('NTSRHG');
    // click to trigger confirmation modal
    cy.contains('Submit move details').click();

    cy.get('[data-testid="modal"]').should('be.visible');

    cy.get('button').contains('Yes, submit').click();
    cy.waitFor(['@patchServiceCounselingCompleted', '@getMoves']);

    // verify success alert
    cy.contains('Move submitted.');
  });

  it('Services Counselor can see errors/warnings for missing data, then make edits', () => {
    navigateToMove('NTSRMN');

    cy.get('[data-testid="ShipmentContainer"]')
      .last()
      .within(($form) => {
        cy.get('[data-icon="chevron-down"]').click();

        cy.get('div[class*="warning"] [data-testid="ntsRecordedWeight"]').should('exist');
        cy.get('div[class*="missingInfoError"] [data-testid="storageFacilityName"]').should('exist');
        cy.get('[data-testid="storageFacilityName"]').contains('Missing');
        cy.get('div[class*="warning"] [data-testid="serviceOrderNumber"]').should('exist');
        cy.get('div[class*="missingInfoError"] [data-testid="storageFacilityAddress"]').should('exist');
        cy.get('[data-testid="storageFacilityAddress"]').contains('Missing');
        cy.get('div[class*="warning"] [data-testid="counselorRemarks"]').should('exist');
        cy.get('div[class*="warning"] [data-testid="tacType"]').should('exist');
        cy.get('div[class*="warning"] [data-testid="sacType"]').should('exist');
      });

    editNTSReleaseShipment();
  });
});

function navigateToMove(moveLocator) {
  //SC Moves queue
  cy.wait(['@getSortedMoves']);
  cy.get('input[name="locator"]').as('moveCodeFilterInput');

  // type in move code/locator to filter
  cy.get('@moveCodeFilterInput').type(moveLocator).blur();

  cy.get('tbody > tr').as('results');

  // click result to navigate to move details page
  cy.get('@results').first().click();
  cy.url().should('include', `/counseling/moves/${moveLocator}/details`);
  cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);
}

function addNTSReleaseShipment() {
  cy.get('[data-testid="dropdown"]').first().select('NTS-release');

  // Previously recorded weight
  cy.get('#ntsRecordedWeight').clear().type('1300').blur();

  // Storage facility info
  cy.get('#facilityName').type('Sample Facility Name');
  cy.get('#facilityPhone').type('999-999-9999');
  cy.get('#facilityEmail').type('sample@example.com');
  cy.get('#facilityServiceOrderNumber').type('999999');

  // Storage facility address
  cy.get('input[name="storageFacility.address.streetAddress1"]').type('148 S East St');
  cy.get('input[name="storageFacility.address.streetAddress2"]').type('Suite 7A');
  cy.get('input[name="storageFacility.address.city"]').type('Sample City');
  cy.get('select[name="storageFacility.address.state"]').select('GA');
  cy.get('input[name="storageFacility.address.postalCode"]').type('30301');
  cy.get('#facilityLotNumber').type('1111111');

  // Requested delivery date
  cy.get('#requestedDeliveryDate').clear().type('20 Mar 2022').blur();

  // Delivery location
  cy.get('input[name="delivery.address.streetAddress1"]').type('448 Washington Blvd NE');
  cy.get('input[name="delivery.address.streetAddress2"]').type('Apt D3');
  cy.get('input[name="delivery.address.city"]').type('Another City');
  cy.get('select[name="delivery.address.state"]').select('AL');
  cy.get('input[name="delivery.address.postalCode"]').type('36101');
  cy.get('#destinationType').select('HOME_OF_RECORD');

  // Receiving agent
  cy.get('input[name="delivery.agent.firstName"]').type('Skyler');
  cy.get('input[name="delivery.agent.lastName"]').type('Hunt');
  cy.get('input[name="delivery.agent.phone"]').type('999-999-9999');
  cy.get('input[name="delivery.agent.email"]').type('skyler.hunt@example.com');

  // Remarks
  cy.get(`[data-testid="remarks"]`).should('not.exist');
  cy.get('[data-testid="counselor-remarks"]').type('NTS-release counselor remarks');

  cy.get('[data-testid="submitForm"]').click();
  // the shipment should be saved with the type
  cy.wait('@createShipment');

  // the new shipment is visible on the Move details page
  cy.get('[data-testid="ShipmentContainer"]')
    .last()
    .within(($form) => {
      cy.get('[data-icon="chevron-down"]').click();

      cy.get('[data-testid="ntsRecordedWeight"]').contains('1,300');
      cy.get('[data-testid="storageFacilityName"]').contains('Sample Facility Name');
      cy.get('[data-testid="serviceOrderNumber"]').contains('999999');
      cy.get('[data-testid="storageFacilityAddress"]').contains('148 S East St, Sample City, GA 30301');
      cy.get('[data-testid="storageFacilityAddress"]').contains('1111111');
      cy.get('[data-testid="requestedDeliveryDate"]').contains('20 Mar 2022');
      cy.get('[data-testid="destinationAddress"]').contains('448 Washington Blvd NE, Another City, AL 36101');
      cy.get('[data-testid="secondaryDeliveryAddress"]').contains('—');
      cy.get('[data-testid="customerRemarks"]').contains('—');
      cy.get('[data-testid="counselorRemarks"]').contains('NTS-release counselor remarks');
      cy.get('[data-testid="tacType"]').contains('—');
      cy.get('[data-testid="sacType"]').contains('—');
    });
}

function editNTSReleaseShipment() {
  cy.get('[data-testid="dropdown"]').first().select('NTS-release');

  // Previously recorded weight
  cy.get('#ntsRecordedWeight').clear().type('1100').blur();

  // Storage facility info
  cy.get('#facilityName').clear().type('AAA Facility Name');
  cy.get('#facilityPhone').clear().type('999-999-9999');
  cy.get('#facilityEmail').clear().type('aaa@example.com');
  cy.get('#facilityServiceOrderNumber').clear().type('123456');

  // Storage facility address
  cy.get('input[name="storageFacility.address.streetAddress1"]').clear().type('9 W 2nd Ave');
  cy.get('input[name="storageFacility.address.streetAddress2"]').clear().type('Bldg 3');
  cy.get('input[name="storageFacility.address.city"]').clear().type('Big City');
  cy.get('select[name="storageFacility.address.state"]').select('SC');
  cy.get('input[name="storageFacility.address.postalCode"]').clear().type('29201');
  cy.get('#facilityLotNumber').clear().type('2222222');

  // Requested delivery date
  cy.get('#requestedDeliveryDate').clear().type('21 Mar 2022').blur();

  // Delivery location
  cy.get('input[name="delivery.address.streetAddress1"]').clear().type('4124 Apache Dr');
  cy.get('input[name="delivery.address.streetAddress2"]').clear().type('Apt 18C');
  cy.get('input[name="delivery.address.city"]').clear().type('Little City');
  cy.get('select[name="delivery.address.state"]').select('GA');
  cy.get('input[name="delivery.address.postalCode"]').clear().type('30901');
  cy.get('#destinationType').select('HOME_OF_RECORD');

  // Receiving agent
  cy.get('input[name="delivery.agent.firstName"]').clear().type('Jody');
  cy.get('input[name="delivery.agent.lastName"]').clear().type('Pitkin');
  cy.get('input[name="delivery.agent.phone"]').clear().type('999-111-1111');
  cy.get('input[name="delivery.agent.email"]').clear().type('jody.pitkin@example.com');

  // Remarks
  cy.get(`[data-testid="remarks"]`).should('not.exist');
  cy.get('[data-testid="counselor-remarks"]').clear().type('NTS-release edited counselor remarks');

  cy.get('[data-testid="submitForm"]').click();
  // the shipment should be saved with the type
  cy.wait('@createShipment');

  // the new shipment is visible on the Move details page
  cy.get('[data-testid="ShipmentContainer"]')
    .last()
    .within(($form) => {
      cy.get('[data-icon="chevron-down"]').click();

      cy.get('[data-testid="ntsRecordedWeight"]').contains('1,100');
      cy.get('[data-testid="storageFacilityName"]').contains('AAA Facility Name');
      cy.get('[data-testid="serviceOrderNumber"]').contains('123456');
      cy.get('[data-testid="storageFacilityAddress"]').contains('9 W 2nd Ave, Big City, SC 29201');
      cy.get('[data-testid="storageFacilityAddress"]').contains('2222222');
      cy.get('[data-testid="requestedDeliveryDate"]').contains('21 Mar 2022');
      cy.get('[data-testid="destinationAddress"]').contains('4124 Apache Dr, Little City, GA 30901');
      cy.get('[data-testid="secondaryDeliveryAddress"]').contains('—');
      cy.get('[data-testid="customerRemarks"]').contains('—');
      cy.get('[data-testid="counselorRemarks"]').contains('NTS-release edited counselor remarks');
      cy.get('[data-testid="tacType"]').contains('—');
      cy.get('[data-testid="sacType"]').contains('—');
    });
}
