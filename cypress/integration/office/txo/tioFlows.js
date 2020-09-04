import { TIOOfficeUserType } from '../../../support/constants';

const completeServiceItemCard = ($serviceItem, approve = false) => {
  // serviceItemAmount
  if (!approve) {
    const inputEl = $serviceItem.find('input[data-testid="rejectRadio"]');
    const id = inputEl.attr('id');
    cy.get(`label[for="${id}"]`).click();
    cy.wrap($serviceItem).contains('Save').click();
  } else {
    const inputEl = $serviceItem.find('input[data-testid="approveRadio"]');
    const id = inputEl.attr('id');
    cy.get(`label[for="${id}"]`).click();
  }
};

describe('TIO user', () => {
  before(() => {
    cy.prepareOfficeApp();
  });

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
    cy.apiSignInAsUser(userId, TIOOfficeUserType);
  });

  // This test performs a mutation so it can only succeed on a fresh DB.
  it('can review a payment request', () => {
    const paymentRequestId = 'ea945ab7-099a-4819-82de-6968efe131dc';

    // TIO Payment Requests queue
    cy.wait(['@getGHCClient', '@getPaymentRequests']);
    cy.contains(paymentRequestId).parents('tr').contains('PENDING');
    cy.contains(paymentRequestId).click();

    // Payment Request detail page
    cy.url().should('include', `/payment-requests/${paymentRequestId}`);
    cy.wait(['@getPaymentRequest', '@getMTOShipments', '@getMTOServiceItems']);
    cy.get('[data-testid="ReviewServiceItems"]');

    // Approve the first service item
    cy.get('[data-testid="ServiceItemCard"]').each((el) => {
      completeServiceItemCard(el, true);
    });
    cy.wait('@patchPaymentServiceItemStatus');
    cy.contains('Next').click();

    // Reject the second
    cy.get('[data-testid="ServiceItemCard"]').each((el) => {
      completeServiceItemCard(el, false);
    });
    cy.wait('@patchPaymentServiceItemStatus');
    cy.contains('Next').click();

    // Complete Request
    cy.contains('Complete request');

    cy.get('[data-testid="requested"]').contains('$1,099.99');
    cy.get('[data-testid="accepted"]').contains('$100.00');
    cy.get('[data-testid="rejected"]').contains('$999.99');

    cy.contains('Authorize payment').click();
    cy.wait('@patchPaymentRequestStatus');

    // Go back to queue
    cy.contains('Payment Requests');
    cy.contains(paymentRequestId).parents('tr').contains('REVIEWED');
  });
});
