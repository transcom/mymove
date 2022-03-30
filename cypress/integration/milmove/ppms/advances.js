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

  const amountRequestedInputSelector = 'input[name="amountRequested"]';
  const errorMessageSelector = '[class="usa-error-message"]';
  const warningMessageSelector = '[data-testid="textInputWarning"]';

  // missing advance
  cy.get(amountRequestedInputSelector).clear().blur();
  cy.get(errorMessageSelector).contains('Required');
  cy.get(errorMessageSelector).next('div').find('input').should('have.id', 'amountRequested');

  // advance violates min
  cy.get(amountRequestedInputSelector).clear().type(0).blur();
  cy.get(errorMessageSelector).contains("The minimum advance request is $1. If you don't want an advance, select No.");
  cy.get(errorMessageSelector).next('div').find('input').should('have.id', 'amountRequested');
  cy.get(amountRequestedInputSelector).clear().type(1).blur();
  cy.get(errorMessageSelector).should('not.exist');

  // a warning is displayed when advance is greater than 60% of incentive
  cy.get(amountRequestedInputSelector).clear().type(8000).blur();
  cy.get(warningMessageSelector).contains('Reminder: your advance can not be more than $6,000');
  cy.get(warningMessageSelector).next('div').find('input').should('have.id', 'amountRequested');
  cy.get(amountRequestedInputSelector).clear().type(1).blur();
  cy.get(warningMessageSelector).should('not.exist');

  // advance violates max (over 100% of incentive)
  cy.get(amountRequestedInputSelector).clear().type(20000).blur();
  cy.get(errorMessageSelector).contains('Enter an amount less than $6,000');
  cy.get(errorMessageSelector).next('div').find('input').should('have.id', 'amountRequested');
  cy.get(amountRequestedInputSelector).clear().type(1).blur();
  cy.get(errorMessageSelector).should('not.exist');
}

function getToAdvancesPage() {
  // estimated_weights@ppm.unsubmitted
  const userId = '4512dc8c-c777-444e-b6dc-7971e398f2dc';

  signInAndNavigateFromHomePageToExistingPPMDateAndLocationPage(userId);
  navigateFromDateAndLocationPageToEstimatedWeightsPage('@patchShipment');
  navigateFromEstimatedWeightsPageToEstimatedIncentivePage();
  navigateFromEstimatedIncentivePageToAdvancesPage();
}
