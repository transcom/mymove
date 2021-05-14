import { TIOOfficeUserType } from '../../../support/constants';

const completeServiceItemCard = ($serviceItem, approve = false) => {
  // serviceItemAmount
  if (!approve) {
    const inputEl = $serviceItem.find('input[data-testid="rejectRadio"]');
    const id = inputEl.attr('id');
    cy.get(`label[for="${id}"]`).click();

    cy.get('textarea[data-testid="textarea"]').type('This is not a valid request');

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
    cy.intercept('**/ghc/v1/swagger.yaml').as('getGHCClient');
    cy.intercept('**/ghc/v1/queues/payment-requests?**').as('getPaymentRequests');
    cy.intercept('**/ghc/v1/queues/payment-requests?sort=age&order=desc&page=1&perPage=20').as(
      'getSortedPaymentRequests',
    );
    cy.intercept('**/ghc/v1/moves/**/payment-requests').as('getMovePaymentRequests');
    cy.intercept('**/ghc/v1/payment-requests/**').as('getPaymentRequest');
    cy.intercept('**/ghc/v1/move/**').as('getMoves');
    cy.intercept('**/ghc/v1/orders/**').as('getOrders');
    cy.intercept('**/ghc/v1/documents/**').as('getDocuments');

    cy.intercept('PATCH', '**/ghc/v1/move-task-orders/**/payment-service-items/**/status').as(
      'patchPaymentServiceItemStatus',
    );
    cy.intercept('PATCH', '**/ghc/v1/payment-requests/**/status').as('patchPaymentRequestStatus');

    const userId = '3b2cc1b0-31a2-4d1b-874f-0591f9127374';
    cy.apiSignInAsUser(userId, TIOOfficeUserType);
  });

  // This test performs a mutation so it can only succeed on a fresh DB.
  it('can use a payment request page to update orders and review a payment request', () => {
    const paymentRequestId = 'ea945ab7-099a-4819-82de-6968efe131dc';

    // TIO Payment Requests queue
    cy.wait(['@getGHCClient', '@getPaymentRequests', '@getSortedPaymentRequests']);
    cy.get('#locator').type('TIOFLO');
    cy.get('th[data-testid="locator"]').first().click();
    cy.get('[data-testid="locator-0"]').click();

    // Payment Requests page
    cy.url().should('include', `/payment-requests`);
    cy.wait(['@getMovePaymentRequests']);
    cy.get('[data-testid="MovePaymentRequests"]');

    // View Orders page
    cy.contains('View orders').click();

    cy.wait(['@getMoves', '@getOrders', '@getDocuments']);
    cy.get('form').within(($form) => {
      cy.get('input[name="tac"]').click().clear().type('E15A');
      cy.get('input[name="sac"]').click().clear().type('4K988AS098F');
      // Edit orders page | Save
      cy.get('button').contains('Save').click();
    });

    cy.url().should('include', `/details`);

    cy.contains('Payment requests').click();

    // Payment Requests page
    cy.url().should('include', `/payment-requests`);
    cy.contains('Review service items').click();

    // Payment Request detail page
    cy.url().should('include', `/payment-requests/${paymentRequestId}`);
    cy.wait(['@getPaymentRequest']);
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

    // Returns to payment requests overview for move
    cy.url().should('include', `/payment-requests`);
    cy.wait(['@getMovePaymentRequests']);
    cy.get('[data-testid="MovePaymentRequests"]');
    cy.get('[data-testid="MovePaymentRequests"] [data-testid="tag"]').contains('Reviewed');
    cy.contains('Review Service Items').should('not.exist');

    // Go back to queue
    cy.get('a[title="Home"]').click();

    cy.contains('Payment requests', { matchCase: false });
    cy.contains('Reviewed', { matchCase: false });
    cy.get('[data-uuid="' + paymentRequestId + '"]').within(() => {
      cy.get('td').eq(2).contains('Reviewed');
    });
  });

  // This test performs a mutation so it can only succeed on a fresh DB.
  // This is a stripped down version of the above to build tests on for MB-7936
  it('can review a payment request', () => {
    const paymentRequestId = 'ea945ab7-099a-4819-82de-6968efe131dc';

    // TIO Payment Requests queue
    cy.wait(['@getGHCClient', '@getPaymentRequests', '@getSortedPaymentRequests']);
    cy.get('[data-uuid="' + paymentRequestId + '"]').click();

    // Payment Requests page
    cy.wait(['@getMovePaymentRequests']);
    cy.contains('Review service items').click();

    // Payment Request detail page
    cy.wait(['@getPaymentRequest']);

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

    cy.contains('Authorize payment').click();
    cy.wait('@patchPaymentRequestStatus');

    // Returns to payment requests overview for move
    cy.wait(['@getMovePaymentRequests']);
    cy.get('[data-testid="MovePaymentRequests"]');

    // Go back to queue
    cy.get('a[title="Home"]').click();

    cy.contains('Payment requests', { matchCase: false });
  });
});
