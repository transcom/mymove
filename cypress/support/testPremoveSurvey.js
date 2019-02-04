/* global cy */

export function selectPreMoveSurveyPanel() {
  // Click on edit Premove Survey
  cy
    .get('.editable-panel-header')
    .contains('Premove')
    .siblings()
    .click();
}

export function fillAndSavePremoveSurvey() {
  // Enter details in form and save orders
  cy
    .get('input[name="survey.pm_survey_planned_pack_date"]')
    .first()
    .type('8/1/2018')
    .blur();
  cy
    .get('input[name="survey.pm_survey_planned_pickup_date"]')
    .first()
    .type('8/2/2018')
    .blur();
  cy
    .get('input[name="survey.pm_survey_planned_delivery_date"]')
    .first()
    .type('8/6/2018')
    .blur();
  cy
    .get('input[name="survey.pm_survey_conducted_date"]')
    .first()
    .type('7/20/2018')
    .blur();
  cy
    .get('input[name="survey.pm_survey_weight_estimate"]')
    .clear()
    .first()
    .type('6000')
    .blur();
  cy
    .get('input[name="survey.pm_survey_progear_weight_estimate"]')
    .clear()
    .first()
    .type('4000')
    .blur();
  cy
    .get('input[name="survey.pm_survey_spouse_progear_weight_estimate"]')
    .clear()
    .first()
    .type('800')
    .blur();
  cy
    .get('textarea[name="survey.pm_survey_notes"]')
    .clear()
    .first()
    .type('Notes notes notes')
    .blur();
  cy.get('select[name="survey.pm_survey_method"]').select('PHONE');

  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save')
    .click();
}

export function testPremoveSurvey() {
  selectPreMoveSurveyPanel();
  fillAndSavePremoveSurvey();

  // Verify data has been saved in the UI
  cy.get('span').contains('4,000 lbs');

  // Refresh browser and make sure changes persist
  cy.patientReload();

  cy
    .get('div.pm_survey_planned_delivery_date')
    .get('span')
    .contains('4,000 lbs');
  cy.get('div.pm_survey_notes').contains('Notes notes notes');
}
