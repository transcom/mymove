import { testPremoveSurvey } from '../../support/testPremoveSurvey';

/* global cy */
describe('TSP User Ships a Shipment', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('tsp user Picks Up a shipment', function() {
    tspUserPacksShipment();
    tspUserPicksUpShipment();
    tspUserDeliversShipment();
  });

  it('tsp user uploads origin documents', function() {
    tspUserVisitsInTransitQueue();
    tspUserViewsInTransitShipment();
    tspUserFillsInDatesPanel();
    tspUserFillsInWeightsPanel();
    tspUserClicksUploadOriginDocs('7e2f1b5f-74f8-47b8-ae49-63b0e58a70bd');
  });
});

function tspUserPacksShipment(locator) {
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
    .contains('SHIPME')
    .dblclick();

  // Click the Pack button
  cy
    .get('div')
    .contains('Enter Packing')
    .click();

  // Done button should be disabled.
  cy
    .get('button')
    .contains('Done')
    .should('be.disabled');

  // Pick a date!
  cy
    .get('div')
    .contains('Actual Pack Date')
    .get('input')
    .click();

  cy
    .get('div.DayPicker-Month')
    .contains('10')
    .click();

  // Cancel
  cy
    .get('button')
    .contains('Cancel')
    .click();

  // Check that the date doesn't appear in dates panel
  cy.get('div.actual_pack_date').contains('missing');

  // Wash, Rinse, Repeat
  // Click the Pack button
  cy
    .get('div')
    .contains('Enter Packing')
    .click();

  // Done button should be disabled.
  cy
    .get('button')
    .contains('Done')
    .should('be.disabled');

  // Pick a date!
  cy
    .get('div')
    .contains('Actual Pack Date')
    .get('input')
    .click();

  cy
    .get('div.DayPicker-Month')
    .contains('10')
    .click();

  cy
    .get('button')
    .contains('Done')
    .click();

  // Appears in dates panel
  cy.get('div.actual_pack_date').contains('10');
}

function tspUserPicksUpShipment() {
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  // Click the Transport button
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
    .contains('Actual Pickup Date')
    .get('input')
    .click();

  cy
    .get('div')
    .contains('11')
    .click();

  // Cancel
  cy
    .get('button')
    .contains('Cancel')
    .click();

  // Wash, Rinse, Repeat
  // Click the Transport button
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
    .contains('Actual Pickup Date')
    .get('input')
    .click();

  cy
    .get('div')
    .contains('11')
    .click();

  cy
    .get('button')
    .contains('Done')
    .click();

  // Appears in dates panel
  cy.get('div.actual_pickup_date').contains('11');

  // New status
  cy.get('li').contains('In_transit');
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
    .get('button')
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
  cy.get('li').contains('Delivered');
}

function tspUserVisitsInTransitQueue() {
  cy.visit('/queues/in_transit');
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/in_transit/);
  });
}

function tspUserViewsInTransitShipment() {
  cy
    .get('div')
    .contains('ORIDOC')
    .dblclick();
}

function tspUserFillsInWeightsPanel() {
  cy
    .get('.editable-panel-header')
    .contains('Weights & Items')
    .siblings()
    .click();

  cy
    .get('input[name="weights.gross_weight"]')
    .type('5000')
    .blur();

  cy
    .get('input[name="weights.tare_weight"]')
    .type('2500')
    .blur();

  cy
    .get('input[name="weights.net_weight"]')
    .type('2500')
    .blur();

  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save')
    .click();
}

function tspUserFillsInDatesPanel() {
  // Conducted Date
  cy
    .get('.editable-panel-header')
    .contains('Dates')
    .siblings()
    .click();

  cy
    .get('input[name="dates.pm_survey_conducted_date"]')
    .first()
    .type('7/20/2018')
    .blur();
  cy.get('select[name="dates.pm_survey_method"]').select('PHONE');

  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save')
    .click();

  cy.reload();

  cy.get('div.pm_survey_conducted_date').contains('20-Jul-18');
  cy.get('div.pm_survey_method').contains('Phone');

  // Pack Dates
  cy
    .get('.editable-panel-header')
    .contains('Dates')
    .siblings()
    .click();

  cy
    .get('input[name="dates.pm_survey_planned_pack_date"]')
    .first()
    .type('8/1/2018')
    .blur();
  cy
    .get('input[name="dates.actual_pack_date"]')
    .first()
    .type('8/2/2018')
    .blur();

  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save')
    .click();

  cy.reload();

  cy.get('div.original_pack_date').contains('15-May-19');
  cy.get('div.pm_survey_planned_pack_date').contains('01-Aug-18');
  cy.get('div.actual_pack_date').contains('02-Aug-18');

  // Pickup Dates
  cy
    .get('.editable-panel-header')
    .contains('Dates')
    .siblings()
    .click();

  cy
    .get('input[name="dates.pm_survey_planned_pickup_date"]')
    .first()
    .type('8/2/2018')
    .blur();
  cy
    .get('input[name="dates.actual_pickup_date"]')
    .first()
    .type('8/3/2018')
    .blur();

  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save')
    .click();

  cy.reload();

  cy.get('div.requested_pickup_date').contains('15-May-19');
  cy.get('div.pm_survey_planned_pickup_date').contains('02-Aug-18');
  cy.get('div.actual_pickup_date').contains('03-Aug-18');

  // Delivery Dates
  cy
    .get('.editable-panel-header')
    .contains('Dates')
    .siblings()
    .click();

  cy
    .get('input[name="dates.pm_survey_planned_delivery_date"]')
    .first()
    .type('10/7/2018')
    .blur();
  cy
    .get('input[name="dates.actual_delivery_date"]')
    .first()
    .type('10/8/2018');

  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save')
    .click();

  cy.reload();

  cy.get('div.original_delivery_date').contains('15-May-19');
  cy.get('div.pm_survey_planned_delivery_date').contains('07-Oct-18');
  cy.get('div.actual_delivery_date').contains('08-Oct-18');
  cy.get('div.rdd').contains('07-Oct-18');

  // Notes
  cy
    .get('.editable-panel-header')
    .contains('Dates')
    .siblings()
    .click();

  cy
    .get('textarea[name="dates.pm_survey_notes"]')
    .first()
    .clear()
    .type('Notes notes notes for dates')
    .blur();

  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save')
    .click();

  cy.reload();

  cy.get('div.pm_survey_notes').contains('Notes notes notes for dates');

  // Verify Premove Survey contains the same data
  cy
    .get('.editable-panel-header')
    .contains('Premove Survey')
    .parent()
    .parent()
    .contains('07-Oct-18');
}

function tspUserClicksUploadOriginDocs(shipmentId) {
  cy.visit(`/shipments/${shipmentId}`, {
    onBeforeLoad(win) {
      cy.stub(win, 'open');
    },
  });

  cy
    .get('button')
    .contains('Upload Origin Docs')
    .click();

  cy
    .window()
    .its('open')
    .should('be.called');
}
