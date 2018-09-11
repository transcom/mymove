/* global cy, Cypress*/
export function fillAndSavePremoveSurvey() {
  // Click on edit Premove Survey
  cy
    .get('.editable-panel-header')
    .contains('Premove')
    .siblings()
    .click();

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
    .type('8/5/2018')
    .blur();
  cy
    .get('input[name="survey.pm_survey_weight_estimate"]')
    .first()
    .type('6000')
    .blur();
  cy
    .get('input[name="survey.pm_survey_progear_weight_estimate"]')
    .first()
    .type('7000')
    .blur();
  cy
    .get('input[name="survey.pm_survey_spouse_progear_weight_estimate"]')
    .first()
    .type('8000')
    .blur();
  cy
    .get('input[name="survey.pm_survey_notes"]')
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
  fillAndSavePremoveSurvey();

  // Verify data has been saved in the UI
  cy.get('span').contains('7,000 lbs');

  // Refresh browser and make sure changes persist
  cy.reload();

  cy
    .get('div.pm_survey_planned_delivery_date')
    .get('span')
    .contains('7,000 lbs');
  cy.get('div.pm_survey_notes').contains('Notes notes notes');
}
