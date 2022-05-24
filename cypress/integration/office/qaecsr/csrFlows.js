import { QAECSROfficeUserType } from '../../../support/constants';

describe('Customer Support User Flows', () => {
  before(() => {
    cy.prepareOfficeApp();
  });

  beforeEach(() => {
    cy.intercept('**/ghc/v1/swagger.yaml').as('getGHCClient');
    cy.intercept('**/ghc/v1/queues/moves?page=1&perPage=20&sort=status&order=asc').as('getSortedOrders');
    cy.intercept('**/ghc/v1/move/**').as('getMoves');
    cy.intercept('GET', '**/ghc/v1/orders/**').as('getOrders');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_shipments').as('getMTOShipments');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_service_items').as('getMTOServiceItems');
    cy.intercept('**/ghc/v1/moves/**/customer-support-remarks').as('getCustomerSupportRemarks');

    // This user has multiple roles, which is the kind of user we use to test in staging.
    // By using this type of user, we can catch bugs like the one fixed in PR 6706.
    const userId = 'b264abd6-52fc-4e42-9e0f-173f7d217bc5';
    cy.apiSignInAsUser(userId, QAECSROfficeUserType);
  });

  // This test performs a mutation so it can only succeed on a fresh DB.
  it('is able to add a remark', () => {
    const moveLocator = 'TEST12';
    const testRemarkText = 'This is a test remark';

    // Moves queue (eventually will come via QAE/CSR move search)
    cy.wait(['@getSortedOrders']);
    cy.contains(moveLocator).click();
    cy.url().should('include', `/moves/${moveLocator}/details`);

    // Move Details page
    cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);

    // Go to Customer support remarks
    cy.contains('Customer support remarks').click();
    cy.url().should('include', `/moves/${moveLocator}/customer-support-remarks`);
    cy.wait(['@getCustomerSupportRemarks']);

    // Validate remarks page content
    cy.get('h1').contains('Customer support remarks');
    cy.get('h2').contains('Remarks');
    cy.get('h3').contains('Past remarks');
    cy.get('small').contains('Use this form to document any customer support provided for this move.');
    cy.get('[data-testid="textarea"]').should('have.attr', 'placeholder', 'Add your remarks here');

    cy.get('[data-testid=form] > [data-testid=button]').should('have.attr', 'disabled');

    // Should not have remarks (yet)
    cy.contains('No remarks yet').should('exist');

    // Add a remark
    cy.get('[data-testid="textarea"]').type(testRemarkText);
    cy.get('[data-testid=form] > [data-testid=button]').should('not.have.attr', 'disabled');
    cy.get('[data-testid=form] > [data-testid=button]').click();

    // Validate remark was added
    cy.wait(['@getCustomerSupportRemarks']);
    cy.contains('No remarks yet').should('not.exist');
    cy.contains(testRemarkText).should('exist');
  });
});
