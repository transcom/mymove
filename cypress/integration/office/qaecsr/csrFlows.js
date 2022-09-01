import { QAECSROfficeUserType, TIOOfficeUserType } from '../../../support/constants';
import { searchForAndNavigateToMove } from './qaeCSRIntegrationUtils';

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
    cy.intercept('**/ghc/v1/moves/search').as('getSearchResults');

    const userId = '2419b1d6-097f-4dc4-8171-8f858967b4db';
    cy.apiSignInAsUser(userId, QAECSROfficeUserType);
  });

  // This test performs a mutation so it can only succeed on a fresh DB.
  it('is able to add, edit, and delete a remark', () => {
    const moveLocator = 'TEST12';
    const testRemarkText = 'This is a test remark';
    const editString = '-edit';
    searchForAndNavigateToMove(moveLocator);

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

    // Open delete modal
    cy.get('[data-testid="modal"]').should('not.exist');
    cy.get('[data-testid="delete-remark-button"]').click();
    cy.get('[data-testid="modal"]').should('exist');
    cy.contains('Are you sure you want to delete this remark').should('exist');
    cy.contains('You cannot undo this action').should('exist');
    cy.contains('Yes, Delete').should('exist');
    cy.contains('No, keep it').should('exist');

    // Exit modal with cancel button
    cy.get('[data-testid=modalBackButton]').click();

    // Open the delete modal again
    cy.get('[data-testid="delete-remark-button"]').click();

    // Exit modal with the X button
    cy.get('[data-testid=modalCloseButton]').click();

    // Delete the remark for real
    cy.get('[data-testid="delete-remark-button"]').click();
    cy.contains('Yes, Delete').click();

    // Make sure success alert is shown
    cy.contains('Your remark has been deleted').should('exist');

    // Validate that the deleted remark is not on the page
    cy.wait(['@getCustomerSupportRemarks']);
    cy.contains(testRemarkText).should('not.exist');
    cy.contains('No remarks yet').should('exist');

    // Add a new remark
    cy.get('[data-testid="textarea"]').type(testRemarkText);
    cy.get('[data-testid=form] > [data-testid=button]').should('not.have.attr', 'disabled');
    cy.get('[data-testid=form] > [data-testid=button]').click();

    // Open edit and cancel
    cy.get('[data-testid="edit-remark-button"]').click();
    cy.get('[data-testid="edit-remark-textarea"]').type(editString);
    cy.get('[data-testid="edit-remark-cancel-button"]').click();

    // Validate remark was not edited
    cy.contains(testRemarkText).should('exist');
    cy.contains(testRemarkText + editString).should('not.exist');

    // Edit the remark
    cy.get('[data-testid="edit-remark-button"]').click();
    cy.get('[data-testid="edit-remark-textarea"]').type(editString);

    // Save the remark edit
    cy.get('[data-testid="edit-remark-save-button"]').click();

    // Validate remark was edited
    cy.wait(['@getCustomerSupportRemarks']);
    cy.contains(testRemarkText + editString).should('exist');
    cy.contains('(edited)').should('exist');

    // Changer user
    cy.contains('Sign out').click();
    cy.apiSignInAsUser('7f45b6bc-1131-4c9a-85ef-24552979d28d', QAECSROfficeUserType);
    searchForAndNavigateToMove(moveLocator);

    // Go to Customer support remarks
    cy.contains('Customer support remarks').click();
    cy.wait(['@getCustomerSupportRemarks']);

    // Edited remark should exist but no edit/delete buttons as I am a different user
    cy.contains(testRemarkText + editString).should('exist');
    cy.get('[data-testid="edit-remark-button"]').should('not.exist');
    cy.get('[data-testid="delete-remark-button"]').should('not.exist');
  });
});
