import {
  navigateFromReviewPageToProGearPage,
  setMobileViewport,
  signInAndNavigateToWeightTicketPage,
  submitWeightTicketPage,
} from '../../../support/ppmCustomerShared';

describe('Progear', function () {
  before(() => {
    cy.prepareCustomerApp();
  });
  beforeEach(() => {
    cy.intercept('GET', '**/internal/moves/**/mto_shipments').as('getShipment');
  });
  const viewportType = [
    { viewport: 'desktop', isMobile: false, userId: '9ec731d8-f347-4d34-8b54-4ce9e6ea3282' }, // actualPPMDateZIPAdvanceDone@ppm.approved
    { viewport: 'mobile', isMobile: true, userId: '2a0146c4-ec9a-4efc-a94c-6c2849c3e167' }, // actualPPMDateZIPAdvanceDone2@ppm.approved
  ];
  viewportType.forEach(({ viewport, isMobile, userId }) => {
    it(`page loads - ${viewport}`, () => {
      if (isMobile) {
        setMobileViewport();
      }

      signInAndNavigateToWeightTicketPage(userId);
      submitWeightTicketPage();
      navigateFromReviewPageToProGearPage();
    });
  });
});
