import { tspUserVerifiesShipmentStatus } from '../../support/testTspStatus';

/* global cy */
describe('TSP User And the Premove Survey Button', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });

  it('tsp user uses action button to enter premove survey for accepted shipment', function() {
    tspUserGoesToShipment('queues/accepted', 'ACC4PM');
    tspUserClicksEnterPreMoveSurvey();
    tspUserFillsInPreMoveSurveyWizard();
    tspUserVerifiesPreMoveSurveyEntered();
    tspUserVerifiesShipmentStatus('Shipment accepted');
  });

  it('tsp user uses action button to enter premove survey for approved shipment', function() {
    tspUserGoesToShipment('queues/approved', 'ENTPMS');
    tspUserClicksEnterPreMoveSurvey();
    tspUserFillsInPreMoveSurveyWizard();
    tspUserVerifiesPreMoveSurveyEntered();
    tspUserVerifiesShipmentStatus('Pre-move survey complete');
  });

  it('tsp user does not see premove survey button if approved, but already completed', function() {
    tspUserGoesToShipment('queues/approved', 'APPPMS');
    tspUserVerifiesPreMoveSurveyEntered();
    tspUserVerifiesShipmentStatus('Pre-move survey complete');
  });
});

function tspUserClicksEnterPreMoveSurvey() {
  cy
    .get('button')
    .contains('Enter pre-move survey')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Enter pre-move survey')
    .click();
}

function tspUserFillsInPreMoveSurveyWizard() {
  cy
    .get('input[name="pm_survey_planned_pack_date"]')
    .type('10/24/2018')
    .blur();

  cy
    .get('a')
    .contains('Cancel')
    .click();

  cy.get('input[name="pm_survey_planned_pack_date"]').should('not.exist');

  cy
    .get('button')
    .contains('Enter pre-move survey')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Enter pre-move survey')
    .click();

  cy
    .get('input[name="pm_survey_planned_pack_date"]')
    .first()
    .type('8/1/2018')
    .blur();
  cy
    .get('input[name="pm_survey_planned_pickup_date"]')
    .first()
    .type('8/2/2018')
    .blur();
  cy
    .get('input[name="pm_survey_planned_delivery_date"]')
    .first()
    .type('8/5/2018')
    .blur();
  cy
    .get('input[name="pm_survey_conducted_date"]')
    .first()
    .type('7/20/2018')
    .blur();
  cy
    .get('input[name="pm_survey_weight_estimate"]')
    .clear()
    .first()
    .type('6000')
    .blur();

  cy
    .get('button')
    .contains('Done')
    .should('be.disabled');

  cy.get('select[name="pm_survey_method"]').select('PHONE');

  cy
    .get('button')
    .contains('Done')
    .should('be.enabled');

  cy
    .get('input[name="pm_survey_progear_weight_estimate"]')
    .clear()
    .first()
    .type('4000')
    .blur();
  cy
    .get('input[name="pm_survey_spouse_progear_weight_estimate"]')
    .clear()
    .first()
    .type('800')
    .blur();
  cy
    .get('textarea[name="pm_survey_notes"]')
    .clear()
    .first()
    .type('Notes notes notes')
    .blur();

  cy
    .get('button')
    .contains('Done')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Done')
    .click();
}

function tspUserVerifiesPreMoveSurveyEntered() {
  cy.get('button').should('not.contain', 'Enter pre-move survey');
}

function tspUserGoesToShipment(queue, locator) {
  cy.patientVisit(queue);

  // Find shipment and open it
  cy
    .get('div')
    .contains(locator)
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });
}
