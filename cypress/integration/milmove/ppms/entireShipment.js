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
    // TODO: need to change id + add email associated with user
    const userId = '3b9360a3-3304-4c60-90f4-83d687884070';
    cy.apiSignInAsUser(userId);
    cy.wait('@getShipment');
    customerChoosesAPPMMove();
    submitsDateAndLocation();
    submitsEstimatedWeightsAndProgear();
    verifyEstimatedIncentivePage();
  });

  it('happy path with edits and backs', () => {
    // TODO: need to change id + add email associated with user
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
    verifyEstimatedIncentivePage();
  });
});

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
