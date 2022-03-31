import {
  navigateFromDateAndLocationPageToEstimatedWeightsPage,
  signInAndNavigateFromHomePageToExistingPPMDateAndLocationPage,
  submitsEstimatedWeights,
  submitsEstimatedWeightsAndProGear,
} from '../../../support/ppmShared';

describe('PPM Onboarding - Add Estimated  Weight and Pro-gear', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.intercept('GET', '**/internal/moves/**/mto_shipments').as('getShipment');
    cy.intercept('PATCH', '**/internal/mto-shipments/**').as('patchShipment');
  });

  it('doesn’t allow SM to progress if form is in an invalid state', () => {
    // For this invalid tests we don't need to wait for the API calls, prioritizing consistency over performance
    getToEstimatedWeightsPage();
    invalidInputs();
  });

  it('can continue to next page', () => {
    getToEstimatedWeightsPage();
    submitsEstimatedWeights();
  });

  it('can continue to next page with progear added', () => {
    getToEstimatedWeightsPage();
    submitsEstimatedWeightsAndProGear();
  });
});

function invalidInputs() {
  cy.contains('Estimated weight');
  cy.url().should('include', '/estimated-weight');
  cy.get('p[class="usa-alert__text"]').contains('Total weight allowance for your move: 5,000 lbs');

  // missing required weight
  cy.get('input[name="estimatedWeight"]').clear().blur();
  cy.get('[class="usa-error-message"]').as('errorMessage');
  cy.get('@errorMessage').contains('Required');
  cy.get('@errorMessage').next('div').get('input').should('have.id', 'estimatedWeight');

  // estimated weight violates min
  cy.get('input[name="estimatedWeight"]').type(0).blur();
  cy.get('@errorMessage').contains('Enter a weight greater than 0 lbs');
  cy.get('@errorMessage').next('div').get('input').should('have.id', 'estimatedWeight');
  cy.get('input[name="estimatedWeight"]').clear().type(500).blur();
  cy.get('@errorMessage').should('not.exist');

  // a warning is displayed when estimated weight is greater than the SM's weight allowance
  cy.get('input[name="estimatedWeight"]').clear().type(17000).blur();
  cy.get('[data-testid="textInputWarning"]').as('warningMessage');
  cy.get('@warningMessage').contains(
    'This weight is more than your weight allowance. Talk to your counselor about what that could mean for your move.',
  );
  cy.get('@warningMessage').next('div').get('input').should('have.id', 'estimatedWeight');
  cy.get('input[name="estimatedWeight"]').clear().type(500).blur();
  cy.get('@warningMessage').should('not.exist');

  // pro gear violates max
  cy.get('input[name="hasProGear"][value="true"]').check({ force: true });
  cy.get('input[name="proGearWeight"]').type(5000).blur();
  cy.get('@errorMessage').contains('Enter a weight less than 2,000 lbs');
  cy.get('input[name="proGearWeight"]').clear().type(500).blur();
  cy.get('@errorMessage').should('not.exist');

  // When hasProGear is true show error if either personal or spouse pro gear isn't specified
  cy.get('input[name="proGearWeight"]').clear().blur();
  cy.get('@errorMessage').contains(
    "Enter a weight into at least one pro-gear field. If you won't have pro-gear, select No above.",
  );
  cy.get('input[name="proGearWeight"]').clear().type(500).blur();
  cy.get('@errorMessage').should('not.exist');

  // spouse pro gear max violation
  cy.get('input[name="spouseProGearWeight"]').clear().type(1000).blur();
  cy.get('@errorMessage').contains('Enter a weight less than 500 lbs');
  cy.get('input[name="spouseProGearWeight"]').clear().type(100).blur();
  cy.get('@errorMessage').should('not.exist');
}

function getToEstimatedWeightsPage() {
  // dates_and_locations@ppm.unsubmitted
  const userId = 'bbb469f3-f4bc-420d-9755-b9569f81715e';

  signInAndNavigateFromHomePageToExistingPPMDateAndLocationPage(userId);
  navigateFromDateAndLocationPageToEstimatedWeightsPage('@patchShipment');
}
