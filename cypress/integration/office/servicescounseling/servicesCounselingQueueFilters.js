import { ServicesCounselorOfficeUserType } from '../../../support/constants';

describe('Services counselor user', () => {
  before(() => {
    cy.prepareOfficeApp();
  });

  beforeEach(() => {
    cy.intercept('**/ghc/v1/swagger.yaml').as('getGHCClient');
    cy.intercept('**/ghc/v1/queues/counseling?page=1&perPage=20&sort=submittedAt&order=asc&needsPPMCloseout=false').as(
      'getSortedMoves',
    );
    cy.intercept(
      '**/ghc/v1/queues/counseling?page=1&perPage=20&sort=closeoutInitiated&order=asc&needsPPMCloseout=true&ppmType=**',
    ).as('getPPMTypeFilteredMoves');
    cy.intercept(
      '**/ghc/v1/queues/counseling?page=1&perPage=20&sort=closeoutInitiated&order=asc&needsPPMCloseout=true&closeoutInitiated=**',
    ).as('getCloseoutInitiatedFilteredMoves');
    cy.intercept(
      '**/ghc/v1/queues/counseling?page=1&perPage=20&sort=closeoutInitiated&order=asc&needsPPMCloseout=true&closeoutLocation=**',
    ).as('getCloseoutLocationFilteredMoves');
    cy.intercept(
      '**/ghc/v1/queues/counseling?page=1&perPage=20&sort=closeoutInitiated&order=asc&destinationDutyLocation=**&needsPPMCloseout=true',
    ).as('getDestinationDutyLocationFilteredMoves');
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
    cy.get('[data-testid="closeout-tab-link"]').click();
    cy.wait(['@getSortedMoves']);
  });

  it('is able to filter moves based on ppm type', () => {
    cy.get('th[data-testid="ppmType"] > div > select').select('Full');
    cy.wait(['@getPPMTypeFilteredMoves']);
    cy.get('h1').contains('Moves (0)');
  });
  it('is able to filter moves based on Closeout initiated', () => {
    cy.get('th[data-testid="closeoutInitiated"] > div > div > input').type('01 Dec 2020');
    cy.wait(['@getCloseoutInitiatedFilteredMoves']);
    cy.get('h1').contains('Moves (0)');
  });
  it('is able to filter moves based on Closeout location', () => {
    cy.get('th[data-testid="closeoutLocation"] > div > input').type('j').blur();
    cy.wait(['@getCloseoutLocationFilteredMoves']);
    cy.get('h1').contains('Moves (0)');
  });
  it('is able to filter moves based on destination duty location', () => {
    cy.get('th[data-testid="destinationDutyLocation"] > div > input').type('Fort').blur();
    cy.wait(['@getDestinationDutyLocationFilteredMoves']);
    cy.get('h1').contains('Moves (4)');
  });
});
