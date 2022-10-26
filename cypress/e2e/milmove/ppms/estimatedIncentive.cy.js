import {
  generalVerifyEstimatedIncentivePage,
  navigateFromDateAndLocationPageToEstimatedWeightsPage,
  navigateFromEstimatedWeightsPageToEstimatedIncentivePage,
  setMobileViewport,
  signInAndNavigateFromHomePageToExistingPPMDateAndLocationPage,
} from '../../../support/ppmCustomerShared';

describe('PPM Onboarding - Estimated Incentive', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.intercept('GET', '**/internal/moves/**/mto_shipments').as('getShipment');
    cy.intercept('PATCH', '**/internal/mto-shipments/**').as('patchShipment');
  });

  it('go to estimated incentives page', () => {
    navigateToEstimatedIncentivePage();
    verifyShipmentSpecificInfo();
    generalVerifyEstimatedIncentivePage(false);
  });

  it('mobile - go to estimated incentives page', () => {
    setMobileViewport();
    navigateToEstimatedIncentivePage();
    verifyShipmentSpecificInfo();
    generalVerifyEstimatedIncentivePage(true);
  });
});

function navigateToEstimatedIncentivePage() {
  // estimated_weights@ppm.unsubmitted
  const userId = '4512dc8c-c777-444e-b6dc-7971e398f2dc';

  signInAndNavigateFromHomePageToExistingPPMDateAndLocationPage(userId);
  navigateFromDateAndLocationPageToEstimatedWeightsPage('@patchShipment');
  navigateFromEstimatedWeightsPageToEstimatedIncentivePage();
}

function verifyShipmentSpecificInfo() {
  cy.get('.container h2').contains('$10,000');
}
