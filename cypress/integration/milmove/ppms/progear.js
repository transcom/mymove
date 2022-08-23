import {
  navigateFromReviewPageToProGearPage,
  setMobileViewport,
  signInAndNavigateToWeightTicketPage,
} from '../../../support/ppmCustomerShared';

describe('Progear', function () {
  before(() => {
    cy.prepareCustomerApp();
  });
  beforeEach(() => {
    cy.intercept('GET', '**/internal/moves/**/mto_shipments').as('getShipment');
  });
  const viewportType = [
    { viewport: 'desktop', isMobile: false, userId: '88007896-6ae7-4600-866a-873d3bc67fd3' }, // actualPPMDateZIPAdvanceDone@ppm.approved
    { viewport: 'mobile', isMobile: true, userId: '22dba194-3d9a-49c6-8328-718dd945292f' }, // actualPPMDateZIPAdvanceDone2@ppm.approved
  ];
  viewportType.forEach(({ viewport, isMobile, userId }) => {
    it(`page loads - ${viewport}`, () => {
      if (isMobile) {
        setMobileViewport();
      }

      signInAndNavigateToWeightTicketPage(userId);
      navigateFromReviewPageToProGearPage();
    });
  });
});
