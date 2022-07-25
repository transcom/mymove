import { ServicesCounselorOfficeUserType } from '../../../support/constants';

describe('Services counselor user', () => {
  before(() => {
    cy.prepareOfficeApp();
  });

  beforeEach(() => {
    cy.intercept('**/ghc/v1/swagger.yaml').as('getGHCClient');
    cy.intercept('**/ghc/v1/queues/counseling?page=1&perPage=20&sort=submittedAt&order=asc').as('getSortedMoves');
    cy.intercept('**/ghc/v1/move/**').as('getMoves');
    cy.intercept('**/ghc/v1/orders/**').as('getOrders');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_shipments').as('getMTOShipments');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_service_items').as('getMTOServiceItems');
    cy.intercept('**/ghc/v1/move-task-orders/**/status/service-counseling-completed').as(
      'patchServiceCounselingCompleted',
    );
    cy.intercept('POST', '**/ghc/v1/mto-shipments').as('createShipment');
    cy.intercept('PATCH', '**/ghc/v1/move_task_orders/**/mto_shipments/**').as('patchShipment');
    cy.intercept('DELETE', '**/ghc/v1/shipments/**').as('deleteShipment');

    const userId = 'a6c8663f-998f-4626-a978-ad60da2476ec';
    cy.apiSignInAsUser(userId, ServicesCounselorOfficeUserType);
  });

  it('Services Counselor can add an NTS shipment to the customer move', () => {
    navigateToMove('NTSHHG');
    addNTSShipment();
  });

  it('Services Counselor can delete/remove an NTS shipment request', () => {
    navigateToMove('NTSHHG');
    addNTSShipment();

    cy.get('[data-testid="ShipmentContainer"] .usa-button').should('have.length', 3);

    cy.get('[data-testid="ShipmentContainer"] .usa-button').last().click();
    cy.wait(['@getMTOServiceItems']);
    // click to trigger confirmation modal
    cy.get('[data-testid="grid"] button').contains('Delete shipment').click();

    cy.get('[data-testid="modal"]').should('be.visible');

    cy.get('[data-testid="modal"] button').contains('Delete shipment').click({ force: true });
    cy.wait(['@deleteShipment', '@getMoves']);

    cy.get('[data-testid="ShipmentContainer"] .usa-button').should('have.length', 2);
  });

  it('Services Counselor can enter accounting codes on the Orders Page', () => {
    navigateToMove('NTSHHG');

    cy.get('[data-testid="ShipmentContainer"] .usa-button').last().click();
    cy.wait(['@getMTOShipments']);

    cy.get('[data-testid="grid"] button').contains('Add or edit codes').click();

    cy.get('form').within(($form) => {
      cy.get('input[name="tac"]').click().clear().type('E15A');
      cy.get('input[name="sac"]').click().clear().type('4K988AS098F');
      cy.get('input[name="ntsTac"]').click().clear().type('F123');
      cy.get('input[name="ntsSac"]').click().clear().type('3L988AS098F');
      cy.get('input[name="ordersNumber"]').click().clear().type('1234');
      cy.get('select[name="departmentIndicator"]').select('21 Army', { force: true });
      cy.get('select[name="ordersType"]').select('Permanent Change Of Station (PCS)');
      cy.get('select[name="ordersTypeDetail"]').select('Shipment of HHG Permitted');
      // Edit orders page | Save
      cy.get('button').contains('Save').click();
    });

    cy.get('[data-testid="tacMDC"]').contains('E15A');
    cy.get('[data-testid="sacSDN"]').contains('4K988AS098F');
    cy.get('[data-testid="NTStac"]').contains('F123');
    cy.get('[data-testid="NTSsac"]').contains('3L988AS098F');
  });

  it('Services Counselor can assign accounting code(s) to a shipment', () => {
    navigateToMove('NTSHHG');

    cy.get('[data-testid="ShipmentContainer"] .usa-button').last().click();
    cy.wait(['@getMTOShipments', '@getMoves']);

    cy.get('[data-testid="radio"] [for="tacType-NTS"]').click();
    cy.get('[data-testid="radio"] [for="sacType-HHG"]').click();

    cy.get('[data-testid="submitForm"]').click();

    cy.wait('@patchShipment');
    cy.get('.usa-alert__text').contains('Your changes were saved.');
  });

  it('Services Counselor can submit a move with an NTS shipment', () => {
    navigateToMove('NTSHHG');
    // click to trigger confirmation modal
    cy.contains('Submit move details').click();

    cy.get('[data-testid="modal"]').should('be.visible');

    cy.get('button').contains('Yes, submit').click();
    cy.wait(['@patchServiceCounselingCompleted', '@getMoves']);

    // verify success alert
    cy.contains('Move submitted.');
  });
});

function navigateToMove(moveLocator) {
  //SC Moves queue
  cy.wait(['@getSortedMoves']);
  cy.get('input[name="locator"]').as('moveCodeFilterInput');

  // type in move code/locator to filter
  cy.get('@moveCodeFilterInput').type(moveLocator).blur();

  cy.get('tbody > tr').as('results');

  // click result to navigate to move details page
  cy.get('@results').first().click();
  cy.url().should('include', `/counseling/moves/${moveLocator}/details`);
  cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);
}

function addNTSShipment() {
  cy.get('[data-testid="dropdown"]').first().select('NTS');

  cy.get('#requestedPickupDate').clear().type('16 Mar 2022').blur();
  cy.get('[data-testid="useCurrentResidence"]').click({ force: true });
  cy.get(`[data-testid="remarks"]`).should('not.exist');
  cy.get('[name="counselorRemarks"]').type('Sample counselor remarks');

  cy.get('[data-testid="submitForm"]').click();
  // the shipment should be saved with the type
  cy.wait('@createShipment');
  // the new shipment is visible on the Move details page
  cy.get('[data-testid="ShipmentContainer"]').last().contains('Sample counselor remarks');
}
