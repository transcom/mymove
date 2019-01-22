import { tspUserVerifiesShipmentStatus } from '../../support/testTspStatus';

/* global cy */
describe('TSP User Ships a Shipment', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('tsp user enters Pack and Pick Up shipment info', function() {
    tspUserEntersPackAndPickUpInfo();
    tspUserDeliversShipment();
  });

  it('tsp user enters a delivery date', function() {
    tspUserVisitsAnInTransitShipment('ENTDEL');
    tspUserVerifiesShipmentStatus('Inbound');
    tspUserCancelsEnteringADeliveryDate();
    tspUserEntersADeliveryDate();
    tspUserVerifiesShipmentStatus('Delivered');
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

function tspUserVisitsAnInTransitShipment(locator) {
  cy.patientVisit('/queues/in_transit');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/in_transit/);
  });

  cy
    .get('div')
    .contains(locator)
    .dblclick();
}

function tspUserEntersPackAndPickUpInfo() {
  cy.patientVisit('/queues/new');

  // Open approved shipments queue
  cy
    .get('div')
    .contains('Approved Shipments')
    .click();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/approved/);
  });

  // Find shipment and open it
  cy
    .get('div')
    .contains('CONGBL')
    .dblclick();

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
    .click();

  cy
    .get('div')
    .contains('11')
    .click();

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
    .click();

  cy
    .get('div.DayPicker-Month')
    .contains('10')
    .click();

  // Pick a Pickup date!
  cy
    .get('div')
    .contains('Actual pickup')
    .get('input[name="actual_pickup_date"]')
    .click();

  cy
    .get('div')
    .contains('11')
    .click();

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
  cy.get('div.actual_pack_date').contains('10');
  cy.get('div.actual_pickup_date').contains('11');

  // Appears in weights panel
  cy.get('.net_weight').should($div => {
    const text = $div.text();
    expect(text).to.include('Actual');
    expect(text).to.include('2,000 lbs');
  });

  // New status
  tspUserVerifiesShipmentStatus('Inbound');
}

function tspUserDeliversShipment() {
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
    .contains('13')
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
    .contains('13')
    .click();

  cy
    .get('button')
    .contains('Done')
    .click();

  // New status
  tspUserVerifiesShipmentStatus('Delivered');
}
