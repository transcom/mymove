import {
  navigateFromAdvancesPageToReviewPage,
  navigateFromDateAndLocationPageToEstimatedWeightsPage,
  navigateFromEstimatedIncentivePageToAdvancesPage,
  navigateFromEstimatedWeightsPageToEstimatedIncentivePage,
  setMobileViewport,
  signInAndNavigateFromHomePageToExistingPPMDateAndLocationPage,
  submitsAdvancePage,
  submitsEstimatedWeights,
  submitsEstimatedWeightsAndProGear,
} from '../../../support/ppmShared';

describe('PPM On-boarding - Advances', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.intercept('GET', '**/internal/moves/**/mto_shipments').as('getShipment');
    cy.intercept('PATCH', '**/internal/mto-shipments/**').as('patchShipment');
  });

  it('does not allow SM to progress if form is in an invalid state', () => {
    getToAdvancesPage();

    invalidInputs();
  });

  const viewportType = [
    { viewport: 'desktop', isMobile: false },
    { viewport: 'mobile', isMobile: true },
  ];

  viewportType.forEach(({ viewport, isMobile }) => {
    it(`can opt to not receive an advance - ${viewport}`, () => {
      if (isMobile) {
        setMobileViewport();
      }

      getToAdvancesPage();

      submitsAdvancePage(false, isMobile);
    });

    it(`can request an advance - ${viewport}`, () => {
      if (isMobile) {
        setMobileViewport();
      }

      getToAdvancesPage();

      submitsAdvancePage(true, isMobile);
    });
  });
});

function invalidInputs() {
  cy.get('input[name="advanceRequested"][value="true"]').check({ force: true });

  cy.get('button').contains('Save & Continue').as('saveButton');

  // missing advance
  cy.get('input[name="amountRequested"]').as('amountRequestedInput');
  cy.get('@amountRequestedInput').clear().blur();
  cy.get('[class="usa-error-message"]').as('errorMessage');
  cy.get('@errorMessage').contains('Required');
  cy.get('@errorMessage').next('div').find('input').should('have.id', 'amountRequested');
  cy.get('@saveButton').should('be.disabled');
  cy.get('@amountRequestedInput').clear().type(1).blur();
  cy.get('@errorMessage').should('not.exist');

  // advance violates min
  cy.get('@amountRequestedInput').clear().type(0).blur();
  cy.get('@errorMessage').contains("The minimum advance request is $1. If you don't want an advance, select No.");
  cy.get('@errorMessage').next('div').find('input').should('have.id', 'amountRequested');
  cy.get('@saveButton').should('be.disabled');
  cy.get('@amountRequestedInput').clear().type(1).blur();
  cy.get('@errorMessage').should('not.exist');

  // a warning is displayed when advance is greater than 60% of incentive
  cy.get('@amountRequestedInput').clear().type(8000).blur();
  cy.get('[data-testid="textInputWarning"]').as('warningMessage');
  cy.get('@warningMessage').contains('Reminder: your advance can not be more than $6,000');
  cy.get('@warningMessage').next('div').find('input').should('have.id', 'amountRequested');
  cy.get('@saveButton').should('be.disabled');
  cy.get('@amountRequestedInput').clear().type(1).blur();
  cy.get('@warningMessage').should('not.exist');

  // advance violates max (over 100% of incentive)
  cy.get('@amountRequestedInput').clear().type(20000).blur();
  cy.get('@errorMessage').contains('Enter an amount less than $6,000');
  cy.get('@errorMessage').next('div').find('input').should('have.id', 'amountRequested');
  cy.get('@saveButton').should('be.disabled');
  cy.get('@amountRequestedInput').clear().type(1).blur();
  cy.get('@errorMessage').should('not.exist');
}

function getToAdvancesPage() {
  // estimated_weights@ppm.unsubmitted
  const userId = '4512dc8c-c777-444e-b6dc-7971e398f2dc';

  signInAndNavigateFromHomePageToExistingPPMDateAndLocationPage(userId);
  navigateFromDateAndLocationPageToEstimatedWeightsPage('@patchShipment');
  navigateFromEstimatedWeightsPageToEstimatedIncentivePage();
  navigateFromEstimatedIncentivePageToAdvancesPage();
}
