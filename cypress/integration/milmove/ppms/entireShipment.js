import {
  customerChoosesAPPMMove,
  submitsDateAndLocation,
  submitsEstimatedWeightsAndProgear,
  verifyEstimatedIncentivePage,
} from '../../../support/ppmShared';

describe('Entire PPM onboarding flow', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.intercept('GET', '**/internal/moves/**/mto_shipments').as('getShipment');
    cy.intercept('POST', '**/internal/mto_shipments').as('createShipment');
    cy.intercept('PATCH', '**/internal/mto-shipments/**').as('patchShipment');
  });

  it('happy path with new shipment', () => {
    // profile_full_ppm@move.draft
    const userId = '9b9ce6ed-70ba-4edf-b016-488c87fc1250';
    navigateHappyPath(userId, false);
  });

  it('mobile - happy path with new shipment', () => {
    // full_ppm_mobile@complete.profile
    const userId = '4fd6726d-2d05-4640-96dd-983bec236a9c';
    cy.viewport(479, 875);
    navigateHappyPath(userId, true);
  });

  it('happy path with edits and backs', () => {
    navigateHappyPathWithEditsAndBacks(false);
  });

  it('mobile - happy path with edits and backs', () => {
    cy.viewport(479, 875);
    navigateHappyPathWithEditsAndBacks(true);
  });
});

function navigateHappyPath(userId, isMobile = false) {
  cy.apiSignInAsUser(userId);
  cy.wait('@getShipment');
  customerChoosesAPPMMove();
  submitsDateAndLocation();
  submitsEstimatedWeightsAndProgear();
  verifyEstimatedIncentivePage(isMobile);
}

function navigateHappyPathWithEditsAndBacks(isMobile = false) {
  // TODO: need to change id to be unique + add email associated with user
  const userId = '4512dc8c-c777-444e-b6dc-7971e398f2dc';
  cy.apiSignInAsUser(userId);
  cy.wait('@getShipment');

  // navigate to existing shipment
  cy.get('[data-testid="shipment-list-item-container"]').click();
  cy.wait('@getShipment');

  submitAndVerifyUpdateDateAndLocation();

  submitsEstimatedWeightsAndProgear();
  verifyEstimatedWeightsAndProgear();

  verifySummaryInfoOnEstimatedIncentivePage();
  verifyEstimatedIncentivePage(isMobile);
}

// update the form values by submitting and then return to the page to verify if the values persist and then return to the next page
function submitAndVerifyUpdateDateAndLocation() {
  cy.get('input[name="pickupPostalCode"]').clear().type('90210').blur();

  cy.get('input[name="destinationPostalCode"]').clear().type('76127');
  cy.get('input[name="expectedDepartureDate"]').clear().type('15 Apr 2022').blur();

  cy.get('[data-testid="ppm-date-and-location-submit"]').contains('Save & Continue').click();
  cy.wait('@patchShipment');

  cy.get('[data-testid="ppm-estimated-weights-back"]').click();

  // verify values
  cy.get('input[name="pickupPostalCode"]').should('have.value', '90210');
  cy.get('input[name="hasSecondaryPickupPostalCode"]').eq(1).should('be.checked').and('have.value', 'false');
  cy.get('input[name="destinationPostalCode"]').should('have.value', '76127');
  cy.get('input[name="hasSecondaryDestinationPostalCode"]').eq(1).should('be.checked').and('have.value', 'false');
  cy.get('input[name="expectedDepartureDate"]').should('have.value', '15 Apr 2022');
  // TODO: We want to update the sit expected value and see if it saves (right now we are not updating this value)
  cy.get('input[name="sitExpected"]').eq(1).should('be.checked').and('have.value', 'false');

  // go to next page
  cy.get('[data-testid="ppm-date-and-location-submit"]').click();
  cy.wait('@patchShipment');
}

// verify page and submit to go to next page
function verifyEstimatedWeightsAndProgear() {
  cy.get('[data-testid="ppm-estimated-incentive-back"]').click();
  cy.get('input[name="estimatedWeight"]').should('have.value', '500');
  cy.get('input[name="hasProGear"][value="true"]').should('be.checked');
  cy.get('input[name="proGearWeight"]').should('be.visible').and('have.value', '500');
  cy.get('input[name="spouseProGearWeight"]').should('be.visible').and('have.value', '0');

  // go to next page
  cy.get('[data-testid="ppm-estimated-weights-submit"]').click();
  cy.wait('@patchShipment');
}

function verifySummaryInfoOnEstimatedIncentivePage() {
  cy.get('[data-testid="ppm-estimated-incentive-relevant-fields"]')
    .should('contain', '500 lbs')
    .and('contain', '90210')
    .and('contain', '76127')
    .and('contain', '15 Apr 2022');
}
