import { TIOOfficeUserType } from '../../../support/constants';

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
    cy.intercept('**/ghc/v1/moves/**/billable-weight').as('getBillableWeight');

    const userId = '3b2cc1b0-31a2-4d1b-874f-0591f9127374';
    cy.apiSignInAsUser(userId, TIOOfficeUserType);
  });

  // This test performs a mutation so it can only succeed on a fresh DB.
  it('can update max billable weight', () => {
    const paymentRequestId = 'ea945ab7-099a-4819-82de-6968efe131dc';

    // TIO Payment Requests queue
    cy.wait(['@getGHCClient', '@getPaymentRequests', '@getSortedPaymentRequests']);
    cy.get('[data-uuid="' + paymentRequestId + '"]').click();

    // Payment Requests page
    cy.wait(['@getMovePaymentRequests']);

    cy.get('#billable-weights').contains('Review weights').click();

    cy.wait(['@getBillableWeight']);

    cy.get('[data-testid="maxBillableWeightEdit"]').contains('Edit').click();
  });
});
