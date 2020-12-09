import { TOOOfficeUserType } from '../../../support/constants';

describe('TOO user', () => {
  before(() => {
    cy.prepareOfficeApp();
  });

  beforeEach(() => {
    cy.intercept('**/ghc/v1/swagger.yaml').as('getGHCClient');
    cy.intercept('**/ghc/v1/queues/moves?**').as('getMoveOrders');
    cy.intercept('**/ghc/v1/queues/moves?page=1&perPage=20&sort=status&order=asc').as('getSortedMoveOrders');
    cy.intercept('**/ghc/v1/move-orders/**/move-task-orders').as('getMoveTaskOrders');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_shipments').as('getMTOShipments');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_service_items').as('getMTOServiceItems');
    cy.intercept('PATCH', '**/ghc/v1/move_task_orders/**/mto_shipments/**/status').as('patchMTOShipmentStatus');
    cy.intercept('PATCH', '**/ghc/v1/move-task-orders/**/status').as('patchMTOStatus');
    cy.intercept('PATCH', '**/ghc/v1/move-task-orders/**/service-items/**/status').as('patchMTOServiceItems');

    const userId = 'dcf86235-53d3-43dd-8ee8-54212ae3078f';
    cy.apiSignInAsUser(userId, TOOOfficeUserType);
    cy.wait(['@getMoveOrders']);
  });

  // This test performs a mutation so it can only succeed on a fresh DB.
  it('is able to approve a shipment', () => {
    const moveOrderId = '6fca843a-a87e-4752-b454-0fac67aa4988';
    const moveLocator = 'TEST12';

    // TOO Moves queue
    cy.wait(['@getSortedMoveOrders']);
    cy.contains(moveLocator).click();
    cy.url().should('include', `/moves/${moveOrderId}/details`);

    // Move Details page
    cy.wait(['@getMoveTaskOrders', '@getMTOShipments', '@getMTOServiceItems']);
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
          .should('contain', 'Shipment management fee')
          .and('contain', 'Counseling fee');
      });

      // Click approve
      cy.contains('Approve and send').click();
      cy.wait(['@patchMTOShipmentStatus', '@patchMTOStatus']);

      // Page refresh
      cy.url().should('include', `/moves/${moveOrderId}/details`);
      cy.get('#approvalConfirmationModal [data-testid="modal"]').should('not.exist');
      cy.wait(['@getMoveTaskOrders', '@getMTOShipments', '@getMTOServiceItems']);
      cy.get('#approvalConfirmationModal [data-testid="modal"]').should('not.exist');
      cy.get('#approved-shipments');
      cy.get('#requested-shipments').should('not.exist');
      cy.contains('Approve selected shipments').should('not.exist');
    });
  });

  it('is able to approve and reject mto service items', () => {
    const moveOrderId = '6fca843a-a87e-4752-b454-0fac67aa4988';
    const moveLocator = 'TEST12';

    // TOO Moves queue
    cy.wait(['@getSortedMoveOrders']);
    cy.contains(moveLocator).click();
    cy.url().should('include', `/moves/${moveOrderId}/details`);
    cy.get('[data-testid="MoveTaskOrder-Tab"]').click();
    cy.wait(['@getMoveTaskOrders', '@getMTOShipments', '@getMTOServiceItems']);
    cy.url().should('include', `/moves/${moveOrderId}/mto`);

    // Move Task Order page
    const shipments = cy.get('[data-testid="ShipmentContainer"]');
    shipments.should('have.length', 1);

    cy.contains('Requested service items (8 items)');
    cy.contains('Rejected service items').should('not.exist');
    cy.contains('Approved service items').should('not.exist');

    cy.get('[data-testid="modal"]').should('not.exist');

    // Approve a requested service item
    cy.get('[data-testid="RequestedServiceItemsTable"]').within(($table) => {
      cy.get('tbody tr').should('have.length', 8);
      cy.get('.acceptButton').first().click();
    });
    cy.contains('Approved service items (1 item)');
    cy.get('[data-testid="ApprovedServiceItemsTable"] tbody tr').should('have.length', 1);

    // Reject a requested service item
    cy.contains('Requested service items (7 items)');
    cy.get('[data-testid="RequestedServiceItemsTable"]').within(($table) => {
      cy.get('tbody tr').should('have.length', 7);
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

    cy.contains('Approved service items (2 items)');
    cy.get('[data-testid="ApprovedServiceItemsTable"] tbody tr').should('have.length', 2);
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

    // Approve the remaining service items
    cy.get('[data-testid="RequestedServiceItemsTable"] .acceptButton').each(($acceptBtn) => {
      $acceptBtn.trigger('click');
    });

    cy.contains('Requested service items').should('not.exist');
    cy.contains('Approved service items (7 items)');
    cy.get('[data-testid="ApprovedServiceItemsTable"] tbody tr').should('have.length', 7);
  });

  it('is able to edit allowances', () => {
    const moveOrderId = '6fca843a-a87e-4752-b454-0fac67aa4988';
    const moveLocator = 'TEST12';

    // TOO Moves queue
    cy.wait(['@getSortedMoveOrders']);
    cy.contains(moveLocator).click();
    cy.url().should('include', `/moves/${moveOrderId}/details`);

    // Move Details page
    cy.wait(['@getMoveTaskOrders', '@getMTOShipments', '@getMTOServiceItems']);

    // Edit allowances page | Save
    cy.get('[data-testid="edit-allowances"]').contains('Edit Allowances').click();
    cy.url().should('include', `/moves/${moveOrderId}/allowances`);
    cy.get('button').contains('Save').click();
    cy.url().should('include', `/moves/${moveOrderId}/details`);

    // Edit allowances page | Cancel
    cy.get('[data-testid="edit-allowances"]').contains('Edit Allowances').click();
    cy.get('button').contains('Cancel').click();
    cy.url().should('include', `/moves/${moveOrderId}/details`);
  });
});
