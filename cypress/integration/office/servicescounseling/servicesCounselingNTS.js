import { ServicesCounselorOfficeUserType } from '../../../support/constants';

describe('Services counselor user', () => {
  before(() => {
    cy.prepareOfficeApp();
  });

  beforeEach(() => {
    cy.intercept('**/ghc/v1/swagger.yaml').as('getGHCClient');
    cy.intercept('**/ghc/v1/queues/counseling?page=1&perPage=20&sort=submittedAt&order=asc').as('getSortedMoves');
    cy.intercept('**/ghc/v1/queues/counseling?page=1&perPage=20&sort=submittedAt&order=asc&locator=SCE4ET').as(
      'getFilterSortedMoves',
    );
    cy.intercept('**/ghc/v1/move/**').as('getMoves');
    cy.intercept('**/ghc/v1/orders/**').as('getOrders');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_shipments').as('getMTOShipments');
    cy.intercept('**/ghc/v1/move_task_orders/**/mto_service_items').as('getMTOServiceItems');
    cy.intercept('**/ghc/v1/move-task-orders/**/status/service-counseling-completed').as(
      'patchServiceCounselingCompleted',
    );
    cy.intercept('**/ghc/v1/moves/**/financial-review-flag').as('financialReviewFlagCompleted');
    cy.intercept('POST', '**/ghc/v1/mto-shipments').as('createShipment');
    cy.intercept('PATCH', '**/ghc/v1/move_task_orders/**/mto_shipments/**').as('patchShipment');
    cy.intercept('PATCH', '**/ghc/v1/counseling/orders/**/allowances').as('patchAllowances');

    const userId = 'a6c8663f-998f-4626-a978-ad60da2476ec';
    cy.apiSignInAsUser(userId, ServicesCounselorOfficeUserType);
  });

  it('Services Counselor can see error indicating missing Move Details information', () => {});

  it('Services Counselor can add an NTS shipment to the customer move', () => {
    // Services Counselor can enter NTS facility information on the Edit Shipment Details screen
  });

  it('Services Counselor can edit an existing NTS shipment request', () => {});
  it('Services Counselor can delete/remove an NTS shipment request', () => {});
  it('Services Counselor can enter accounting codes on the Orders Page', () => {});
  it('Services Counselor can assign accounting code(s) to a shipment', () => {});
  it('Services Counselor can edit existing accounting codes', () => {});
  it('Services Counselor cannot edit Customer Remarks field', () => {});

  it('Services Counselor can manually enter the shipment weight or reweigh weight on the Edit Shipment Details screen for NTS-Release', () => {});
});
