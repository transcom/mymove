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
    { viewport: 'desktop', isMobile: false, userId: '33eabbb6-416d-4d91-ba5b-bfd7d35e3037' }, // progearWeightTicket@ppm.approved
    // TODO: fill in another userId here:
    { viewport: 'mobile', isMobile: true, userId: '33eabbb6-416d-4d91-ba5b-bfd7d35e3037' }, // progearWeightTicket2@ppm.approved
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
