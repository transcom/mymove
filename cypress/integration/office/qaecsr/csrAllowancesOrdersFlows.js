import { QAECSROfficeUserType } from '../../../support/constants';

describe('QAE/CSR orders and allowances read only view', () => {
  before(() => {
    cy.prepareOfficeApp();
  });

  beforeEach(() => {
    cy.intercept('**/ghc/v1/swagger.yaml').as('getGHCClient');
    cy.intercept('**/ghc/v1/moves/search').as('getSearchResults');
    cy.intercept('**/ghc/v1/move/**').as('getMoves');
    cy.intercept('GET', '**/ghc/v1/orders/**').as('getOrders');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_shipments').as('getMTOShipments');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_service_items').as('getMTOServiceItems');

    const userId = '2419b1d6-097f-4dc4-8171-8f858967b4db';
    cy.apiSignInAsUser(userId, QAECSROfficeUserType);
  });

  it('is able to see orders and form is read only', () => {
    const moveLocator = 'QAEHLP';
    cy.get('input[name="searchText"]').as('moveCodeSearchInput');
    cy.get('@moveCodeSearchInput').type(moveLocator).blur();
    cy.get('button').contains('Search').click();
    cy.wait(['@getSearchResults']);

    cy.get('tbody > tr').as('results');
    cy.get('@results').should('have.length', 1);
    cy.get('@results').first().contains(moveLocator);

    // click result to navigate to move details page
    cy.get('@results').first().click();
    cy.url().should('include', `/moves/${moveLocator}/details`);

    // Navigate to view orders page
    cy.get('[data-testid="view-orders"]').contains('View orders').click();

    cy.get('input[name="issueDate"]').should('be.disabled');
    cy.get('input[name="reportByDate"]').should('be.disabled');
    cy.get('select[name="departmentIndicator"]').should('be.disabled');
    cy.get('input[name="ordersNumber"]').should('be.disabled');
    cy.get('select[name="ordersType"]').should('be.disabled');
    cy.get('select[name="ordersTypeDetail"]').should('be.disabled');
    cy.get('input[name="tac"]').should('be.disabled');
    cy.get('input[name="sac"]').should('be.disabled');
    // no save button should exist
    cy.get('button').contains('Save').should('not.exist');
  });

  it('is able to see allowances and the form is read only', () => {
    const moveLocator = 'QAEHLP';
    cy.get('input[name="searchText"]').as('moveCodeSearchInput');
    cy.get('@moveCodeSearchInput').type(moveLocator).blur();
    cy.get('button').contains('Search').click();
    cy.wait(['@getSearchResults']);

    cy.get('tbody > tr').as('results');
    cy.get('@results').should('have.length', 1);
    cy.get('@results').first().contains(moveLocator);

    // click result to navigate to move details page
    cy.get('@results').first().click();
    cy.url().should('include', `/moves/${moveLocator}/details`);

    // Move Details page
    cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);

    // Navigate to view allowances page
    cy.get('[data-testid="view-allowances"]').contains('View allowances').click();

    cy.wait(['@getMoves', '@getOrders']);

    // read only pro-gear, pro-gear spouse, RME, SIT, and OCIE fields
    cy.get('input[name="proGearWeight"]').should('be.disabled');
    cy.get('input[name="proGearWeightSpouse"]').should('be.disabled');
    cy.get('input[name="requiredMedicalEquipmentWeight"]').should('be.disabled');
    cy.get('input[name="storageInTransit"]').should('be.disabled');
    cy.get('input[name="organizationalClothingAndIndividualEquipment"]').should('be.disabled');

    // read only grade and authorized weight
    cy.get('select[name=agency]').should('be.disabled');
    cy.get('select[name=agency]').should('be.disabled');
    cy.get('select[name="grade"]').should('be.disabled');
    cy.get('select[name="grade"]').should('be.disabled');
    cy.get('input[name="authorizedWeight"]').should('be.disabled');
    cy.get('input[name="dependentsAuthorized"]').should('be.disabled');

    // no save button should exist
    cy.get('button').contains('Save').should('not.exist');
  });
});
