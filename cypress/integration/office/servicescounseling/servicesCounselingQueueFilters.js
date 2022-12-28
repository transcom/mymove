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

  it('is able to filter partial vs full moves based on ppm type', () => {
    // PPMCHE is a Partial PPM move, so when we search for Partial, we should see it in the results
    cy.get('th[data-testid="locator"] > div > input').type('PPMCHE').blur();
    cy.wait(['@getLocatorFilteredMoves']);
    cy.get('th[data-testid="ppmType"] > div > select').select('Partial');
    cy.wait(['@getPPMTypeAndLocatorFilteredMoves']);
    cy.get('td').contains('PPMCHE');

    // When we search for Full PPM moves, PPMCHE should not come up
    cy.get('th[data-testid="ppmType"] > div > select').select('Full');
    cy.wait(['@getPPMTypeAndLocatorFilteredMoves']);
    cy.get('h1').contains('Moves (0)');
  });

  it('is able to filter moves based on PPM Closeout initiated', () => {
    cy.get('th[data-testid="closeoutInitiated"] > div > div > input').type('11 Dec 2020');
    cy.wait(['@getCloseoutInitiatedFilteredMoves']);
    cy.get('h1').contains('Moves (0)');
  });

  it('is able to filter moves based on PPM Closeout location', () => {
    // add filter for move code CLSOFF (which routes to PPM closeout tab and has a closeout office set to Los Angeles AFB)
    cy.get('th[data-testid="locator"] > div > input').type('CLSOFF').blur();
    // add another filter for the closeout office column checking it's not case sensitive
    cy.get('th[data-testid="closeoutLocation"] > div > input').type('LOS ANGELES').blur();
    cy.wait(['@getCloseoutLocationFilteredMoves']);
    cy.get('td').contains('CLSOFF');
    // Add some nonsense z text to our filter
    cy.get('th[data-testid="closeoutLocation"] > div > input').type('z').blur();
    cy.wait(['@getCloseoutLocationFilteredMoves']);
    // now we should get no results
    cy.get('h1').contains('Moves (0)');
  });

  it('is able to filter moves based on destination duty location', () => {
    // add filter for move code PPMCHE (PPM closeout that has Fort Gordon as its destination duty location)
    cy.get('th[data-testid="locator"] > div > input').type('PPMCHE').blur();
    cy.wait(['@getLocatorFilteredMoves']);
    // Add destination duty location filter 'fort'
    cy.get('th[data-testid="destinationDutyLocation"] > div > input').type('fort').blur();
    cy.wait(['@getDestinationDutyLocationAndLocatorFilteredMoves']);
    // We should still see our move
    cy.get('td').contains('PPMCHE');
    // Add nonsense string to our filter (so now we're searching for 'fortzzzz')
    cy.get('th[data-testid="destinationDutyLocation"] > div > input').type('zzzz').blur();
    cy.wait(['@getDestinationDutyLocationAndLocatorFilteredMoves']);
    // Now we shouldn't see any results
    cy.get('h1').contains('Moves (0)');
  });
});
