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
      '**/ghc/v1/queues/counseling?page=1&perPage=20&sort=closeoutInitiated&order=asc&locator=**&needsPPMCloseout=true',
    ).as('getLocatorFilteredMoves');
    cy.intercept(
      '**/ghc/v1/queues/counseling?page=1&perPage=20&sort=closeoutInitiated&order=asc&locator=**&needsPPMCloseout=true&ppmType=**',
    ).as('getPPMTypeAndLocatorFilteredMoves');
    cy.intercept(
      '**/ghc/v1/queues/counseling?page=1&perPage=20&sort=closeoutInitiated&order=asc&needsPPMCloseout=true&closeoutInitiated=**',
    ).as('getCloseoutInitiatedFilteredMoves');
    cy.intercept(
      '**/ghc/v1/queues/counseling?page=1&perPage=20&sort=closeoutInitiated&order=asc&locator=**&needsPPMCloseout=true&closeoutLocation=**',
    ).as('getCloseoutLocationFilteredMoves');
    cy.intercept(
      '**/ghc/v1/queues/counseling?page=1&perPage=20&sort=closeoutInitiated&order=asc&locator=**&destinationDutyLocation=**&needsPPMCloseout=true',
    ).as('getDestinationDutyLocationAndLocatorFilteredMoves');
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
    // PPMSC1 is a Partial PPM move, so when we search for Partial, we should see it in the results
    cy.get('th[data-testid="locator"] > div > input').type('PPMSC1').blur();
    cy.wait(['@getLocatorFilteredMoves']);
    cy.get('th[data-testid="ppmType"] > div > select').select('Partial');
    cy.wait(['@getPPMTypeAndLocatorFilteredMoves']);
    cy.get('td').contains('PPMSC1');

    // When we search for Full PPM moves, PPMSC1 should not come up
    cy.get('th[data-testid="ppmType"] > div > select').select('Full');
    cy.wait(['@getPPMTypeAndLocatorFilteredMoves']);
    cy.get('h1').contains('Moves (0)');
  });

  it('is able to filter moves based on Closeout initiated', () => {
    cy.get('th[data-testid="closeoutInitiated"] > div > div > input').type('01 Dec 2020');
    cy.wait(['@getCloseoutInitiatedFilteredMoves']);
    cy.get('h1').contains('Moves (0)');
  });

  it('is able to filter moves based on Closeout location', () => {
    // add filter for move code CLSOFF (which has a closeout office set to JPPSO Testy McTest)
    cy.get('th[data-testid="locator"] > div > input').type('CLSOFF').blur();
    // add another filter for the closeout office column
    cy.get('th[data-testid="closeoutLocation"] > div > input').type('jppso testy').blur();
    cy.wait(['@getCloseoutLocationFilteredMoves']);
    cy.get('td').contains('CLSOFF');
    // Add some nonsense text to our filter
    cy.get('th[data-testid="closeoutLocation"] > div > input').type('z').blur();
    cy.wait(['@getCloseoutLocationFilteredMoves']);
    // now we should get no results
    cy.get('h1').contains('Moves (0)');
  });

  it('is able to filter moves based on destination duty location', () => {
    // add filter for move code PPMSC1 (which has Fort Gordon as its destination duty location)
    cy.get('th[data-testid="locator"] > div > input').type('PPMSC1').blur();
    cy.wait(['@getLocatorFilteredMoves']);
    // Add destination duty location filter 'fort'
    cy.get('th[data-testid="destinationDutyLocation"] > div > input').type('fort').blur();
    cy.wait(['@getDestinationDutyLocationAndLocatorFilteredMoves']);
    // We should still see our move
    cy.get('td').contains('PPMSC1');
    // Add nonsense string to our filter (so now we're searching for 'fortzzzz')
    cy.get('th[data-testid="destinationDutyLocation"] > div > input').type('zzzz').blur();
    cy.wait(['@getDestinationDutyLocationAndLocatorFilteredMoves']);
    // Now we shouldn't see any results
    cy.get('h1').contains('Moves (0)');
  });
});
