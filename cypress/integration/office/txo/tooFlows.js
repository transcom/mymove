import { TOOOfficeUserType } from '../../../support/constants';

describe('TOO user', () => {
  before(() => {
    const userId = 'dcf86235-53d3-43dd-8ee8-54212ae3078f';
    // cy.signInAsNewTOOUser();
    cy.signInAsUserPostRequest(TOOOfficeUserType, userId);
  });

  beforeEach(() => {
    Cypress.Cookies.preserveOnce('masked_gorilla_csrf', 'office_session_token', '_gorilla_csrf');
  });

  describe('approve shipment flow', () => {
    /* This test performs a mutation so it can only succeed on a fresh DB.
     * This can be improved once the TOO queue UI contains information about
     * moves with requested shipments
     */
    it('clicks on a move with requested shipments', () => {
      const moveOrderId = '6fca843a-a87e-4752-b454-0fac67aa4988';
      cy.contains(moveOrderId).click();
      cy.url().should('include', `/moves/${moveOrderId}/details`);
    });

    it('does not show any Approved Shipments', () => {
      cy.contains('#approved-shipments').should('not.exist');
    });

    it('shows Requested Shipments', () => {
      cy.contains('#requested-shipments');
    });

    it('can’t click the Approve Shipments button yet', () => {
      cy.contains('Approve selected shipments').should('be.disabled');
    });

    it('doesn’t show the Approve Shipments modal', () => {
      cy.get('#approvalConfirmationModal [data-testid="modal"]').should('not.be.visible');
    });

    it('selects all shipments and open the Approve Shipments modal', () => {
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
          cy.get('#approvalConfirmationModal [data-testid="ShipmentContainer"]').should(
            'have.length',
            $shipments.length,
          );
          cy.contains('Approved service items for this move')
            .next('table')
            .should('contain', 'Shipment management fee')
            .and('contain', 'Counseling fee');
        });

        // Click approve
        cy.contains('Approve and send').click();
      });
    });

    it('the page refreshes to show new data', () => {
      cy.get('#approvalConfirmationModal [data-testid="modal"]').should('not.be.visible');
      cy.contains('#approved-shipments');
      cy.contains('#requested-shipments').should('not.exist');
      cy.contains('Approve selected shipments').should('not.exist');
    });
  });
});
