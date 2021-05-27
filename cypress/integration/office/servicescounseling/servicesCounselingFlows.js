import { ServicesCounselorOfficeUserType } from '../../../support/constants';

describe('Services counselor user', () => {
  before(() => {
    cy.prepareOfficeApp();
  });

  beforeEach(() => {
    cy.intercept('**/ghc/v1/swagger.yaml').as('getGHCClient');
    cy.intercept('**/ghc/v1/queues/counseling?page=1&perPage=20&sort=submittedAt&order=asc').as('getSortedOrders');
    // cy.intercept('**/ghc/v1/move/**').as('getMoves');
    // cy.intercept('**/ghc/v1/orders/**').as('getOrders');
    // cy.intercept('**/ghc/v1/orders/**/move-task-orders').as('getMoveTaskOrders');
    // cy.intercept('**/ghc/v1/move_task_orders/**/mto_shipments').as('getMTOShipments');
    // cy.intercept('**/ghc/v1/move_task_orders/**/mto_service_items').as('getMTOServiceItems');
    // cy.intercept('PATCH', '**/ghc/v1/move_task_orders/**/mto_shipments/**/status').as('patchMTOShipmentStatus');
    // cy.intercept('PATCH', '**/ghc/v1/move-task-orders/**/status').as('patchMTOStatus');
    // cy.intercept('PATCH', '**/ghc/v1/move-task-orders/**/service-items/**/status').as('patchMTOServiceItems');

    const userId = 'a6c8663f-998f-4626-a978-ad60da2476ec';
    cy.apiSignInAsUser(userId, ServicesCounselorOfficeUserType);
  });
  // it('is able to edit orders', () => {
  //   const moveLocator = 'TEST12';
  //
  //   // TOO Moves queue
  //   cy.wait(['@getSortedOrders']);
  //   cy.contains(moveLocator).click();
  //   cy.url().should('include', `/moves/${moveLocator}/details`);
  //
  //   // Move Details page
  //   cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);
  //
  //   // Navigate to Edit orders page
  //   cy.get('[data-testid="edit-orders"]').contains('Edit orders').click();
  //
  //   // Toggle between Edit Allowances and Edit Orders page
  //   cy.get('[data-testid="view-allowances"]').click();
  //   cy.url().should('include', `/moves/${moveLocator}/allowances`);
  //   cy.get('[data-testid="view-orders"]').click();
  //   cy.url().should('include', `/moves/${moveLocator}/orders`);
  //
  //   // Edit orders fields
  //
  //   cy.get('form').within(($form) => {
  //     cy.get('[class*="-control"]')
  //       .first()
  //       .click(0, 0)
  //       .type('Fort Irwin')
  //       .get('[class*="-menu"]')
  //       .find('[class*="-option"]')
  //       .first()
  //       .click(0, 0);
  //
  //     cy.get('[class*="-control"]')
  //       .eq(1)
  //       .click(0, 0)
  //       .type('JB McGuire-Dix-Lakehurst')
  //       .get('[class*="-menu"]')
  //       .find('[class*="-option"]')
  //       .eq(1)
  //       .click(0, 0);
  //
  //     cy.get('input[name="issueDate"]').click({ force: true }).clear().type('16 Mar 2018');
  //     cy.get('input[name="reportByDate"]').click({ force: true }).clear().type('22 Mar 2018');
  //     cy.get('select[name="departmentIndicator"]').select('21 Army', { force: true });
  //     cy.get('input[name="ordersNumber"]').click().clear().type('ORDER66');
  //     cy.get('select[name="ordersType"]').select('Permanent Change Of Station (PCS)');
  //     cy.get('select[name="ordersTypeDetail"]').select('Shipment of HHG Permitted');
  //     cy.get('input[name="tac"]').click().clear().type('F123');
  //     cy.get('input[name="sac"]').click().clear().type('4K988AS098F');
  //
  //     // Edit orders page | Save
  //     cy.get('button').contains('Save').click();
  //   });
  //
  //   // Verify edited values are saved
  //   cy.url().should('include', `/moves/${moveLocator}/details`);
  //   cy.get('[data-testid="currentDutyStation"]').contains('Fort Irwin');
  //   cy.get('[data-testid="newDutyStation"]').contains('JB Lewis-McChord');
  //   cy.get('[data-testid="issuedDate"]').contains('16 Mar 2018');
  //   cy.get('[data-testid="reportByDate"]').contains('22 Mar 2018');
  //   cy.get('[data-testid="departmentIndicator"]').contains('Army');
  //   cy.get('[data-testid="ordersNumber"]').contains('ORDER66');
  //   cy.get('[data-testid="ordersType"]').contains('Permanent Change Of Station (PCS)');
  //   cy.get('[data-testid="ordersTypeDetail"]').contains('Shipment of HHG Permitted');
  //   cy.get('[data-testid="tacMDC"]').contains('F123');
  //   cy.get('[data-testid="sacSDN"]').contains('4K988AS098F');
  //
  //   // Edit orders page | Cancel
  //   cy.get('[data-testid="edit-orders"]').contains('Edit orders').click();
  //   cy.get('button').contains('Cancel').click();
  //   cy.url().should('include', `/moves/${moveLocator}/details`);
  // });

  // it('is able to edit allowances', () => {
  //   const moveLocator = 'TEST12';
  //
  //   // TOO Moves queue
  //   cy.wait(['@getSortedOrders']);
  //   cy.contains(moveLocator).click();
  //   cy.url().should('include', `/moves/${moveLocator}/details`);
  //
  //   // Move Details page
  //   cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);
  //
  //   // Navigate to Edit allowances page
  //   cy.get('[data-testid="edit-allowances"]').contains('Edit allowances').click();
  //
  //   // Toggle between Edit Allowances and Edit Orders page
  //   cy.get('[data-testid="view-orders"]').click();
  //   cy.url().should('include', `/moves/${moveLocator}/orders`);
  //   cy.get('[data-testid="view-allowances"]').click();
  //   cy.url().should('include', `/moves/${moveLocator}/allowances`);
  //
  //   cy.get('form').within(($form) => {
  //     // Edit pro-gear, pro-gear spouse, RME, and OCIE fields
  //     cy.get('input[name="proGearWeight"]').clear().type('1999');
  //     cy.get('input[name="proGearWeightSpouse"]').clear().type('499');
  //     cy.get('input[name="requiredMedicalEquipmentWeight"]').clear().type('999');
  //     cy.get('input[name="organizationalClothingAndIndividualEquipment"]').click({ force: true });
  //
  //     // Edit grade and authorized weight
  //     cy.get('select[name=agency]').contains('Army');
  //     cy.get('select[name=agency]').select('Navy');
  //     cy.get('select[name="grade"]').contains('E-1');
  //     cy.get('select[name="grade"]').select('W-2');
  //     cy.get('input[name="authorizedWeight"]').clear().type('11111');
  //
  //     //Edit DependentsAuthorized
  //     cy.get('input[name="dependentsAuthorized"]').click({ force: true });
  //
  //     // Edit allowances page | Save
  //     cy.get('button').contains('Save').click();
  //   });
  //
  //   // Verify edited values are saved
  //   cy.url().should('include', `/moves/${moveLocator}/details`);
  //   cy.get('[data-testid="progear"]').contains('1,999');
  //   cy.get('[data-testid="spouseProgear"]').contains('499');
  //   cy.get('[data-testid="rme"]').contains('999');
  //   cy.get('[data-testid="ocie"]').contains('Unauthorized');
  //
  //   cy.get('[data-testid="authorizedWeight"]').contains('11,111');
  //   cy.get('[data-testid="branchRank"]').contains('Navy');
  //   cy.get('[data-testid="branchRank"]').contains('W-2');
  //   cy.get('[data-testid="dependents"]').contains('Unauthorized');
  //
  //   // Edit allowances page | Cancel
  //   cy.get('[data-testid="edit-allowances"]').contains('Edit allowances').click();
  //   cy.get('button').contains('Cancel').click();
  //   cy.url().should('include', `/moves/${moveLocator}/details`);
  // });
});
