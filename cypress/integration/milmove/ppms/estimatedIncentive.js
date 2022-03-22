import { verifyEstimatedIncentivePage } from '../../../support/ppmShared';

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
    verifyEstimatedIncentivePage(false);
  });

  it('mobile - go to estimated incentives page', () => {
    cy.viewport(479, 875);
    navigateToEstimatedIncentivePage();
    verifyEstimatedIncentivePage(true);
  });
});

function navigateToEstimatedIncentivePage() {
  // estimated_weights@ppm.unsubmitted
  const userId = '4512dc8c-c777-444e-b6dc-7971e398f2dc';
  cy.apiSignInAsUser(userId);
  cy.wait('@getShipment');

  // navigate to existing shipment
  cy.get('[data-testid="shipment-list-item-container"]').click();
  cy.wait('@getShipment');

  cy.get('[data-testid="ppm-date-and-location-submit"]').should('not.be.disabled').click();
  cy.wait('@patchShipment');

  cy.get('[data-testid="ppm-estimated-weights-submit"]').should('not.be.disabled').click();
  cy.wait('@patchShipment');
}
