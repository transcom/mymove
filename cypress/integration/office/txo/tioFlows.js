import { TIOOfficeUserType } from '../../../support/constants';

describe('TIO user', () => {
  before(() => {
    const userId = '3b2cc1b0-31a2-4d1b-874f-0591f9127374';
    cy.signInAsUserPostRequest(TIOOfficeUserType, userId);
  });

  beforeEach(() => {
    cy.persistSessionCookies();
  });

  describe('review payment request flow', () => {
    // This test performs a mutation so it can only succeed on a fresh DB.
    const paymentRequestId = 'a2c34dba-015f-4f96-a38b-0c0b9272e208';
    it('payment request status is pending', () => {
      cy.contains(paymentRequestId).parents('tr').contains('PENDING');
    });

    it('clicks on a pending payment request', () => {
      cy.contains(paymentRequestId).click();
      cy.url().should('include', `/payment-requests`);
    });

    it('shows the Review service items sidebar', () => {
      cy.get('[data-testid="ReviewServiceItems"]');
    });

    it('can approve the first service item', () => {
      cy.get('input[data-testid="approveRadio"]').then((el) => {
        const id = el.attr('id');
        cy.get(`label[for="${id}"]`).click();
      });
    });

    it('can go to the next service item', () => {
      cy.contains('Next').click();
    });

    it('can reject the second service item', () => {
      cy.get('input[data-testid="rejectRadio"]').then((el) => {
        const id = el.attr('id');
        cy.get(`label[for="${id}"]`).click();
        cy.contains('Save').click();
        cy.contains('Next').click();
      });
    });

    it('approves the remaining service items', () => {
      // 3
      cy.get('input[data-testid="approveRadio"]').then((el) => {
        const id = el.attr('id');
        cy.get(`label[for="${id}"]`).click();
      });
      cy.contains('Next').click();
      // 4
      cy.get('input[data-testid="approveRadio"]').then((el) => {
        const id = el.attr('id');
        cy.get(`label[for="${id}"]`).click();
      });
      cy.contains('Next').click();
      // 5
      cy.get('input[data-testid="approveRadio"]').then((el) => {
        const id = el.attr('id');
        cy.get(`label[for="${id}"]`).click();
      });
      cy.contains('Next').click();
      // 6
      cy.get('input[data-testid="approveRadio"]').then((el) => {
        const id = el.attr('id');
        cy.get(`label[for="${id}"]`).click();
      });
      cy.contains('Next').click();
      // 7
      cy.get('input[data-testid="approveRadio"]').then((el) => {
        const id = el.attr('id');
        cy.get(`label[for="${id}"]`).click();
      });
      cy.contains('Next').click();
      // 8
      cy.get('input[data-testid="approveRadio"]').then((el) => {
        const id = el.attr('id');
        cy.get(`label[for="${id}"]`).click();
      });
      cy.contains('Next').click();
    });

    it('shows the Complete Request step', () => {
      cy.contains('Complete request');
      cy.get('[data-testid="requested"]').contains('$2,169.88');
      cy.get('[data-testid="accepted"]').contains('$1,919.88');
      cy.get('[data-testid="rejected"]').contains('$250.00');
      cy.contains('Authorize payment').click();
    });

    it('navigates back to the payment request index page', () => {
      cy.contains('Payment Requests');
      cy.contains(paymentRequestId).parents('tr').contains('REVIEWED');
    });
  });
});
