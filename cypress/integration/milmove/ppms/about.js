import {
  deleteShipment,
  navigateFromAdvancesPageToReviewPage,
  navigateFromDateAndLocationPageToEstimatedWeightsPage,
  navigateFromEstimatedIncentivePageToAdvancesPage,
  navigateFromEstimatedWeightsPageToEstimatedIncentivePage,
  navigateFromHomePageToReviewPage,
  navigateFromReviewPageToHomePage,
  navigateToAgreementAndSign,
  setMobileViewport,
  signInAndNavigateFromHomePageToExistingPPMDateAndLocationPage,
  signInAndNavigateFromHomePageToReviewPage,
  signInAndNavigateToAboutPage,
  submitsAdvancePage,
} from '../../../support/ppmCustomerShared';

describe('About Your PPM - With Advances', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.intercept('GET', '**/internal/moves/**/mto_shipments').as('getShipment');
    cy.intercept('PATCH', '**/internal/mto-shipments/**').as('patchShipment');
    cy.intercept('DELETE', '**/internal/mto-shipments/**').as('deleteShipment');
    cy.intercept('GET', '**/internal/moves/**/signed_certifications').as('signedCertifications');
  });
  // Get another user for desktop
  const viewportType = [
    // { viewport: 'desktop', isMobile: false, userId: 'cde987a1-a717-4a61-98b5-1f05e2e0844d' }, // readyToFinish@ppm.approved
    { viewport: 'mobile', isMobile: true, userId: 'cde987a1-a717-4a61-98b5-1f05e2e0844d' }, // readyToFinish@ppm.approved
  ];

  // TODO: Add a test for when a user logs off but logs back in - they should see the fields on teh page populated with data that they had already filled out in other forms (zips and Advance question)
  viewportType.forEach(({ viewport, isMobile, userId }) => {
    it(`can upload with an advance - ${viewport}`, () => {
      if (isMobile) {
        setMobileViewport();
      }
      signInAndNavigateToAboutPage(userId);

      //
      //     saveAndContinueAboutPage(true, isMobile);
      //   });
      //
      //   it(`can request an advance - ${viewport}`, () => {
      //     if (isMobile) {
      //       setMobileViewport();
      //     }
      //
      //     getToAdvancesPage();
      //
      //     submitsAdvancePage(true, isMobile);
    });
  });
});
