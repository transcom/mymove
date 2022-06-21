import { TOOOfficeUserType } from '../../../support/constants';

describe('TOO user', () => {
  before(() => {
    cy.prepareOfficeApp();
  });

  beforeEach(() => {
    cy.intercept('**/ghc/v1/swagger.yaml').as('getGHCClient');
    cy.intercept('**/ghc/v1/queues/moves?page=1&perPage=20&sort=status&order=asc').as('getSortedOrders');
    cy.intercept('**/ghc/v1/move/**').as('getMoves');
    cy.intercept('GET', '**/ghc/v1/orders/**').as('getOrders');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_shipments').as('getMTOShipments');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_service_items').as('getMTOServiceItems');
    cy.intercept('POST', '**/ghc/v1/shipments/**/approve').as('approveShipment');
    cy.intercept('POST', '**/ghc/v1/shipments/**/request-cancellation').as('requestShipmentCancellation');
    cy.intercept('PATCH', '**/ghc/v1/move-task-orders/**/status').as('patchMTOStatus');
    cy.intercept('PATCH', '**/ghc/v1/move-task-orders/**/service-items/**/status').as('patchMTOServiceItems');
    cy.intercept('PATCH', '**/ghc/v1/move_task_orders/**/mto_shipments/**').as('patchShipment');
    cy.intercept('PATCH', '**/ghc/v1/orders/**/allowances').as('patchAllowances');
    cy.intercept('**/ghc/v1/moves/**/financial-review-flag').as('financialReviewFlagCompleted');

    // This user has multiple roles, which is the kind of user we use to test in staging.
    // By using this type of user, we can catch bugs like the one fixed in PR 6706.
    const userId = 'dcf86235-53d3-43dd-8ee8-54212ae3078f';
    cy.apiSignInAsUser(userId, TOOOfficeUserType);
  });

  // TODO FOR NTS-RELEASE
  // Enter/edit shipment weight if needed -- user will look up in TOPS and manually enter
  // TOO has to enter Service Order Number (SON) for NTS-RELEASE shipments prior to posting to the MTO

  // This test covers editing the NTS shipment and prepares it for approval
  it('TOO can edit a request for Domestic NTS Shipment handled by the Prime', () => {
    navigateToMove('PRINTR');

    cy.contains('Approve selected').should('be.disabled');
    cy.get('[data-testid="ShipmentContainer"]')
      .last()
      .within(() => {
        cy.get('[data-testid="shipment-display-checkbox"]').should('be.disabled');
        cy.get('[data-icon="chevron-down"]').click();

        cy.get('div[class*="missingInfoError"] [data-testid="ntsRecordedWeight"]').should('exist');
        cy.get('div[class*="missingInfoError"] [data-testid="storageFacilityName"]').should('exist');
        cy.get('[data-testid="storageFacilityName"]').contains('Missing');
        cy.get('div[class*="missingInfoError"] [data-testid="serviceOrderNumber"]').should('exist');
        cy.get('div[class*="missingInfoError"] [data-testid="storageFacilityAddress"]').should('exist');
        cy.get('[data-testid="storageFacilityAddress"]').contains('Missing');
        cy.get('div[class*="missingInfoError"] [data-testid="tacType"]').should('exist');
      });

    editTacSac();

    // Edit shipments to enter missing info
    cy.get('[data-testid="ShipmentContainer"] .usa-button').last().click();
    cy.wait(['@getMTOShipments', '@getMoves']);

    // Basic info
    cy.get('#ntsRecordedWeight').clear().type('3000').blur();

    // Storage facility info
    cy.get('#facilityName').clear().type('Sample Facility Name').blur();
    cy.get('#facilityPhone').clear().type('999-999-9999').blur();
    cy.get('#facilityEmail').clear().type('sample@example.com').blur();
    cy.get('#facilityServiceOrderNumber').clear().type('999999').blur();

    // Storage facility address
    cy.get('input[name="storageFacility.address.streetAddress1"]').clear().type('148 S East St').blur();
    cy.get('input[name="storageFacility.address.streetAddress2"]').clear().type('Suite 7A').blur();
    cy.get('input[name="storageFacility.address.city"]').clear().type('Sample City').blur();
    cy.get('select[name="storageFacility.address.state"]').select('GA');
    cy.get('input[name="storageFacility.address.postalCode"]').clear().type('30301').blur();
    cy.get('#facilityLotNumber').clear().type('1111111').blur();

    // Delivery info
    cy.get('#requestedDeliveryDate').clear().type('16 Mar 2022').blur();

    cy.get('input[name="delivery.address.streetAddress1"]').clear().type('148 S East St').blur();
    cy.get('input[name="delivery.address.streetAddress2"]').clear().type('Suite 7A').blur();
    cy.get('input[name="delivery.address.city"]').clear().type('Sample City').blur();
    cy.get('select[name="delivery.address.state"]').select('GA');
    cy.get('input[name="delivery.address.postalCode"]').clear().type('30301').blur();
    cy.get('#destinationType').select('HOME_OF_RECORD');

    // TAC and SAC
    cy.get('[data-testid="radio"] [for="tacType-NTS"]').click();
    cy.get('[data-testid="radio"] [for="sacType-HHG"]').click();

    cy.get('[data-testid="submitForm"]').click();

    cy.wait('@patchShipment');

    // edit the NTS shipment to be handled by an external vendor
    cy.get('[data-testid="ShipmentContainer"] .usa-button').last().click();
    cy.wait(['@getMTOShipments', '@getMoves']);

    cy.get('label[for="vendorExternal"]').click();
    cy.get('[data-testid="submitForm"]').click();

    cy.wait('@patchShipment');
    cy.get('[data-testid="ShipmentContainer"] [data-testid="tag"]').contains('external vendor');

    // edit the NTS shipment back to being handled by the GHC Prime contractor
    cy.get('[data-testid="ShipmentContainer"] .usa-button').last().click();
    cy.wait(['@getMTOShipments', '@getMoves']);

    cy.get('[data-testid="alert"]').contains('The GHC prime contractor is not handling the shipment.');

    cy.get('label[for="vendorPrime"]').click();
    cy.get('[data-testid="submitForm"]').click();

    cy.wait('@patchShipment');

    cy.get('[data-testid="ShipmentContainer"]')
      .last()
      .within(() => {
        cy.get('[data-testid="shipment-display-checkbox"]').should('be.enabled');
        cy.get('[data-testid="tag"]').should('not.exist');
        cy.get('[data-icon="chevron-down"]').click();

        cy.get('div[class*="missingInfoError"] [data-testid="ntsRecordedWeight"]').should('not.exist');
        cy.get('div[class*="missingInfoError"] [data-testid="storageFacilityName"]').should('not.exist');
        cy.get('div[class*="missingInfoError"] [data-testid="serviceOrderNumber"]').should('not.exist');
        cy.get('div[class*="missingInfoError"] [data-testid="storageFacilityAddress"]').should('not.exist');
        cy.get('div[class*="missingInfoError"] [data-testid="tacType"]').should('not.exist');
      });
  });

  it('TOO can approve an NTS shipment handled by the Prime', () => {
    // Make sure that it shows all the relevant information on the approve page
    // captures the information about the NTS Facility and relevant storage information
    // verifies that all the information is shown for NTS shipments handled by the GHC Contractor
    navigateToMove('PRINTR');

    cy.get('#approved-shipments').should('not.exist');
    cy.get('#requested-shipments');
    cy.contains('Approve selected').should('be.disabled');

    // Select & approve items
    cy.get('input[data-testid="shipment-display-checkbox"]').should('have.length', 2);
    cy.get('input[data-testid="shipment-display-checkbox"]').then(($shipmentsCheckbox) => {
      // Select each shipment
      $shipmentsCheckbox.each((i, el) => {
        const { id } = el;
        cy.get(`label[for="${id}"]`).click({ force: true }); // force because of shipment wrapping bug
      });
      // Select additional service items
      cy.get('label[for="shipmentManagementFee"]').click();
      cy.get('label[for="counselingFee"]').click();
      // Open modal
      const button = cy.contains('Approve selected');
      button.should('be.enabled');
      button.click();
      cy.get('#approvalConfirmationModal [data-testid="modal"]').then(($modal) => {
        cy.get($modal).should('be.visible');
        // Verify modal content
        cy.contains('Preview and post move task order');
        cy.get('#approvalConfirmationModal [data-testid="ShipmentContainer"]').should(
          'have.length',
          $shipmentsCheckbox.length,
        );
        cy.contains('Approved service items for this move')
          .next('table')
          .should('contain', 'Move management')
          .and('contain', 'Counseling');

        cy.get('[data-testid="modal"] [data-testid="ShipmentContainer"]')
          .last()
          .within(() => {
            cy.get('[data-testid="ntsRecordedWeight"]').should('exist');
            cy.get('[data-testid="destinationAddress"]').should('exist');
            cy.get('[data-testid="tacType"]').should('exist');
            cy.get('[data-testid="sacType"]').should('exist');
          });
        // Click approve
        cy.contains('Approve and send').click();
      });
    });

    cy.wait(['@approveShipment', '@patchMTOStatus']);

    // Redirected to Move Task Order page
    cy.url().should('include', `/moves/PRINTR/mto`);
  });

  it('TOO can view and edit Domestic NTS Shipments handled by the Prime on the MTO page', () => {
    navigateToMove('PRINTR');
    cy.get('[data-testid="MoveTaskOrder-Tab"]').click();
    cy.wait(['@getMTOShipments', '@getMTOServiceItems']);

    cy.get('[id="move-weights"] div').contains('1 shipment not moved by GHC prime.').should('not.exist');

    cy.get('[data-testid="ShipmentContainer"]').should('have.length', 2);

    cy.get('[data-testid="ShipmentContainer"]')
      .last()
      .within(() => {
        // non-temp storage header
        cy.get('h2').contains('Non-temp storage');
        // pickup address header
        cy.get('[class*="ShipmentAddresses_mtoShipmentAddresses"]').contains('Delivery address');
        // facility address header
        cy.get('[class*="ShipmentAddresses_mtoShipmentAddresses"]').contains('Facility address');

        // Confirming expected data
        // facility info and address
        // service order number
        // accounting codes
        cy.get('[class*="ShipmentDetailsSidebar"]').within(() => {
          cy.get('section header').first().contains('Facility info and address');
          cy.get('section').first().contains('148 S East St');

          cy.get('section header').eq(1).contains('Service order number');
          cy.get('section').eq(1).contains('999999');

          cy.get('section header').last().contains('Accounting codes');
          cy.get('section').last().contains('F123');
        });

        // edit facility info and address
        cy.get('[data-testid="edit-facility-info-modal-open"]').click();
      });

    cy.get('[data-testid="modal"]').within(($modal) => {
      cy.get($modal).should('be.visible');
      // Storage facility info
      cy.get('#facilityName').clear().type('New Facility Name').blur();
      cy.get('#facilityPhone').clear().type('999-999-9999').blur();
      cy.get('#facilityEmail').clear().type('new@example.com').blur();
      cy.get('#facilityServiceOrderNumber').clear().type('098098').blur();

      // Storage facility address
      cy.get('input[name="storageFacility.address.streetAddress1"]').clear().type('265 S East St').blur();
      cy.get('#facilityLotNumber').clear().type('1111111').blur();

      cy.get('button[type="submit"]').click();
    });

    cy.get('[data-testid="ShipmentContainer"]')
      .last()
      .within(() => {
        cy.get('[class*="ShipmentDetailsSidebar"]').within(() => {
          cy.get('section header').first().contains('Facility info and address');
          cy.get('section').first().contains('New Facility Name');
          cy.get('section').first().contains('265 S East St');
          cy.get('section').first().contains('Lot 1111111');
        });

        // edit service order number
        cy.get('[data-testid="service-order-number-modal-open"]').click();
      });

    cy.get('[data-testid="modal"]').within(($modal) => {
      cy.get($modal).should('be.visible');

      cy.get('[data-testid="textInput"]').clear().type('ORDER456').blur();

      cy.get('button[type="submit"]').click();
    });

    cy.get('[data-testid="ShipmentContainer"]')
      .last()
      .within(() => {
        cy.get('[class*="ShipmentDetailsSidebar"] section').eq(1).contains('ORDER456');

        // edit accounting codes
        cy.get('[data-testid="edit-accounting-codes-modal-open"]').click();
      });

    cy.get('[data-testid="modal"]').within(($modal) => {
      cy.get($modal).should('be.visible');

      cy.get('[data-testid="radio"] [for="tacType-HHG"]').click();
      cy.get('[data-testid="radio"] [for="sacType-NTS"]').click();

      cy.get('button[type="submit"]').click();
    });

    cy.get('[data-testid="ShipmentContainer"]')
      .last()
      .within(() => {
        cy.get('[class*="ShipmentDetailsSidebar"]').within(() => {
          cy.get('section').last().contains('F123');
          cy.get('section').last().contains('4K988AS098F');
        });

        cy.get('[data-testid="ApprovedServiceItemsTable"] h3').last().contains('Approved service items (5 items)');
      });
  });

  it('TOO can approve an HHG shipment with an NTS Shipment handled by an external vendor', () => {
    navigateToMove('PRXNTR');

    cy.get('#approved-shipments').should('not.exist');
    cy.get('#requested-shipments');
    cy.contains('Approve selected').should('be.disabled');

    cy.get('[data-testid="ShipmentContainer"]')
      .last()
      .within(() => {
        cy.get('[data-icon="chevron-down"]').click();

        cy.get('div[class*="missingInfoError"] [data-testid="storageFacilityName"]').should('exist');
        cy.get('[data-testid="storageFacilityName"]').contains('Missing');
        cy.get('div[class*="missingInfoError"] [data-testid="serviceOrderNumber"]').should('exist');
        cy.get('div[class*="missingInfoError"] [data-testid="storageFacilityAddress"]').should('exist');
        cy.get('[data-testid="storageFacilityAddress"]').contains('Missing');
        cy.get('div[class*="missingInfoError"] [data-testid="tacType"]').should('exist');
      });

    editTacSac();

    // Select & approve items
    cy.get('input[data-testid="shipment-display-checkbox"]').then(($shipmentsCheckbox) => {
      // Select each shipment
      $shipmentsCheckbox.each((i, el) => {
        const { id } = el;
        cy.get(`label[for="${id}"]`).click({ force: true }); // force because of shipment wrapping bug
      });
      // Select additional service items
      cy.get('label[for="shipmentManagementFee"]').click();
      cy.get('label[for="counselingFee"]').click();
      // Open modal
      const button = cy.contains('Approve selected');
      button.should('be.enabled');
      button.click();
      cy.get('#approvalConfirmationModal [data-testid="modal"]').then(($modal) => {
        cy.get($modal).should('be.visible');
        // Verify modal content
        cy.contains('Preview and post move task order');
        cy.get('#approvalConfirmationModal [data-testid="ShipmentContainer"]').should('have.length', 1);
        cy.contains('Approved service items for this move')
          .next('table')
          .should('contain', 'Move management')
          .and('contain', 'Counseling');

        // Click approve
        cy.contains('Approve and send').click();
      });
    });

    cy.wait(['@approveShipment', '@patchMTOStatus']);

    // Redirected to Move Task Order page
    cy.url().should('include', `/moves/PRXNTR/mto`);
    cy.wait(['@getMTOShipments']);
    // Confirm estimated weight shows expected extra shipment detail link
    cy.get('[id="move-weights"] div').contains('1 shipment not moved by GHC prime.');
  });

  it('TOO can submit service items on an NTS-only move handled by an external vendor', () => {
    navigateToMove('EXTNTR');

    cy.get('#approved-shipments').should('not.exist');
    cy.get('#requested-shipments');
    cy.contains('Approve selected').should('be.disabled');

    cy.get('[data-testid="ShipmentContainer"]')
      .last()
      .within(() => {
        cy.get('[data-testid="shipment-display-checkbox"]').should('not.exist');
        cy.get('[data-icon="chevron-down"]').click();

        cy.get('div[class*="missingInfoError"] [data-testid="storageFacilityName"]').should('exist');
        cy.get('[data-testid="storageFacilityName"]').contains('Missing');
        cy.get('div[class*="missingInfoError"] [data-testid="serviceOrderNumber"]').should('exist');
        cy.get('div[class*="missingInfoError"] [data-testid="storageFacilityAddress"]').should('exist');
        cy.get('[data-testid="storageFacilityAddress"]').contains('Missing');
        cy.get('div[class*="missingInfoError"] [data-testid="tacType"]').should('exist');
      });

    editTacSac();

    // Select & approve items
    cy.get('input[data-testid="shipment-display-checkbox"]').should('not.exist');
    // Select additional service items
    cy.get('label[for="shipmentManagementFee"]').click();
    cy.get('label[for="counselingFee"]').click();
    // Open modal
    const button = cy.contains('Approve selected');
    button.should('be.enabled');
    button.click();
    cy.get('#approvalConfirmationModal [data-testid="modal"]').then(($modal) => {
      cy.get($modal).should('be.visible');
      // Verify modal content
      cy.contains('Preview and post move task order');
      cy.get('#approvalConfirmationModal [data-testid="ShipmentContainer"]').should('not.exist');
      cy.contains('Approved service items for this move')
        .next('table')
        .should('contain', 'Move management')
        .and('contain', 'Counseling');

      // Click approve
      cy.contains('Approve and send').click();
    });

    cy.wait(['@patchMTOStatus']);

    // Redirected to Move Task Order page
    cy.url().should('include', `/moves/EXTNTR/mto`);

    cy.get('[role="main"]').contains('This move does not have any approved shipments yet.');
  });
});

function editTacSac() {
  cy.get('[data-testid="ShipmentContainer"] .usa-button').last().click();
  cy.wait(['@getMTOShipments']);

  cy.get('[data-testid="grid"] button').contains('Add or edit codes').click();

  cy.get('form').within(($form) => {
    cy.get('select[name="departmentIndicator"]').select('21 Army', { force: true });
    cy.get('input[name="ordersNumber"]').click().clear().type('ORDER66');
    cy.get('select[name="ordersType"]').select('Permanent Change Of Station (PCS)');
    cy.get('select[name="ordersTypeDetail"]').select('Shipment of HHG Permitted');

    cy.get('[data-testid="hhgTacInput"]').click().clear().type('E15A');
    cy.get('[data-testid="hhgSacInput"]').click().clear().type('4K988AS098F');
    cy.get('[data-testid="ntsTacInput"]').click().clear().type('F123');
    cy.get('[data-testid="ntsSacInput"]').click().clear().type('3L988AS098F');
    // Edit orders page | Save
    cy.get('[data-testid="button"]').contains('Save').click();
  });

  cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);

  cy.get('[data-testid="tacMDC"]').contains('E15A');
  cy.get('[data-testid="sacSDN"]').contains('4K988AS098F');
  cy.get('[data-testid="NTStac"]').contains('F123');
  cy.get('[data-testid="NTSsac"]').contains('3L988AS098F');
}

function navigateToMove(moveLocator) {
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
