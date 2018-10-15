import { fillAndSavePremoveSurvey } from '../../support/testPremoveSurvey';

/* global cy */
describe('TSP Interacts With the Weights & Items Panel', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });

  it('tsp user enters actual weight', function() {
    tspUserEntersActualWeight();
  });

  it('tsp user sees estimated weights from the customer and TSP', function() {
    tspUserSeesEstimatedWeights();
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
    expect(text).to.include('7,000 lbs');
  });
  cy.get('.spouse_progear_weight_estimate').should($div => {
    const text = $div.text();
    expect(text).to.include('Customer estimate');
    expect(text).to.include('312 lbs');
  });
  cy.get('.pm_survey_spouse_progear_weight_estimate').should($div => {
    const text = $div.text();
    expect(text).to.include('TSP estimate');
    expect(text).to.include('8,000 lbs');
  });
}

function tspUserSeesEstimatedWeights() {
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find shipment and open it
  cy
    .get('div')
    .contains('BACON4')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

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

function tspUserEntersActualWeight() {
  // Open new shipments queue
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/queues\/new/);
  });

  // Find shipment and open it
  cy
    .get('div')
    .contains('BACON4')
    .dblclick();

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
  });

  // Check that the initial display view for actual weight is correct
  withinWeightsAndItemsPanel(() => {
    cy.get('.actual_weight').should($div => {
      const text = $div.text();
      expect(text).to.include('Actual Weight');
      expect(text).to.include('missing');
    });
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

  // Fill out the actual weight and save it
  withinWeightsAndItemsPanel(() => {
    cy.get('label[for="weights.actual_weight"]').should('have.text', 'Actual Weight');
    cy
      .get('input[name="weights.actual_weight"]')
      .first()
      .type('40000')
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
    cy.get('.actual_weight').should($div => {
      const text = $div.text();
      expect(text).to.include('Actual Weight');
      expect(text).to.include('40,000 lbs');
    });
  });
}
