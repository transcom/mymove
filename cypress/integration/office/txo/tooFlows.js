import { TOOOfficeUserType } from '../../../support/constants';

describe('TOO user', () => {
  beforeEach(() => {
    cy.removeFetch();
    cy.server();
    cy.route('GET', '/api/v1/swagger.yaml').as('getGHCClient');
    cy.route('GET', '**/move-orders').as('getMoveOrders');

    const userId = 'dcf86235-53d3-43dd-8ee8-54212ae3078f';
    cy.signInAsUserPostRequest(TOOOfficeUserType, userId);
  });

  it('is able to approve a shipment', () => {
    cy.wait('@getGHCClient');
    cy.wait('@getMoveOrders');

    // This test performs a mutation so it can only succeed on a fresh DB.
    const moveOrderId = '6fca843a-a87e-4752-b454-0fac67aa4988';
    cy.contains(moveOrderId).click();
    cy.url().should('include', `/moves/${moveOrderId}/details`);

    cy.get('#approved-shipments').should('not.exist');

    cy.get('#requested-shipments');
    cy.contains('Approve selected shipments').should('be.disabled');
    cy.get('#approvalConfirmationModal [data-testid="modal"]').should('not.be.visible');
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
      cy.get('#approvalConfirmationModal [data-testid="modal"]').should('not.be.visible');
      cy.get('#approved-shipments');
      cy.get('#requested-shipments').should('not.exist');
      cy.contains('Approve selected shipments').should('not.exist');
    });
  });
});
