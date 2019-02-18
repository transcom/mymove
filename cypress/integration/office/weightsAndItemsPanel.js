import { selectPreMoveSurveyPanel, fillAndSavePremoveSurvey } from '../../support/testPremoveSurvey';

/* global cy */
describe('Office User Interacts With the Weights & Items Panel', function() {
  beforeEach(() => {
    cy.signIntoOffice();
  });

  it('office user enters net weight', function() {
    officeUserEntersNetWeight();
  });

  it('office user sees estimated weights from the customer and office', function() {
    officeUserSeesEstimatedWeights();
  });
});

function withinWeightsAndItemsPanel(func) {
  cy
    .get('.editable-panel')
    .contains('Weights & Items')
    .get('.editable-panel-content')
    .within(panel => {
      func(panel);
    });
}

function testReadOnlyWeights() {
  cy.get('.weight_estimate').should($div => {
    const text = $div.text();
    expect(text).to.include('Customer estimate');
    expect(text).to.include('2,000 lbs');
  });
  cy.get('.pm_survey_weight_estimate').should($div => {
    const text = $div.text();
    expect(text).to.include('TSP estimate');
    expect(text).to.include('6,000 lbs');
  });
  cy.get('.progear_weight_estimate').should($div => {
    const text = $div.text();
    expect(text).to.include('Customer estimate');
    expect(text).to.include('225 lbs');
  });
  cy.get('.pm_survey_progear_weight_estimate').should($div => {
    const text = $div.text();
    expect(text).to.include('TSP estimate');
    expect(text).to.include('4,000 lbs');
  });
  cy.get('.spouse_progear_weight_estimate').should($div => {
    const text = $div.text();
    expect(text).to.include('Customer estimate');
    expect(text).to.include('3,120 lbs');
  });
  cy.get('.pm_survey_spouse_progear_weight_estimate').should($div => {
    const text = $div.text();
    expect(text).to.include('TSP estimate');
    expect(text).to.include('800 lbs');
  });
}

function openMoveHhgPanel() {
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find shipment and open it
  cy.selectQueueItemMoveLocator('WTSPNL');

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
  });

  cy
    .get('.nav-tab')
    .contains('HHG')
    .click();
}

function officeUserSeesEstimatedWeights() {
  openMoveHhgPanel();
  selectPreMoveSurveyPanel();
  fillAndSavePremoveSurvey();

  // Check that the display view is correct for the estimated weights
  withinWeightsAndItemsPanel(() => testReadOnlyWeights);

  cy
    .get('.editable-panel-header')
    .contains('Weights & Items')
    .siblings()
    .click();

  // Check that the edit view is correct for the estimated weights
  withinWeightsAndItemsPanel(() => testReadOnlyWeights);

  // Verify the user can cancel
  cy
    .get('button')
    .contains('Cancel')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Cancel')
    .click();
}

function officeUserEntersNetWeight() {
  openMoveHhgPanel();
  // Check that the initial display view for net weight is correct
  withinWeightsAndItemsPanel(() => {
    cy.get('.net_weight').should('not.exist');
  });

  cy
    .get('.editable-panel-header')
    .contains('Weights & Items')
    .siblings()
    .click();

  cy
    .get('button')
    .contains('Save')
    .should('not.be.enabled');

  // Fill out the net weight and save it
  withinWeightsAndItemsPanel(() => {
    cy.get('label[for="weights.pm_survey_weight_estimate"]').should('have.text', 'TSP estimateOptional');
    cy
      .get('input[name="weights.pm_survey_weight_estimate"]')
      .clear()
      .first()
      .type('5000')
      .blur();
    cy.get('label[for="weights.gross_weight"]').should('have.text', 'Gross');
    cy
      .get('input[name="weights.gross_weight"]')
      .first()
      .type('30000')
      .blur();
    cy.get('label[for="weights.tare_weight"]').should('have.text', 'Tare');
    cy
      .get('input[name="weights.tare_weight"]')
      .first()
      .type('10000')
      .blur();
    cy.get('label[for="weights.net_weight"]').should('have.text', 'Net (Gross - Tare)');
    cy
      .get('input[name="weights.net_weight"]')
      .first()
      .type('40000')
      .blur();

    cy.get('label[for="weights.pm_survey_progear_weight_estimate"]').should('have.text', 'Service memberOptional');
    cy
      .get('input[name="weights.pm_survey_progear_weight_estimate"]')
      .clear()
      .first()
      .type('4000')
      .blur();
    cy.get('label[for="weights.pm_survey_spouse_progear_weight_estimate"]').contains('Spouse');
    cy
      .get('input[name="weights.pm_survey_spouse_progear_weight_estimate"]')
      .clear()
      .first()
      .type('800')
      .blur();
  });

  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save')
    .click();

  // Verify that the entered weight displays properly
  withinWeightsAndItemsPanel(() => {
    cy.get('.net_weight').should($div => {
      const text = $div.text();
      expect(text).to.include('Actual');
      expect(text).to.include('40,000 lbs');
    });
  });
}
