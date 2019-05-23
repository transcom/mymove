import { fillAndSaveStorageInTransit, editAndSaveStorageInTransit } from '../../support/testCreateStorageInTransit';

/* global cy */
describe('TSP user interacts with storage in transit panel', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('TSP user creates storage in transit request', function() {
    tspUserCreatesSitRequest();
  });
  it('TSP user starts and then cancels storage in transit request', function() {
    tspUserStartsAndCancelsSitRequest();
  });
  it('TSP user edits and saves storage in transit request', function() {
    tspUserEditsSitRequest();
  });
  it('TSP user starts and then cancels then completes Place into SIT form', function() {
    tspUserGoesToAcceptedSIT();
    tspUserStartsAndCancelsSitPlaceInSit();
    tspUserEntersInvalidActualStartDate();
    tspUserSubmitsPlaceInSit();
  });
  it('TSP user views remaining days and status of shipment in SIT (with frozen clock)', function() {
    tspUserEntitlementRemainingDays();
  });
  it('TSP user views remaining days and status of shipment expired in SIT (with frozen clock)', function() {
    tspUserEntitlementRemainingDaysExpired();
  });
});

// need to simulate a form submit

function tspUserCreatesSitRequest() {
  // Open accepted shipments queue
  cy.patientVisit('/queues/accepted');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/accepted/);
  });

  // Find shipment and open it
  cy.selectQueueItemMoveLocator('SITPAN');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  // Click on Request SIT and see SIT Request form
  cy
    .get('[data-cy=storage-in-transit-panel]')
    .should($div => {
      const text = $div.text();
      expect(text).to.include('Entitlement: 90 days');
      expect(text).to.not.include(' remaining)');
    })
    .get('[data-cy=storage-in-transit-panel] [data-cy=add-request]')
    .contains('Request SIT')
    .click()
    .get('.storage-in-transit-form')
    .should($div => {
      const text = $div.text();
      expect(text).to.include('SIT location');
      expect(text).to.include('Warehouse ID number');
      expect(text).to.include('Warehouse name');
      expect(text).to.include('Address line 1');
    });

  // fill out and submit the form
  fillAndSaveStorageInTransit();
}

function tspUserStartsAndCancelsSitRequest() {
  // Open accepted shipments queue
  cy.patientVisit('/queues/accepted');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/accepted/);
  });

  // Find shipment and open it
  cy.selectQueueItemMoveLocator('SITPAN');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  cy
    .get('[data-cy=storage-in-transit-panel] [data-cy=add-request]')
    .contains('Request SIT')
    .click();

  cy
    .get('.usa-button-secondary')
    .contains('Cancel')
    .click()
    .get('[data-cy=storage-in-transit-panel] [data-cy=add-request]')
    .should($div => {
      const text = $div.text();
      expect(text).to.not.include('Sit Location');
    });
}

function tspUserEditsSitRequest() {
  // Open accepted shipments queue
  cy.patientVisit('/queues/accepted');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/accepted/);
  });

  // Find shipment and open it
  cy.selectQueueItemMoveLocator('SITPAN');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  cy
    .get('[data-cy=storage-in-transit-panel] [data-cy=add-request]')
    .contains('Request SIT')
    .click();

  fillAndSaveStorageInTransit();

  // click edit
  cy
    .get('.sit-edit a')
    .first()
    .click();

  editAndSaveStorageInTransit();
}

function tspUserGoesToAcceptedSIT() {
  // Open in_transit shipments queue
  cy.patientVisit('/queues/in_transit');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/in_transit/);
  });

  // Find shipment and open it
  cy.selectQueueItemMoveLocator('SITAPR');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });
}

function tspUserStartsAndCancelsSitPlaceInSit() {
  // User starts from Accepted SIT
  cy
    .get('[data-cy=storage-in-transit-panel] [data-cy=place-in-sit-link]')
    .contains('Place into SIT')
    .click();

  cy.get('input[name=actual_start_date]').should('have.value', '3/26/2019');
  cy
    .get('[data-cy=place-into-sit-cancel]')
    .contains('Cancel')
    .click()
    .get('[data-cy=storage-in-transit-panel] [data-cy=place-in-sit-link]')
    .should($div => {
      const text = $div.text();
      expect(text).to.not.include('Actual start date');
    });
}

function tspUserEntersInvalidActualStartDate() {
  // User starts from Accepted SIT
  cy
    .get('[data-cy=storage-in-transit-panel] [data-cy=place-in-sit-link]')
    .contains('Place into SIT')
    .click();

  // calendar date picker date should be disabled
  cy
    .get('input[name=actual_start_date]')
    .should('have.value', '3/26/2019')
    // chooses invalid 3/22/2019
    .click()
    .get('.DayPickerInput-Overlay .DayPicker-Day')
    .contains('25')
    .should('have.class', 'DayPicker-Day--disabled');

  // submit should be disabled typed invalid date input
  cy
    .get('input[name=actual_start_date]')
    .should('have.value', '3/26/2019')
    .click()
    .clear()
    .type('3/25/2019');

  cy.get('[data-cy=storage-in-transit-panel]').click(); // click out of input field to hide datepicker
  cy.get('input[name=actual_start_date]').should('have.value', '3/25/2019');
  // expect submit to be disabled
  cy.get('[data-cy=place-in-sit-button]').should('be.disabled');
  cy
    .get('[data-cy=place-into-sit-cancel]')
    .contains('Cancel')
    .click();
}

function tspUserSubmitsPlaceInSit() {
  // User starts from Accepted SIT
  cy
    .get('[data-cy=storage-in-transit-panel] [data-cy=place-in-sit-link]')
    .contains('Place into SIT')
    .click();

  cy
    .get('input[name=actual_start_date')
    .should('have.value', '3/26/2019')
    // Chooses valid 3/30/2019
    .click()
    .get('.DayPickerInput-Overlay .DayPicker-Day')
    .contains('30')
    .click();
  cy.get('input[name=actual_start_date]').should('have.value', '3/30/2019');
  cy
    .get('[data-cy=place-in-sit-button]')
    .contains('Place Into SIT')
    .click()
    .get('[data-cy=storage-in-transit-panel]')
    .should($div => {
      const text = $div.text();
      expect(text).to.include('Entitlement: 90 days');
      expect(text).to.include(' remaining)');
      expect(text).to.include('Actual start date');
      expect(text).to.include('SIT Number');
      expect(text).to.include('Days used');
      expect(text).to.include('Expires');
    });
}

function tspUserGoesToPlacedSIT() {
  // Open in_transit shipments queue
  cy.patientVisit('/queues/in_transit');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/in_transit/);
  });

  // Find shipment and open it
  cy.selectQueueItemMoveLocator('SITIN1');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });
}

function tspUserEntitlementRemainingDays() {
  // Freeze the clock so we can test a specific remaining days.
  let now = new Date(Date.UTC(2019, 3, 10)).getTime(); // 4/10/2019
  cy.clock(now);

  tspUserGoesToPlacedSIT();

  cy
    .get('[data-cy=storage-in-transit-panel]')
    .should($div => {
      const text = $div.text();
      expect(text).to.include('In SIT');
      expect(text).to.include('Entitlement: 90 days (78 remaining)');
    })
    .get('[data-cy=storage-in-transit] [data-cy=sit-days-used]')
    .contains('12 days')
    .get('[data-cy=storage-in-transit] [data-cy=sit-expires]')
    .contains('28-Jun-2019');
}

function tspUserEntitlementRemainingDaysExpired() {
  // Freeze the clock so we can test a specific remaining days.
  let now = new Date(Date.UTC(2019, 6, 10)).getTime(); // 7/10/2019
  cy.clock(now);

  tspUserGoesToPlacedSIT();

  cy
    .get('[data-cy=storage-in-transit-panel]')
    .should($div => {
      const text = $div.text();
      expect(text).to.include('In SIT - SIT Expired');
      expect(text).to.include('Entitlement: 90 days (-13 remaining)');
    })
    .get('[data-cy=storage-in-transit] [data-cy=sit-days-used]')
    .contains('103 days')
    .get('[data-cy=storage-in-transit] [data-cy=sit-expires]')
    .contains('28-Jun-2019');
}
