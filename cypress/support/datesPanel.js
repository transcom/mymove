/* global cy */

export function userEntersDates() {
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

  cy.patientReload();

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

  cy.patientReload();

  cy.get('div.original_pack_date').contains('11-May-18');
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

  cy.patientReload();

  cy.get('div.requested_pickup_date').contains('15-May-18');
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

  cy.patientReload();

  cy.get('div.original_delivery_date').contains('21-May-18');
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

  cy.patientReload();

  cy.get('div.pm_survey_notes').contains('Notes notes notes for dates');
}

export function userEntersAndRemovesDates() {
  // Enter details in form and save dates

  // Conducted Date

  cy
    .get('input[name="dates.pm_survey_conducted_date"]')
    .first()
    .type('7/20/2018')
    .blur();
  cy.get('select[name="dates.pm_survey_method"]').select('PHONE');

  // Pack Dates
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
  cy
    .get('input[name="dates.pm_survey_planned_delivery_date"]')
    .first()
    .type('10/7/2018')
    .blur();
  cy
    .get('input[name="dates.actual_delivery_date"]')
    .first()
    .type('10/8/2018')
    .blur();

  // Notes

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

  cy.patientReload();

  cy.get('div.pm_survey_conducted_date').contains('20-Jul-18');
  cy.get('div.pm_survey_method').contains('Phone');
  cy.get('div.original_pack_date').contains('11-May-18');
  cy.get('div.pm_survey_planned_pack_date').contains('01-Aug-18');
  cy.get('div.actual_pack_date').contains('02-Aug-18');
  cy.get('div.requested_pickup_date').contains('15-May-18');
  cy.get('div.pm_survey_planned_pickup_date').contains('02-Aug-18');
  cy.get('div.actual_pickup_date').contains('03-Aug-18');
  cy.get('div.original_delivery_date').contains('21-May-18');
  cy.get('div.pm_survey_planned_delivery_date').contains('07-Oct-18');
  cy.get('div.actual_delivery_date').contains('08-Oct-18');
  cy.get('div.rdd').contains('07-Oct-18');
  cy.get('div.pm_survey_notes').contains('Notes notes notes for dates');

  // Now remove all the dates
  const zeroTime = '{selectall}01/01/0001{enter}';

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
    .clear()
    .type(zeroTime)
    .blur();
  cy.get('select[name="dates.pm_survey_method"]').select('PHONE');

  // Pack Dates
  cy
    .get('input[name="dates.pm_survey_planned_pack_date"]')
    .first()
    .clear()
    .type(zeroTime)
    .blur();
  cy
    .get('input[name="dates.actual_pack_date"]')
    .first()
    .clear()
    .type(zeroTime)
    .blur();

  // Pickup Dates
  cy
    .get('input[name="dates.pm_survey_planned_pickup_date"]')
    .first()
    .clear()
    .type(zeroTime)
    .blur();
  cy
    .get('input[name="dates.actual_pickup_date"]')
    .first()
    .clear()
    .type(zeroTime)
    .blur();

  // Delivery Dates
  cy
    .get('input[name="dates.pm_survey_planned_delivery_date"]')
    .first()
    .clear()
    .type(zeroTime)
    .blur();
  cy
    .get('input[name="dates.actual_delivery_date"]')
    .first()
    .clear()
    .type(zeroTime)
    .blur();

  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save')
    .click();

  cy.patientReload();

  cy.get('div.pm_survey_conducted_date').contains('missing');
  cy.get('div.original_pack_date').contains('11-May-18');
  cy.get('div.pm_survey_planned_pack_date').contains('missing');
  cy.get('div.actual_pack_date').contains('missing');
  cy.get('div.requested_pickup_date').contains('15-May-18');
  cy.get('div.pm_survey_planned_pickup_date').contains('missing');
  cy.get('div.actual_pickup_date').contains('missing');
  cy.get('div.original_delivery_date').contains('21-May-18');
  cy.get('div.pm_survey_planned_delivery_date').contains('missing');
  cy.get('div.actual_delivery_date').contains('missing');
  cy.get('div.rdd').contains('21-May-18');
}
