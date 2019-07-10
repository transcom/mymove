import { tspUserVerifiesShipmentStatus } from '../../support/testTspStatus';
import { nthDayOfCurrentMonth } from '../../support/utils';

/* global cy */
describe('TSP User Ships a Shipment', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('tsp user enters Pack and Pick Up shipment info', function() {
    cy.nextAvailable(nthDayOfCurrentMonth(1)).then(availableDays => {
      tspUserEntersPackAndPickUpInfo();
      tspUserDeliversShipment(availableDays);
    });
  });

  it('tsp user enters a delivery date', function() {
    tspUserVisitsAnInTransitShipment('ENTDEL');
    tspUserVerifiesShipmentStatus('Inbound');
    tspUserCancelsEnteringADeliveryDate();
    tspUserEntersADeliveryDate();
    tspUserVerifiesShipmentStatus('Delivered');
    tspUserVerifiesSITStatus();
  });
});

function tspUserCancelsEnteringADeliveryDate() {
  cy
    .get('button')
    .contains('Enter Delivery')
    .should('exist');

  cy
    .get('button')
    .contains('Enter Delivery')
    .click();

  cy
    .get('input[name="actual_delivery_date"]')
    .type('10/24/2018')
    .blur();

  cy
    .get('a')
    .contains('Cancel')
    .click();

  cy.get('input[name="actual_delivery_date"]').should('not.exist');
}

function tspUserEntersADeliveryDate() {
  cy
    .get('button')
    .contains('Enter Delivery')
    .click();

  cy.get('input[name="actual_delivery_date"]').should('be.empty');

  cy
    .get('input[name="actual_delivery_date"]')
    .type('10/24/2018')
    .blur();

  cy
    .get('button')
    .contains('Done')
    .click();

  cy.get('button').should('not.contain', 'Enter Delivery');
}

function tspUserVerifiesSITStatus() {
  // IN_SIT Destination Sits are in delivered status and have an out date
  cy.get('[data-cy=storage-in-transit-panel]').contains('SIT Delivered');

  cy.get('[data-cy=storage-in-transit-panel] [data-cy=sit-dates]').contains('Date out');
}

function tspUserVisitsAnInTransitShipment(locator) {
  cy.patientVisit('/queues/in_transit');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/in_transit/);
  });

  cy.selectQueueItemMoveLocator(locator);
}

function tspUserEntersPackAndPickUpInfo() {
  cy.patientVisit('/queues/new');

  // Open approved shipments queue
  cy
    .get('div')
    .contains('All Shipments')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/all/);
  });

  // Find shipment and open it
  cy.selectQueueItemMoveLocator('CONGBL');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  // Click the Enter Pickup button
  cy
    .get('div')
    .contains('Enter Pickup')
    .click();

  // Done button should be disabled.
  cy
    .get('button')
    .contains('Done')
    .should('be.disabled');

  // Pick a date!
  cy
    .get('div')
    .contains('Actual pickup')
    .get('input[name="actual_pickup_date"]')
    .type('10/20/2018')
    .blur();

  // Cancel
  cy
    .get('a')
    .contains('Cancel')
    .click();

  // Wash, Rinse, Repeat
  // Click the Enter Pickup button
  cy
    .get('div')
    .contains('Enter Pickup')
    .click();

  // Done button should be disabled.
  cy
    .get('button')
    .contains('Done')
    .should('be.disabled');

  // Pick a Pack date!
  cy
    .get('div')
    .contains('Actual packing (first day)')
    .get('input[name="actual_pack_date"]')
    .type('10/18/2018')
    .blur();

  // Pick a Pickup date!
  cy
    .get('div')
    .contains('Actual pickup')
    .get('input[name="actual_pickup_date"]')
    .type('10/19/2018')
    .blur();

  // Done button should STILL be disabled.
  cy
    .get('button')
    .contains('Done')
    .should('be.disabled');
  cy
    .get('input[name="net_weight"]')
    .first()
    .clear()
    .type('2000')
    .blur();
  // Done button should be enabled.
  cy
    .get('button')
    .contains('Done')
    .should('be.enabled');
  cy
    .get('input[name="gross_weight"]')
    .first()
    .clear()
    .type('3000')
    .blur();
  cy
    .get('input[name="tare_weight"]')
    .first()
    .clear()
    .type('1000')
    .blur();
  cy
    .get('button')
    .contains('Done')
    .click();

  // Appears in dates panel
  cy.get('div.actual_pack_date').contains('18');
  cy.get('div.actual_pickup_date').contains('19');

  // Appears in weights panel
  cy.get('.net_weight').should($div => {
    const text = $div.text();
    expect(text).to.include('Actual');
    expect(text).to.include('2,000 lbs');
  });

  // New status
  tspUserVerifiesShipmentStatus('Inbound');
}

function tspUserDeliversShipment(availableDays) {
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  // Click the Transport button
  cy
    .get('div')
    .contains('Enter Delivery')
    .click();

  // Done button should be disabled.
  cy
    .get('button')
    .contains('Done')
    .should('be.disabled');

  // Pick a date!
  cy
    .get('div')
    .contains('Actual Delivery Date')
    .get('input')
    .click();

  cy
    .get('div')
    .contains(availableDays[0].format('DD'))
    .click();

  // Cancel
  cy
    .get('a')
    .contains('Cancel')
    .click();

  // Wash, Rinse, Repeat
  // Click the Transport button
  cy
    .get('div')
    .contains('Enter Delivery')
    .click();

  // Done button should be disabled.
  cy
    .get('button')
    .contains('Done')
    .should('be.disabled');

  // Pick a date!
  cy
    .get('div')
    .contains('Actual Delivery Date')
    .get('input')
    .click();

  cy
    .get('div')
    .contains(availableDays[1].format('DD'))
    .click();

  cy
    .get('button')
    .contains('Done')
    .click();

  // New status
  tspUserVerifiesShipmentStatus('Delivered');
}
