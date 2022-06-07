import {
  customerStartsAddingAPPMShipment,
  deleteShipment,
  generalVerifyEstimatedIncentivePage,
  navigateFromDateAndLocationPageToEstimatedWeightsPage,
  navigateFromEstimatedWeightsPageToEstimatedIncentivePage,
  navigateToAgreementAndSign,
  setMobileViewport,
  signInAndNavigateFromHomePageToExistingPPMDateAndLocationPage,
  submitsAdvancePage,
  submitsDateAndLocation,
  submitsEstimatedWeightsAndProGear,
} from '../../../support/ppmCustomerShared';

import { submitMove } from '../../mymove/utilities/customer';

describe('Entire PPM onboarding flow', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.intercept('GET', '**/internal/moves/**/mto_shipments').as('getShipment');
    cy.intercept('POST', '**/internal/mto_shipments').as('createShipment');
    cy.intercept('PATCH', '**/internal/mto-shipments/**').as('patchShipment');
    cy.intercept('DELETE', '**/internal/mto-shipments/**').as('deleteShipment');
    cy.intercept('GET', '**/internal/moves/**/signed_certifications').as('signedCertifications');
  });

  const viewportType1 = [
    { viewport: 'desktop', isMobile: false, userId: '9b9ce6ed-70ba-4edf-b016-488c87fc1250' }, // profile_full_ppm@move.draft
    { viewport: 'mobile', isMobile: true, userId: '4fd6726d-2d05-4640-96dd-983bec236a9c' }, // full_ppm_mobile@complete.profile
  ];

  viewportType1.forEach(({ viewport, isMobile, userId }) => {
    it(`flows through happy path for existing shipment - ${viewport}`, () => {
      if (isMobile) {
        setMobileViewport();
      }

      navigateHappyPath(userId, isMobile);
    });
  });

  const viewportType2 = [
    { viewport: 'desktop', isMobile: false, userId: 'b54d5368-a633-4e3e-a8df-22133b9f8c7c' }, // happyPathWithEdits@ppm.unsubmitted
    { viewport: 'mobile', isMobile: true, userId: '9365990e-5813-4031-aa42-170886150912' }, // happyPathWithEditsMobile@ppm.unsubmitted
  ];

  viewportType2.forEach(({ viewport, isMobile, userId }) => {
    it(`happy path with edits and backs - ${viewport}`, () => {
      if (isMobile) {
        setMobileViewport();
      }

      navigateHappyPathWithEditsAndBacks(userId, isMobile);
    });
  });

  const viewportType3 = [
    { viewport: 'desktop', isMobile: false, userId: '57d58062-93ac-4eb7-b1da-21dd137e4f65' }, // deleteShipment@ppm.unsubmitted
    { viewport: 'mobile', isMobile: true, userId: '781cf194-4eb2-4def-9da6-01abdc62333d' }, // deleteShipmentMobile@ppm.unsubmitted
  ];

  viewportType3.forEach(({ viewport, isMobile, userId }) => {
    it(`deletes shipment - ${viewport}`, () => {
      if (isMobile) {
        setMobileViewport();
      }

      navigateDeletingShipment(userId, isMobile);
    });
  });
});

function navigateHappyPath(userId, isMobile = false) {
  signInAndNavigateFromHomePageToExistingPPMDateAndLocationPage(userId);
  submitsDateAndLocation();
  submitsEstimatedWeightsAndProGear();
  generalVerifyEstimatedIncentivePage(isMobile);
  submitsAdvancePage(true, isMobile);
  navigateToAgreementAndSign();
  submitMove('@signedCertifications');
  verifyStep5ExistsAndBtnIsDisabled();
}

function navigateDeletingShipment(userId, isMobile = false) {
  cy.apiSignInAsUser(userId);
  cy.wait('@getShipment');
  customerDeletesExistingShipment();
}

function navigateHappyPathWithEditsAndBacks(userId, isMobile = false) {
  signInAndNavigateFromHomePageToExistingPPMDateAndLocationPage(userId);

  submitAndVerifyUpdateDateAndLocation();

  submitsEstimatedWeightsAndProGear();
  verifyEstimatedWeightsAndProGear();

  verifyShipmentSpecificInfoOnEstimatedIncentivePage();
  generalVerifyEstimatedIncentivePage(isMobile);

  submitsAdvancePage(true, isMobile);

  navigateToAgreementAndSign();

  submitMove('@signedCertifications');
  verifyStep5ExistsAndBtnIsDisabled();
}

// update the form values by submitting and then return to the page to verify if the values persist and then return to the next page
function submitAndVerifyUpdateDateAndLocation() {
  cy.get('input[name="pickupPostalCode"]').clear().type('90210').blur();
  cy.get('input[name="secondaryPickupPostalCode"]').clear().type('90212').blur();
  cy.get('input[name="destinationPostalCode"]').clear().type('76127');
  // TODO: The user has secondary destination zips. We should test clearing this value by selecting the no radio btn. This doesn't work atm
  cy.get('input[name="sitExpected"][value="false"]').check();
  cy.get('input[name="expectedDepartureDate"]').clear().type('01 Feb 2022').blur();

  navigateFromDateAndLocationPageToEstimatedWeightsPage('@patchShipment');

  cy.get('button').contains('Back').click();

  // verify values
  cy.get('input[name="pickupPostalCode"]').should('have.value', '90210');
  cy.get('input[name="hasSecondaryPickupPostalCode"]').eq(0).should('be.checked').and('have.value', 'true');
  cy.get('input[name="secondaryPickupPostalCode"]').should('have.value', '90212');
  cy.get('input[name="destinationPostalCode"]').should('have.value', '76127');
  cy.get('input[name="hasSecondaryDestinationPostalCode"]').eq(0).should('be.checked').and('have.value', 'true');
  cy.get('input[name="expectedDepartureDate"]').should('have.value', '01 Feb 2022');
  cy.get('input[name="sitExpected"]').last().should('be.checked').and('have.value', 'false');

  navigateFromDateAndLocationPageToEstimatedWeightsPage('@patchShipment');
}

// verify page and submit to go to next page
function verifyEstimatedWeightsAndProGear() {
  cy.get('button').contains('Back').click();

  cy.get('input[name="estimatedWeight"]').should('have.value', '4,000');
  cy.get('input[name="hasProGear"][value="true"]').should('be.checked');
  cy.get('input[name="proGearWeight"]').should('be.visible').and('have.value', '500');
  cy.get('input[name="spouseProGearWeight"]').should('be.visible').and('have.value', '400');

  navigateFromEstimatedWeightsPageToEstimatedIncentivePage();
}

function verifyShipmentSpecificInfoOnEstimatedIncentivePage() {
  cy.get('.container li')
    .should('contain', '4,000 lbs')
    .and('contain', '90210')
    .and('contain', '76127')
    .and('contain', '01 Feb 2022');
}

function customerDeletesExistingShipment() {
  cy.get('[data-testid="shipment-list-item-container"]').as('shipmentListContainer');
  deleteShipment('@shipmentListContainer', 0);
}

function verifyStep5ExistsAndBtnIsDisabled() {
  cy.get('[data-testid="stepContainer5"]').within(() => {
    cy.get('button').contains('Upload PPM Documents').should('be.disabled');
    cy.get('p').contains('After a counselor approves your PPM, you will be able to:');
  });
}
