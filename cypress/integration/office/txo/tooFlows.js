import { TOOOfficeUserType } from '../../../support/constants';

describe('TOO user', () => {
  before(() => {
    cy.prepareOfficeApp();
  });

  beforeEach(() => {
    cy.intercept('**/ghc/v1/swagger.yaml').as('getGHCClient');
    cy.intercept('**/ghc/v1/queues/moves?page=1&perPage=20&sort=status&order=asc').as('getSortedOrders');
    cy.intercept('**/ghc/v1/move/**').as('getMoves');
    cy.intercept('**/ghc/v1/orders/**').as('getOrders');
    cy.intercept('**/ghc/v1/orders/**/move-task-orders').as('getMoveTaskOrders');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_shipments').as('getMTOShipments');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_service_items').as('getMTOServiceItems');
    cy.intercept('PATCH', '**/ghc/v1/move_task_orders/**/mto_shipments/**/status').as('patchMTOShipmentStatus');
    cy.intercept('PATCH', '**/ghc/v1/move-task-orders/**/status').as('patchMTOStatus');
    cy.intercept('PATCH', '**/ghc/v1/move-task-orders/**/service-items/**/status').as('patchMTOServiceItems');

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
    cy.contains('Approve selected shipments').should('be.disabled');
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
      const button = cy.contains('Approve selected shipments');
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
      cy.wait(['@patchMTOShipmentStatus', '@patchMTOStatus']);
    });

    // Redirected to Move Task Order page
    cy.url().should('include', `/moves/${moveLocator}/mto`);
    cy.wait(['@getMoveTaskOrders', '@getMTOShipments', '@getMTOServiceItems']);
    cy.get('[data-testid="ShipmentContainer"]');
    cy.get('[data-testid="ApprovedServiceItemsTable"] h4').contains('Approved service items (6 items)');

    // Navigate back to Move Details
    cy.get('[data-testid="MoveDetails-Tab"]').click();
    cy.url().should('include', `/moves/${moveLocator}/details`);
    cy.get('#approvalConfirmationModal [data-testid="modal"]').should('not.exist');
    cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);
    cy.get('#approvalConfirmationModal [data-testid="modal"]').should('not.exist');
    cy.get('#approved-shipments');
    cy.get('#requested-shipments').should('not.exist');
    cy.contains('Approve selected shipments').should('not.exist');
  });

  it('is able to approve and reject mto service items', () => {
    const moveLocator = 'TEST12';

    // TOO Moves queue
    cy.wait(['@getSortedOrders']);
    cy.contains(moveLocator).click();
    cy.url().should('include', `/moves/${moveLocator}/details`);
    cy.get('[data-testid="MoveTaskOrder-Tab"]').click();
    cy.wait(['@getMoveTaskOrders', '@getMTOShipments', '@getMTOServiceItems']);
    cy.url().should('include', `/moves/${moveLocator}/mto`);

    // Move Task Order page
    const shipments = cy.get('[data-testid="ShipmentContainer"]');
    shipments.should('have.length', 1);

    cy.contains('Approved service items (6 items)');
    cy.contains('Rejected service items').should('not.exist');

    cy.get('[data-testid="modal"]').should('not.exist');

    // Approve a requested service item
    cy.get('[data-testid="RequestedServiceItemsTable"]').within(($table) => {
      cy.get('tbody tr').should('have.length', 2);
      cy.get('.acceptButton').first().click();
    });
    cy.contains('Approved service items (7 items)');
    cy.get('[data-testid="ApprovedServiceItemsTable"] tbody tr').should('have.length', 7);

    // Reject a requested service item
    cy.contains('Requested service items (1 item)');
    cy.get('[data-testid="RequestedServiceItemsTable"]').within(($table) => {
      cy.get('tbody tr').should('have.length', 1);
      cy.get('.rejectButton').first().click();
    });

    cy.get('[data-testid="modal"]').within(($modal) => {
      expect($modal).to.be.visible;
      cy.get('button[type="submit"]').should('be.disabled');
      cy.get('[data-testid="textInput"]').type('my very valid reason');
      cy.get('button[type="submit"]').click();
    });

    cy.get('[data-testid="modal"]').should('not.exist');

    cy.contains('Rejected service items (1 item)');
    cy.get('[data-testid="RejectedServiceItemsTable"] tbody tr').should('have.length', 1);

    // Accept a previously rejected service item
    cy.get('[data-testid="RejectedServiceItemsTable"] button').click();

    cy.contains('Approved service items (8 items)');
    cy.get('[data-testid="ApprovedServiceItemsTable"] tbody tr').should('have.length', 8);
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
    cy.contains('Approved service items (7 items)');
    cy.get('[data-testid="ApprovedServiceItemsTable"] tbody tr').should('have.length', 7);
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
    cy.get('[data-testid="view-allowances"]').click();
    cy.url().should('include', `/moves/${moveLocator}/allowances`);
    cy.get('[data-testid="view-orders"]').click();
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
        .eq(1)
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
    cy.get('[data-testid="currentDutyStation"]').contains('Fort Irwin');
    cy.get('[data-testid="newDutyStation"]').contains('JB Lewis-McChord');
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
    cy.get('[data-testid="edit-allowances"]').contains('Edit Allowances').click();

    // Toggle between Edit Allowances and Edit Orders page
    cy.get('[data-testid="view-orders"]').click();
    cy.url().should('include', `/moves/${moveLocator}/orders`);
    cy.get('[data-testid="view-allowances"]').click();
    cy.url().should('include', `/moves/${moveLocator}/allowances`);

    cy.get('form').within(($form) => {
      // Edit grade and authorized weight
      cy.get('select[name=agency]').contains('Army');
      cy.get('select[name=agency]').select('Navy');
      cy.get('select[name="grade"]').contains('E-1');
      cy.get('select[name="grade"]').select('W-2');
      cy.get('input[name="authorizedWeight"]').clear().type('11111');

      //Edit DependentsAuthorized
      cy.get('input[name="dependentsAuthorized"]').click();

      // Edit allowances page | Save
      cy.get('button').contains('Save').click();
    });

    // Verify edited values are saved
    cy.url().should('include', `/moves/${moveLocator}/details`);
    cy.get('[data-testid="authorizedWeight"]').contains('11,111 lbs');
    cy.get('[data-testid="branchRank"]').contains('Navy');
    cy.get('[data-testid="branchRank"]').contains('W-2');
    cy.get('[data-testid="dependents"]').contains('Unauthorized');

    // Edit allowances page | Cancel
    cy.get('[data-testid="edit-allowances"]').contains('Edit Allowances').click();
    cy.get('button').contains('Cancel').click();
    cy.url().should('include', `/moves/${moveLocator}/details`);
  });

  it('is able to request cancellation for a shipment', () => {
    const moveLocator = 'TEST12';

    // TOO Moves queue
    cy.wait(['@getSortedOrders']);
    cy.contains(moveLocator).click();
    cy.url().should('include', `/moves/${moveLocator}/details`);
    cy.get('[data-testid="MoveTaskOrder-Tab"]').click();
    cy.wait(['@getMoveTaskOrders', '@getMTOShipments', '@getMTOServiceItems']);
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

    cy.wait(['@patchMTOShipmentStatus']);
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
});
