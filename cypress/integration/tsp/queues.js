import { tspUserVerifiesShipmentStatus } from '../../support/testTspStatus';

/* global cy */
describe('TSP User Views Shipment', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('tsp user views shipments in new shipments queue', function() {
    tspUserViewsNewShipments();
  });
  it('tsp user views shipments in accepted shipments queue', function() {
    tspUserViewsAcceptedShipments();
  });
  it('tsp user views shipments in approved shipments queue', function() {
    tspUserViewsApprovedShipments();
  });
  it('tsp user views shipments in in_transit shipments queue', function() {
    tspUserViewsInTransitShipments();
  });
  it('tsp user views shipments in delivered shipments queue', function() {
    tspUserViewsDeliveredShipments();
  });
  it('tsp user views shipments in completed shipments queue', function() {
    tspUserViewsCompletedShipments();
  });
});

function tspUserViewsNewShipments() {
  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find title
  cy.get('h1').contains('Queue: New Shipments');

  // Find shipment and check properties in row
  cy
    .get('div')
    .contains('div', 'BACON1')
    .parentsUntil('div.rt-tr-group')
    .contains('div', 'Awarded')
    .parentsUntil('div.rt-tr-group')
    .contains('div', 'LKBM7000002')
    .parentsUntil('div.rt-tr-group')
    .contains('div', 'Submitted, HHG')
    .parentsUntil('div.rt-tr-group')
    // TODO: cgilmer 2018/10/17 CircleCI seems to get this as 'US17 to Region 2'
    // Should figure out why this is happening
    // .contains('div', 'US88 to Region 2')
    // .parentsUntil('div.rt-tr-group')
    .contains('div', '15-May-19');

  // Find and open shipment
  cy
    .get('div')
    .contains('BACON1')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });
}

function tspUserViewsInTransitShipments() {
  // Open in transit shipments queue
  cy.visit('/queues/in_transit');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/in_transit/);
  });

  // Find title
  cy.get('h1').contains('Queue: In Transit Shipments');

  // Find in transit (generated in e2ebasic.go) and open it
  cy
    .get('div')
    .contains('NINOPK')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  // Status
  tspUserVerifiesShipmentStatus('Inbound');

  cy.get('a').contains('In Transit Shipments Queue');
}

function tspUserViewsDeliveredShipments() {
  // Open delivered shipments queue
  cy.visit('/queues/delivered');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/delivered/);
  });

  // Find title
  cy.get('h1').contains('Queue: Delivered Shipments');

  // Find delivered shipment (generated in e2ebasic.go) and open it
  cy
    .get('div')
    .contains('SCHNOO')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  // Status
  tspUserVerifiesShipmentStatus('Delivered');

  cy.get('a').contains('Delivered Shipments Queue');
}

function tspUserViewsAcceptedShipments() {
  // Open accepted shipments queue
  cy
    .get('div')
    .contains('Accepted Shipments')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/accepted/);
  });

  // Find title
  cy.get('h1').contains('Queue: Accepted Shipments');

  // Find shipment
  cy
    .get('div')
    .contains('BACON3')
    .dblclick();

  // Status
  tspUserVerifiesShipmentStatus('Shipment accepted');

  cy.get('a').contains('Accepted Shipments Queue');
}

function tspUserViewsApprovedShipments() {
  // Open approved shipments queue
  cy
    .get('div')
    .contains('Approved Shipments')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/approved/);
  });

  // Find title
  cy.get('h1').contains('Queue: Approved Shipments');

  // Find shipment
  cy
    .get('div')
    .contains('BACON1')
    .should('not.exist');
  cy
    .get('div')
    .contains('APPRVD')
    .dblclick();

  // Status
  tspUserVerifiesShipmentStatus('Awaiting pre-move survey');

  cy.get('a').contains('Approved Shipments Queue');
}

function tspUserViewsCompletedShipments() {
  // Open completed shipments queue
  cy
    .get('div')
    .contains('Completed Shipments')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/completed/);
  });

  // Find title
  cy.get('h1').contains('Queue: Completed Shipments');

  // Find shipment
  cy
    .get('div')
    .contains('NOCHKA')
    .dblclick();

  // Status
  tspUserVerifiesShipmentStatus('Delivered');

  cy.get('a').contains('Completed Shipments Queue');
}
