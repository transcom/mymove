import { TIOOfficeUserType } from '../../../support/constants';

describe('TIO user', () => {
  beforeEach(() => {
    cy.removeFetch();
    cy.server();
    cy.route('GET', '/ghc/v1/swagger.yaml').as('getGHCClient');
    cy.route('GET', '/ghc/v1/payment-requests').as('getPaymentRequests');
    cy.route('GET', '/ghc/v1/payment-requests/**').as('getPaymentRequest');
    cy.route('GET', '/ghc/v1/move_task_orders/**/mto_shipments').as('getMTOShipments');
    cy.route('GET', '/ghc/v1/move_task_orders/**/mto_service_items').as('getMTOServiceItems');

    cy.route('PATCH', '/ghc/v1/move-task-orders/**/payment-service-items/**/status').as(
      'patchPaymentServiceItemStatus',
    );
    cy.route('PATCH', '/ghc/v1/payment-requests/**/status').as('patchPaymentRequestStatus');

    const userId = '3b2cc1b0-31a2-4d1b-874f-0591f9127374';
    cy.signInAsUserPostRequest(TIOOfficeUserType, userId);
  });

  // This test performs a mutation so it can only succeed on a fresh DB.
  it('can review a payment request', () => {
    const paymentRequestId = 'a2c34dba-015f-4f96-a38b-0c0b9272e208';

    // TIO Payment Requests queue
    cy.wait(['@getGHCClient', '@getPaymentRequests']);
    cy.contains(paymentRequestId).parents('tr').contains('PENDING');
    cy.contains(paymentRequestId).click();

    // Payment Request detail page
    cy.wait(['@getPaymentRequest', '@getMTOShipments', '@getMTOServiceItems']);
    cy.url().should('include', `/payment-requests/${paymentRequestId}`);
    cy.get('[data-testid="ReviewServiceItems"]');

    // Approve #1
    cy.get('input[data-testid="approveRadio"]').then((el) => {
      const id = el.attr('id');
      cy.get(`label[for="${id}"]`).click();
      cy.wait('@patchPaymentServiceItemStatus');
    });
    cy.contains('Next').click();

    // Reject #2
    cy.get('input[data-testid="rejectRadio"]').then((el) => {
      const id = el.attr('id');
      cy.get(`label[for="${id}"]`).click();
      cy.contains('Save').click();
      cy.wait('@patchPaymentServiceItemStatus');
    });
    cy.contains('Next').click();

    // Approve 3-8
    cy.get('input[data-testid="approveRadio"]').then((el) => {
      const id = el.attr('id');
      cy.get(`label[for="${id}"]`).click();
      cy.wait('@patchPaymentServiceItemStatus');
    });
    cy.contains('Next').click();
    // 4
    cy.get('input[data-testid="approveRadio"]').then((el) => {
      const id = el.attr('id');
      cy.get(`label[for="${id}"]`).click();
      cy.wait('@patchPaymentServiceItemStatus');
    });
    cy.contains('Next').click();
    // 5
    cy.get('input[data-testid="approveRadio"]').then((el) => {
      const id = el.attr('id');
      cy.get(`label[for="${id}"]`).click();
      cy.wait('@patchPaymentServiceItemStatus');
    });
    cy.contains('Next').click();
    // 6
    cy.get('input[data-testid="approveRadio"]').then((el) => {
      const id = el.attr('id');
      cy.get(`label[for="${id}"]`).click();
      cy.wait('@patchPaymentServiceItemStatus');
    });
    cy.contains('Next').click();
    // 7
    cy.get('input[data-testid="approveRadio"]').then((el) => {
      const id = el.attr('id');
      cy.get(`label[for="${id}"]`).click();
      cy.wait('@patchPaymentServiceItemStatus');
    });
    cy.contains('Next').click();
    // 8
    cy.get('input[data-testid="approveRadio"]').then((el) => {
      const id = el.attr('id');
      cy.get(`label[for="${id}"]`).click();
      cy.wait('@patchPaymentServiceItemStatus');
    });
    cy.contains('Next').click();

    // Complete Request
    cy.contains('Complete request');
    cy.get('[data-testid="requested"]').contains('$2,169.88');
    cy.get('[data-testid="accepted"]').contains('$1,919.88');
    cy.get('[data-testid="rejected"]').contains('$250.00');
    cy.contains('Authorize payment').click();
    cy.wait('@patchPaymentRequestStatus');

    // Go back to queue
    cy.contains('Payment Requests');
    cy.contains(paymentRequestId).parents('tr').contains('REVIEWED');
  });
});
