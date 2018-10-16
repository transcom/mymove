/* global cy */
describe('TSP User Completes Dates Panel', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });
  it('tsp user completes dates panel', function() {
    tspUserEntersDates();
  });
});

function tspUserEntersDates() {
  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find shipment and open it
  cy
    .get('div')
    .contains('DATESP')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  cy
    .get('.editable-panel-header')
    .contains('Dates')
    .siblings()
    .click();

  // Enter details in form and save dates

  // Conducted Date
  cy
    .get('input[name="dates.pm_survey_conducted_date"]')
    .first()
    .type('7/20/2018')
    .blur();
  cy.get('select[name="dates.pm_survey_method"]').select('PHONE');
  // Pack Dates
  // TODO: ADD original_pack_date
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
  // Pickup Dates
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
  // Delivery Dates
  // TODO: Add original_delivery_date
  cy
    .get('input[name="dates.pm_survey_planned_delivery_date"]')
    .first()
    .type('10/7/2018')
    .blur();
  cy
    .get('input[name="dates.actual_delivery_date"]')
    .first()
    .type('10/8/2018');
  // Notes
  cy
    .get('textarea[name="dates.pm_survey_notes"]')
    .first()
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

  // Verify data has been saved in the UI
  cy.get('span').contains('Notes notes notes for dates');

  // Refresh browser and make sure changes persist
  cy.reload();

  cy.get('div.pm_survey_conducted_date').contains('20-Jul-18');
  cy.get('div.original_pack_date').contains('TODO');
  cy.get('div.pm_survey_planned_pack_date').contains('01-Aug-18');
  cy.get('div.actual_pack_date').contains('02-Aug-18');
  cy.get('div.requested_pickup_date').contains('15-May-19');
  cy.get('div.pm_survey_planned_pickup_date').contains('02-Aug-18');
  cy.get('div.actual_pickup_date').contains('03-Aug-18');
  cy.get('div.original_delivery_date').contains('TODO');
  cy.get('div.pm_survey_planned_delivery_date').contains('07-Oct-18');
  cy.get('div.actual_delivery_date').contains('08-Oct-18');
  cy.get('div.rdd').contains('08-Oct-18');

  // Verify Premove Survey contains the same data
  cy
    .get('.editable-panel-header')
    .contains('Premove Survey')
    .parent()
    .parent()
    .contains('07-Oct-18');
}
