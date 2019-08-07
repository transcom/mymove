import { fillAndSaveStorageInTransit, editAndSaveStorageInTransit } from '../../support/testCreateStorageInTransit';

/* global cy */
describe('TSP user interacts with storage in transit panel', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('TSP user Delivers Shipment with two (2) SITs 30 mi or less', function() {
    tspUserDeliversShipmentSIT030();
  });
  it('TSP user Delivers Shipment with one (1) SIT 50 mi or less', function() {
    tspUserDeliversShipmentSIT050();
  });
  it('TSP user Delivers Shipment with one (1) SIT more than 50 mi', function() {
    tspUserDeliversShipmentSIT051();
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
    tspUserGoesToApprovedSit();
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
  it('TSP user views denied SIT', function() {
    tspUserViewsDeniedSit();
  });
  it('TSP user releases SIT IN-SIT at ORIGIN', function() {
    tspUserReleasesOriginSit();
  });
  it('TSP user cancels delete, then actually deletes SIT', function() {
    tspUserDeletesSitRequest();
  });
  it('TSP user edits IN-SIT SIT request', function() {
    tspUserEditsSitRequestInSit();
  });
  it('TSP user edits RELEASED SIT request', function() {
    tspUserEditsReleasedSitRequest();
  });
  it('TSP user edits DELIVERED SIT request', function() {
    tspUserEditsDeliveredSitRequest();
  });
});

function tspUserDeliversShipmentSIT030() {
  // SIT030
  // Origin SIT
  // Destination SIT

  // Release Origin SIT before entering delivery
  let tspQueue = 'in_transit';
  let moveLocator = 'SIT030';
  let releaseOnDate = '5/22/2019';
  let dateOut = '22-May-2019';
  tspUserSubmitsReleaseSit(tspQueue, moveLocator, releaseOnDate, dateOut);

  // Enter delivery for shipment
  let deliveryDate = '5/28/2019';
  tspUserDeliversShipment(tspQueue, moveLocator, deliveryDate);

  // Verify Destination SIT was updated on panel
  dateOut = '28-May-2019';
  tspUserVerifySitReleasedDelivered('Destination', 'Delivered', dateOut);

  // Verify Invoice contains expected 210* line item(s)
  // 210A Origin SIT Line Item
  tspUserInvoiceContains('210ASIT Pup/Del - 30 or Less MilesO9 mi');
  // 210A Destination SIT Line Item
  tspUserInvoiceContains('210ASIT Pup/Del - 30 or Less MilesD13 mi');
  cy.get('[data-cy=invoice-panel] [data-cy=basic-panel-content] tbody > tr').should('not.contain', /^210B\w+/);
  cy.get('[data-cy=invoice-panel] [data-cy=basic-panel-content] tbody > tr').should('not.contain', /^210C\w+/);
}
function tspUserDeliversShipmentSIT050() {
  // SIT050
  // Origin SIT

  // Release Origin SIT before entering delivery
  let tspQueue = 'in_transit';
  let moveLocator = 'SIT050';
  let releaseOnDate = '5/22/2019';
  let dateOut = '22-May-2019';
  tspUserSubmitsReleaseSit(tspQueue, moveLocator, releaseOnDate, dateOut);

  // Enter delivery for shipment
  let deliveryDate = '5/28/2019';
  tspUserDeliversShipment(tspQueue, moveLocator, deliveryDate);

  // Verify Invoice contains expected 210* line item(s)
  // 210A Origin SIT Line Item
  tspUserInvoiceContains('210ASIT Pup/Del - 30 or Less MilesO43 mi');
  // 210B Origin SIT Line Item
  tspUserInvoiceContains('210BSIT Pup/Del 31 - 50 MilesO43 mi');
  cy.get('[data-cy=invoice-panel] [data-cy=basic-panel-content] tbody > tr').should('not.contain', /^210C\w+/);
}
function tspUserDeliversShipmentSIT051() {
  // SIT051
  // Destination SIT

  // Enter delivery for shipment
  let tspQueue = 'in_transit';
  let moveLocator = 'SIT051';
  let deliveryDate = '5/28/2019';
  tspUserDeliversShipment(tspQueue, moveLocator, deliveryDate);

  // Verify Destination SIT was updated on panel
  let dateOut = '28-May-2019';
  tspUserVerifySitReleasedDelivered('Destination', 'Delivered', dateOut);

  // Verify Invoice contains expected 210* line item(s)
  // 210C Destination SIT Line Item
  tspUserInvoiceContains('210CSIT Pup/Del Over 50 MilesD226 mi');
  cy.get('[data-cy=invoice-panel] [data-cy=basic-panel-content] tbody > tr').should('not.contain', /^210A\w+/);
  cy.get('[data-cy=invoice-panel] [data-cy=basic-panel-content] tbody > tr').should('not.contain', /^210B\w+/);
}

function tspUserReleasesOriginSit() {
  // SITOIN - SIT added to shipment in transit, ready to be
  // placed inSIT
  let tspQueue = 'in_transit';
  let moveLocator = 'SITOIN';
  let releaseOnDate = '5/26/2019';
  let dateOut = '26-May-2019';
  tspUserSubmitsReleaseSit(tspQueue, moveLocator, releaseOnDate, dateOut);

  // DISIT1 - Origin SIT added after Shipment is Delivered
  tspQueue = 'delivered';
  moveLocator = 'DISIT1';
  releaseOnDate = '5/26/2019';
  dateOut = '26-May-2019';
  tspUserSubmitsReleaseSit(tspQueue, moveLocator, releaseOnDate, dateOut);

  // DISIT2 - Origin SIT added to Shipment in Transit and then Shipment is Delivered
  tspQueue = 'delivered';
  moveLocator = 'DISIT2';
  releaseOnDate = '5/26/2019';
  dateOut = '26-May-2019';
  tspUserSubmitsReleaseSit(tspQueue, moveLocator, releaseOnDate, dateOut);
}

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
      expect(text).to.include('SIT Location');
      expect(text).to.include('Warehouse ID number');
      expect(text).to.include('Warehouse name');
      expect(text).to.include('Address line 1');
    });

  // fill out and submit the form
  fillAndSaveStorageInTransit();

  // Verify action links
  cy
    .get('[data-cy=storage-in-transit-panel] [data-cy=sit-edit-link]')
    .contains('Edit')
    .get('[data-cy=storage-in-transit-panel] [data-cy=sit-delete-link]')
    .contains('Delete');
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
    .get('[data-cy="sit-edit-link"]')
    .first()
    .click();

  editAndSaveStorageInTransit();
}

function tspUserGoesToApprovedSit() {
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
  // User starts from Approved SIT

  // Verify action links
  cy.get('[data-cy=storage-in-transit-panel] [data-cy=sit-delete-link]').contains('Delete');

  cy
    .get('[data-cy=storage-in-transit-panel] [data-cy=sit-place-into-sit-link]')
    .contains('Place into SIT')
    .click();

  cy.get('input[name=actual_start_date]').should('have.value', '3/26/2019');
  cy
    .get('[data-cy=place-into-sit-cancel]')
    .contains('Cancel')
    .click()
    .get('[data-cy=storage-in-transit-panel] [data-cy=sit-place-into-sit-link]')
    .should($div => {
      const text = $div.text();
      expect(text).to.not.include('Actual start date');
    });
}

function tspUserEntersInvalidActualStartDate() {
  // User starts from Approved SIT
  cy
    .get('[data-cy=storage-in-transit-panel] [data-cy=sit-place-into-sit-link]')
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
  // User starts from Approved SIT
  cy
    .get('[data-cy=storage-in-transit-panel] [data-cy=sit-place-into-sit-link]')
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

  // Verify action links
  cy
    .get('[data-cy=storage-in-transit-panel] [data-cy=sit-edit-link]')
    .contains('Edit')
    .get('[data-cy=storage-in-transit-panel] [data-cy=sit-delete-link]')
    .contains('Delete');
}

function tspUserGoesToPlacedSit() {
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

  tspUserGoesToPlacedSit();

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

  tspUserGoesToPlacedSit();

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

function tspUserViewsDeniedSit() {
  // Open in_transit shipments queue
  cy.patientVisit('/queues/in_transit');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/in_transit/);
  });

  // Find shipment and open it
  cy.selectQueueItemMoveLocator('SITDN2');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  // Verify action links
  cy.get('[data-cy=storage-in-transit-panel] [data-cy=sit-delete-link]').contains('Delete');
}

function tspUserInvoiceContains(lineItem) {
  // The invoice table should display the unbilled line items.
  cy.get('[data-cy=invoice-panel] [data-cy=basic-panel-content] tbody').contains('tr', lineItem);
}

// tspUserVerifySitReleasedDelivered
// params: sitLocation is 'Origin' or 'Destination'
//         status is 'Released' or 'Delivered'
//         dateOut is in the formation dd-mmm-yyyy
function tspUserVerifySitReleasedDelivered(sitLocation, status, dateOut) {
  cy.get('[data-cy=storage-in-transit-panel]').should($div => {
    const text = $div.text();
    expect(text).to.include(sitLocation + ' SIT');
    expect(text).to.include(status);
    expect(text).to.include('Entitlement: 90 days');
    expect(text).to.include('Actual start date');
    expect(text).to.include('SIT Number');
    expect(text).to.include('Days used');
    expect(text).to.include('Expires');
    expect(text).to.include('Date out' + dateOut);
  });
}

function tspUserVisitsQueue(tspQueue) {
  cy.patientVisit('/queues/' + tspQueue);
  let re;
  switch (tspQueue) {
    case 'in_transit':
      re = new RegExp('^/queues/in_transit');
      break;
    case 'delivered':
      re = new RegExp('^/queues/delivered');
      break;
    case 'accepted':
      re = new RegExp('^/queues/accepted');
      break;
    case 'new':
      re = new RegExp('^/queues/new');
      break;
    case 'all':
      re = new RegExp('^/queues/all');
      break;
    default:
      re = new RegExp('^/queues/undefined');
  }
  cy.location().should(loc => {
    expect(loc.pathname).to.match(re);
  });
}

function tspUserDeliversShipment(tspQueue, moveLocator, deliveryDate) {
  // Open TSP shipments queue
  tspUserVisitsQueue(tspQueue);

  // Find shipment with moveLocator
  cy.selectQueueItemMoveLocator(moveLocator);

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  // click Enter Delivery button
  cy
    .get('[data-cy=tsp-enter-delivery]')
    .contains('Enter Delivery')
    .click();

  // enter Actual Delivery Date
  cy
    .get('input[name=actual_delivery_date]')
    .type(deliveryDate)
    .blur();

  // press button to release shipment and confirm the expected information on the panel
  // after releasing the shipment
  cy
    .get('[data-cy=tsp-enter-delivery-submit]')
    .contains('Done')
    .click();
}

function tspUserSubmitsReleaseSit(tspQueue, moveLocator, releaseOnDate, dateOut) {
  // Open TSP shipments queue
  tspUserVisitsQueue(tspQueue);

  // Find shipment with moveLocator
  cy.selectQueueItemMoveLocator(moveLocator);

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  // click release shipment link
  cy
    .get('[data-cy=storage-in-transit-panel] [data-cy=sit-release-from-sit-link]')
    .contains('Release from SIT')
    .click();

  // Test canceling after clicking release shipment link
  cy
    .get('[data-cy=release-from-sit-cancel]')
    .contains('Cancel')
    .click();

  cy.get('[data-cy=storage-in-transit-panel] [data-cy=sit-release-from-sit-link]').should($div => {
    const text = $div.text();
    expect(text).to.not.include('Date out');
  });

  // click release shipment link
  cy
    .get('[data-cy=storage-in-transit-panel] [data-cy=sit-release-from-sit-link]')
    .contains('Release from SIT')
    .click();

  // enter in date released on
  cy
    .get('input[name=released_on]')
    .focus()
    .type(releaseOnDate)
    .blur();

  // press button to release shipment and confirm the expected information on the panel
  // after releasing the shipment
  cy
    .get('[data-cy=release-from-sit-button]')
    .contains('Done')
    .click();

  tspUserVerifySitReleasedDelivered('Origin', 'Released', dateOut);
}

function tspUserDeletesSitRequest() {
  // Open accepted shipments queue
  cy.patientVisit('/queues/accepted');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/accepted/);
  });

  // Find shipment and open it
  cy.selectQueueItemMoveLocator('SITDEL');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  // Click on Delete SIT and see SIT Delete warning, then click cancel and it should go away.
  cy
    .get('[data-cy=storage-in-transit-panel] [data-cy=sit-delete-link]')
    .click()
    .get('[data-cy=sit-delete-warning] [data-cy=sit-delete-cancel]')
    .click()
    .get('[data-cy=sit-delete-warning]')
    .should('not.exist');

  // Now click on Delete SIT again, then actually delete it this time.
  cy
    .get('[data-cy=storage-in-transit-panel] [data-cy=sit-delete-link]')
    .click()
    .get('[data-cy=sit-delete-warning] [data-cy=sit-delete-delete]')
    .click()
    .get('[data-cy=storage-in-transit]')
    .should('not.exist');
}

function tspUserEditsSitRequestInSit() {
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

  cy
    .get('[data-cy=storage-in-transit-panel] [data-cy=sit-edit-link]')
    .contains('Edit')
    .click();
  cy
    .get('input[name=actual_start_date]')
    .should('have.value', '3/30/2019')
    .click();
  cy
    .get('.DayPickerInput-Overlay .DayPicker-Day')
    .contains('29')
    .click();
  cy.get('input[name=actual_start_date]').should('have.value', '3/29/2019');
  cy.get('.usa-button-primary').click();
  cy.get('[data-cy=storage-in-transit-panel]').should($div => {
    const text = $div.text();
    expect(text).to.include('29-Mar-2019');
  });
}

function tspUserEditsReleasedSitRequest() {
  // Open in transit shipments queue
  cy.patientVisit('/queues/in_transit');

  //
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/in_transit/);
  });

  // Find shipment that is inSIT at ORIGIN and open it
  cy.selectQueueItemMoveLocator('SITOIN');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  cy
    .get('[data-cy=storage-in-transit-panel] [data-cy=sit-edit-link]')
    .contains('Edit')
    .click();
  cy
    .get('input[name=out_date]')
    .should('have.value', '5/26/2019')
    .click();
  cy
    .get('.DayPickerInput-Overlay .DayPicker-Day')
    .contains('29')
    .click();
  cy.get('input[name=out_date]').should('have.value', '5/29/2019');
  cy.get('.usa-button-primary').click();
  cy.get('[data-cy=storage-in-transit-panel]').should($div => {
    const text = $div.text();
    expect(text).to.include('29-May-2019');
  });
}

function tspUserEditsDeliveredSitRequest() {
  cy.patientVisit('/queues/delivered');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/delivered/);
  });

  // Find shipment and open it
  cy.selectQueueItemMoveLocator('SITDST');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  cy
    .get('[data-cy=storage-in-transit-panel] [data-cy=sit-edit-link]')
    .contains('Edit')
    .click();
  cy
    .get('input[name=out_date]')
    .should('have.value', '3/27/2019')
    .click();
  cy
    .get('.DayPickerInput-Overlay .DayPicker-Day')
    .contains('22')
    .click();
  // cannot change date_out to a date before the actual start date
  cy
    .get('input[name=out_date]')
    .should('have.value', '3/27/2019')
    .click();
  cy
    .get('.DayPickerInput-Overlay .DayPicker-Day')
    .contains('29')
    .click();
  cy.get('input[name=out_date]').should('have.value', '3/29/2019');
  cy.get('.usa-button-primary').click();
  cy.get('[data-cy=storage-in-transit-panel]').should($div => {
    const text = $div.text();
    expect(text).to.include('29-Mar-2019');
  });
  cy
    .get('.editable-panel-header')
    .contains('Dates')
    .siblings()
    .click()
    .get('input[name="dates.actual_delivery_date"]')
    //When SIT date out is changed, actual delivery date in the dates panel is also changed
    .should('have.value', '3/29/2019')
    .click()
    .get('.DayPickerInput-Overlay .DayPicker-Day')
    .contains('27')
    .click()
    .get('input[name="dates.actual_delivery_date"]')
    .should('have.value', '3/27/2019');

  cy
    .get('[data-cy=storage-in-transit-panel] [data-cy=sit-edit-link]')
    .contains('Edit')
    .click();
  cy
    //When shipment actual delivery date is changed, SIT out date is unchanged
    .get('input[name=out_date]')
    .should('have.value', '3/29/2019');
}
