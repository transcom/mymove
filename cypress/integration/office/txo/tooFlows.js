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

  // This test performs a mutation so it can only succeed on a fresh DB.
  it('is able to approve a shipment', () => {
    const moveLocator = 'TEST12';

    // TOO Moves queue
    cy.wait(['@getSortedOrders']);
    cy.contains(moveLocator).click();
    cy.url().should('include', `/moves/${moveLocator}/details`);

    // Move Details page
    cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);
    cy.get('#approved-shipments').should('not.exist');
    cy.get('#requested-shipments');
    cy.contains('Approve selected').should('be.disabled');
    cy.get('#approvalConfirmationModal [data-testid="modal"]').should('not.be.visible');

    // Select & approve items
    cy.get('input[data-testid="shipment-display-checkbox"]').then(($shipments) => {
      // Select each shipment
      $shipments.each((i, el) => {
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
        cy.get('#approvalConfirmationModal [data-testid="ShipmentContainer"]').should('have.length', $shipments.length);
        cy.contains('Approved service items for this move')
          .next('table')
          .should('contain', 'Move management')
          .and('contain', 'Counseling');
      });

      // Click approve
      cy.contains('Approve and send').click();
      cy.wait(['@approveShipment', '@patchMTOStatus']);
    });

    // Redirected to Move Task Order page
    cy.url().should('include', `/moves/${moveLocator}/mto`);
    // cy.wait(['@getMTOShipments', '@getMTOServiceItems']);
    cy.get('[data-testid="ShipmentContainer"]');
    cy.get('[data-testid="ApprovedServiceItemsTable"] h3').contains('Approved service items (12 items)');

    // Navigate back to Move Details
    cy.get('[data-testid="MoveDetails-Tab"]').click();
    cy.url().should('include', `/moves/${moveLocator}/details`);
    cy.get('#approvalConfirmationModal [data-testid="modal"]').should('not.exist');
    cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);
    cy.get('#approvalConfirmationModal [data-testid="modal"]').should('not.exist');
    cy.get('#approved-shipments');
    cy.get('#requested-shipments').should('not.exist');
    cy.contains('Approve selected').should('not.exist');
  });

  it('is able to flag a move for financial review', () => {
    cy.wait(['@getSortedOrders']);
    // It doesn't matter which move we click on in the queue.
    cy.get('td').first().click();
    cy.url().should('include', `details`);
    cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);

    // click to trigger financial review modal
    cy.contains('Flag move for financial review').click();

    // Enter information in modal and submit
    cy.get('label').contains('Yes').click();
    cy.get('textarea').type('Something is rotten in the state of Denmark');

    // Click save on the modal
    cy.get('button').contains('Save').click();

    // Verify sucess alert and tag
    cy.contains('Move flagged for financial review.');
    cy.contains('Flagged for financial review');
  });

  it('is able to unflag a move for financial review', () => {
    cy.wait(['@getSortedOrders']);
    // It doesn't matter which move we click on in the queue.
    cy.get('td').first().click();
    cy.url().should('include', `details`);
    cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);

    // click to trigger financial review modal
    cy.contains('Edit').click();

    // Enter information in modal and submit
    cy.get('label').contains('No').click();

    // Click save on the modal
    cy.get('button').contains('Save').click();

    // Verify sucess alert and tag
    cy.contains('Move unflagged for financial review.');
  });

  it('is able to approve and reject mto service items', () => {
    const moveLocator = 'TEST12';

    // TOO Moves queue
    cy.wait(['@getSortedOrders']);
    cy.contains(moveLocator).click();
    cy.url().should('include', `/moves/${moveLocator}/details`);
    cy.get('[data-testid="MoveTaskOrder-Tab"]').click();
    cy.wait(['@getMTOShipments', '@getMTOServiceItems']);
    cy.url().should('include', `/moves/${moveLocator}/mto`);

    // Move Task Order page
    cy.get('[data-testid="ShipmentContainer"]').should('have.length', 1);

    cy.contains('Approved service items (12 items)');
    cy.contains('Rejected service items').should('not.exist');

    cy.get('[data-testid="modal"]').should('not.exist');

    // Approve a requested service item
    cy.get('[data-testid="RequestedServiceItemsTable"]').within(($table) => {
      cy.get('tbody tr').should('have.length', 2);
      cy.get('.acceptButton').first().click();
    });
    cy.contains('Approved service items (12 items)');
    cy.get('[data-testid="ApprovedServiceItemsTable"] tbody tr').should('have.length', 13);

    // Reject a requested service item
    cy.contains('Requested service items (1 item)');
    cy.get('[data-testid="RequestedServiceItemsTable"]').within(($table) => {
      cy.get('tbody tr').should('have.length', 1);
      cy.get('.rejectButton').first().click();
    });

    cy.get('[data-testid="modal"]').within(($modal) => {
      cy.get($modal).should('be.visible');
      cy.get('button[type="submit"]').should('be.disabled');
      cy.get('[data-testid="textInput"]').type('my very valid reason');
      cy.get('button[type="submit"]').click();
    });

    cy.get('[data-testid="modal"]').should('not.exist');

    cy.contains('Rejected service items (1 item)');
    cy.get('[data-testid="RejectedServiceItemsTable"] tbody tr').should('have.length', 1);

    // Accept a previously rejected service item
    cy.get('[data-testid="RejectedServiceItemsTable"] button').click();

    cy.contains('Approved service items (13 items)');
    cy.get('[data-testid="ApprovedServiceItemsTable"] tbody tr').should('have.length', 13);
    cy.contains('Rejected service items (1 item)').should('not.exist');

    // Reject a previously accpeted service item
    cy.get('[data-testid="ApprovedServiceItemsTable"] button').first().click();

    cy.get('[data-testid="modal"]').within(($modal) => {
      expect($modal).to.be.visible;
      cy.get('button[type="submit"]').should('be.disabled');
      cy.get('[data-testid="textInput"]').type('changed my mind about this one');
      cy.get('button[type="submit"]').click();
    });

    cy.get('[data-testid="modal"]').should('not.exist');

    cy.contains('Rejected service items (1 item)');
    cy.get('[data-testid="RejectedServiceItemsTable"] tbody tr').should('have.length', 1);

    cy.contains('Requested service items').should('not.exist');
    cy.contains('Approved service items (13 items)');
    cy.get('[data-testid="ApprovedServiceItemsTable"] tbody tr').should('have.length', 13);
  });

  it('is able to edit orders', () => {
    const moveLocator = 'TEST12';

    // TOO Moves queue
    cy.wait(['@getSortedOrders']);
    cy.contains(moveLocator).click();
    cy.url().should('include', `/moves/${moveLocator}/details`);

    // Move Details page
    cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);

    // Navigate to Edit orders page
    cy.get('[data-testid="edit-orders"]').contains('Edit orders').click();

    // Toggle between Edit Allowances and Edit Orders page
    cy.get('[data-testid="view-allowances"]').should('be.visible').click();
    cy.url().should('include', `/moves/${moveLocator}/allowances`);
    cy.get('[data-testid="view-orders"]').should('be.visible').click();
    cy.url().should('include', `/moves/${moveLocator}/orders`);

    // Edit orders fields

    cy.get('form').within(($form) => {
      cy.get('[class*="-control"]')
        .first()
        .click(0, 0)
        .type('Fort Irwin')
        .get('[class*="-menu"]')
        .find('[class*="-option"]')
        .first()
        .click(0, 0);

      cy.get('[class*="-control"]')
        .eq(1)
        .click(0, 0)
        .type('JB McGuire-Dix-Lakehurst')
        .get('[class*="-menu"]')
        .find('[class*="-option"]')
        .eq(5)
        .click(0, 0);

      cy.get('input[name="issueDate"]').click({ force: true }).clear().type('16 Mar 2018');
      cy.get('input[name="reportByDate"]').click({ force: true }).clear().type('22 Mar 2018');
      cy.get('select[name="departmentIndicator"]').select('21 Army', { force: true });
      cy.get('input[name="ordersNumber"]').click().clear().type('ORDER66');
      cy.get('select[name="ordersType"]').select('Permanent Change Of Station (PCS)');
      cy.get('select[name="ordersTypeDetail"]').select('Shipment of HHG Permitted');
      cy.get('input[name="tac"]').click().clear().type('F123');
      cy.get('input[name="sac"]').click().clear().type('4K988AS098F');

      // Edit orders page | Save
      cy.get('button').contains('Save').click();
    });

    // Verify edited values are saved
    cy.url().should('include', `/moves/${moveLocator}/details`);
    cy.get('[data-testid="currentDutyLocation"]').contains('Fort Irwin');
    cy.get('[data-testid="newDutyLocation"]').contains('Joint Base Lewis-McChord (McChord AFB)');
    cy.get('[data-testid="issuedDate"]').contains('16 Mar 2018');
    cy.get('[data-testid="reportByDate"]').contains('22 Mar 2018');
    cy.get('[data-testid="departmentIndicator"]').contains('Army');
    cy.get('[data-testid="ordersNumber"]').contains('ORDER66');
    cy.get('[data-testid="ordersType"]').contains('Permanent Change Of Station (PCS)');
    cy.get('[data-testid="ordersTypeDetail"]').contains('Shipment of HHG Permitted');
    cy.get('[data-testid="tacMDC"]').contains('F123');
    cy.get('[data-testid="sacSDN"]').contains('4K988AS098F');

    // Edit orders page | Cancel
    cy.get('[data-testid="edit-orders"]').contains('Edit orders').click();
    cy.get('button').contains('Cancel').click();
    cy.url().should('include', `/moves/${moveLocator}/details`);
  });

  it('is able to edit allowances', () => {
    const moveLocator = 'TEST12';

    // TOO Moves queue
    cy.wait(['@getSortedOrders']);
    cy.contains(moveLocator).click();
    cy.url().should('include', `/moves/${moveLocator}/details`);

    // Move Details page
    cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);

    // Navigate to Edit allowances page
    cy.get('[data-testid="edit-allowances"]').contains('Edit allowances').click();

    // Toggle between Edit Allowances and Edit Orders page
    cy.get('[data-testid="view-orders"]').click();
    cy.url().should('include', `/moves/${moveLocator}/orders`);
    cy.get('[data-testid="view-allowances"]').click();
    cy.url().should('include', `/moves/${moveLocator}/allowances`);

    cy.wait(['@getMoves', '@getOrders']);

    cy.get('form').within(($form) => {
      // Edit pro-gear, pro-gear spouse, RME, SIT, and OCIE fields
      cy.get('input[name="proGearWeight"]').clear().type('1999');
      cy.get('input[name="proGearWeightSpouse"]').clear().type('499');
      cy.get('input[name="requiredMedicalEquipmentWeight"]').clear().type('999');
      cy.get('input[name="storageInTransit"]').clear().type('199');
      cy.get('input[name="organizationalClothingAndIndividualEquipment"]').siblings('label[for="ocieInput"]').click();

      // Edit grade and authorized weight
      cy.get('select[name=agency]').contains('Army');
      cy.get('select[name=agency]').select('Navy');
      cy.get('select[name="grade"]').contains('E-1');
      cy.get('select[name="grade"]').select('W-2');
      cy.get('input[name="authorizedWeight"]').clear().type('11111');

      //Edit DependentsAuthorized
      cy.get('input[name="dependentsAuthorized"]').siblings('label[for="dependentsAuthorizedInput"]').click();

      // Edit allowances page | Save
      cy.get('button').contains('Save').should('be.enabled').click().should('be.disabled');
    });

    cy.wait(['@patchAllowances']);

    // Verify edited values are saved
    cy.url().should('include', `/moves/${moveLocator}/details`);

    cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);

    cy.get('[data-testid="progear"]').contains('1,999');
    cy.get('[data-testid="spouseProgear"]').contains('499');
    cy.get('[data-testid="rme"]').contains('999');
    cy.get('[data-testid="storageInTransit"]').contains('199');
    cy.get('[data-testid="ocie"]').contains('Unauthorized');

    cy.get('[data-testid="authorizedWeight"]').contains('11,111');
    cy.get('[data-testid="branchRank"]').contains('Navy');
    cy.get('[data-testid="branchRank"]').contains('W-2');
    cy.get('[data-testid="dependents"]').contains('Unauthorized');

    // Edit allowances page | Cancel
    cy.get('[data-testid="edit-allowances"]').contains('Edit allowances').click();
    cy.get('button').contains('Cancel').click();
    cy.url().should('include', `/moves/${moveLocator}/details`);
  });

  it('is able to edit shipment', () => {
    const moveLocator = 'TEST12';
    const deliveryDate = new Date().toLocaleDateString('en-US');

    // TOO Moves queue
    cy.wait(['@getSortedOrders']);
    cy.contains(moveLocator).click();
    cy.url().should('include', `/moves/${moveLocator}/details`);
    // Edit the shipment
    cy.get('[data-testid="ShipmentContainer"] .usa-button').first().click();
    // fill out some changes on the form
    cy.get('#requestedDeliveryDate').clear().type(deliveryDate).blur();
    cy.get('input[name="delivery.address.streetAddress1"]').clear().type('7 q st');
    cy.get('input[name="delivery.address.city"]').clear().type('city');
    cy.get('select[name="delivery.address.state"]').select('OH');
    cy.get('input[name="delivery.address.postalCode"]').clear().type('90210');
    cy.get('[data-testid="submitForm"]').click();
    // the shipment should be saved successfully
    cy.wait('@patchShipment');
  });

  it('is able to edit shipment for retiree', () => {
    const moveLocator = 'R3T1R3';
    const deliveryDate = new Date().toLocaleDateString('en-US');

    // TOO Moves queue
    cy.wait(['@getSortedOrders']);
    cy.contains(moveLocator).click();
    cy.url().should('include', `/moves/${moveLocator}/details`);
    // Edit the shipment
    cy.get('[data-testid="ShipmentContainer"] .usa-button').first().click();
    // fill out some changes on the form
    cy.get('#requestedDeliveryDate').clear().type(deliveryDate).blur();
    cy.get('input[name="delivery.address.streetAddress1"]').clear().type('7 q st');
    cy.get('input[name="delivery.address.city"]').clear().type('city');
    cy.get('select[name="delivery.address.state"]').select('OH');
    cy.get('input[name="delivery.address.postalCode"]').clear().type('90210');
    cy.get('select[name="destinationType"]').select('HOME_OF_SELECTION');
    cy.get('[data-testid="submitForm"]').click();
    // the shipment should be saved successfully
    cy.wait('@patchShipment');
  });

  it('is able to request cancellation for a shipment', () => {
    const moveLocator = 'TEST12';

    // TOO Moves queue
    cy.wait(['@getSortedOrders']);
    cy.contains(moveLocator).click();
    cy.url().should('include', `/moves/${moveLocator}/details`);
    cy.get('[data-testid="MoveTaskOrder-Tab"]').click();
    // cy.wait(['@getMTOShipments', '@getMTOServiceItems']);
    cy.url().should('include', `/moves/${moveLocator}/mto`);

    // Move Task Order page
    const shipments = cy.get('[data-testid="ShipmentContainer"]');
    shipments.should('have.length', 1);

    // Click requestCancellation button and display modal
    cy.get('.shipment-heading').find('button').should('contain', 'Request Cancellation').click();

    cy.get('[data-testid="modal"]').within(($modal) => {
      expect($modal).to.be.visible;
      cy.get('button[type="submit"]').should('exist');
      cy.get('button[type="submit"]').click();
    });

    cy.wait(['@requestShipmentCancellation']);
    // After updating, the button is disabeld and an alert is shown
    cy.get('[data-testid="request-cancellation-modal"]').should('not.exist');
    cy.get('.shipment-heading').find('button').should('be.disabled').and('contain', 'Cancellation Requested');
    cy.get('[data-testid="alert"]')
      .should('exist')
      .and('contain', 'The request to cancel that shipment has been sent to the movers.');

    // Alert should disappear if focus changes
    cy.get('[data-testid="rejectTextButton"]').first().click();
    cy.get('[data-testid="closeRejectServiceItem"]').click();
    cy.get('[data-testid="alert"]').should('not.exist');
  });

  it('is able to view SIT and create and edit SIT extensions', () => {
    const moveLocator = 'TEST12';

    // TOO Moves queue
    cy.wait(['@getSortedOrders']);
    cy.contains(moveLocator).click();
    cy.url().should('include', `/moves/${moveLocator}/details`);
    cy.get('[data-testid="MoveTaskOrder-Tab"]').click();
    cy.wait(['@getMTOShipments', '@getMTOServiceItems']);
    cy.url().should('include', `/moves/${moveLocator}/mto`);

    // View SIT display
    cy.get('[data-testid="sitExtensions"]');

    // Total SIT
    cy.contains('379 authorized');
    // cy.contains('60 used');
    // cy.contains('210 remaining');
    // cy.contains('Ends 26 Apr 2022');

    // Current SIT
    cy.contains('Current location: destination');
    cy.contains('60');
    // cy.contains('29 Aug 2021');
    // cy.contains('Ends 26 Apr 2022');

    // Previous SIT
    // cy.contains('30 days at origin (30 Jul 2021 - 29 Aug 2021)');
    cy.contains('30 days at origin');

    // SIT extensions
    cy.contains('90 days added');
    // cy.contains('on 28 Sep 2021');
    cy.contains('Serious illness of the member');
    cy.contains('The customer requested an extension.');
    cy.contains('The service member is unable to move into their new home at the expected time');
  });
});
