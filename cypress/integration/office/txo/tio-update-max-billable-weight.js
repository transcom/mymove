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
    cy.intercept('**/ghc/v1/move-task-orders/**/billable-weights-reviewed-at').as('getBillableWeight');
    cy.intercept('**/ghc/v1/move/**').as('getMoves');
    cy.intercept('**/ghc/v1/orders/**').as('getOrders');
    cy.intercept('**/ghc/v1/orders/**/update-max-billable-weight/tio').as('updateBillableWeight');

    const userId = '3b2cc1b0-31a2-4d1b-874f-0591f9127374';
    cy.apiSignInAsUser(userId, TIOOfficeUserType);
  });

  // This test performs a mutation so it can only succeed on a fresh DB.
  it('can update max billable weight', () => {
    // Go to known NTS-R move
    cy.get('#locator').type('NTSRT2');
    cy.get('th[data-testid="locator"]').first().click();
    cy.get('[data-testid="locator-0"]').click();

    // Verify we are on the Payment Requests page
    cy.url().should('include', `/payment-requests`);
    cy.wait(['@getMovePaymentRequests', '@getMoves', '@getOrders']);
    cy.get('[data-testid="MovePaymentRequests"]');

    cy.get('#billable-weights').contains('Review weights').click();

    cy.wait(['@getBillableWeight']);

    cy.get('[data-testid="button"]').contains('Edit').click();

    cy.wait(['@getMoves']);
    cy.wait(250);

    cy.get('fieldset').within((/* $form */) => {
      cy.get('input#billableWeight').click().clear().type('7400');
      cy.get('textarea#billableWeightJustification').click().clear().type('Some basic remarks.');
    });

    cy.get('button').contains('Save changes').click();

    cy.wait(['@updateBillableWeight']);

    cy.get('[data-testid="billableWeightValue"]').contains('7,400 lbs');
    cy.get('[data-testid="billableWeightRemarks"]').contains('Some basic remarks.');
  });
});
