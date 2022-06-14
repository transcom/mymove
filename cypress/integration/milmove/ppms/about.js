import {
  deleteShipment,
  setMobileViewport,
  signInAndNavigateToAboutPageWithAdvance,
  signInAndNavigateToAboutPageWithoutAdvance,
} from '../../../support/ppmCustomerShared';

describe('About Your PPM', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.intercept('GET', '**/internal/moves/**/mto_shipments').as('getShipment');
    cy.intercept('PATCH', '**/internal/mto-shipments/**').as('patchShipment');
  });
  // Get another user for desktop
  const viewportType = [
    { viewport: 'desktop', isMobile: false },
    { viewport: 'mobile', isMobile: true },
  ];

  viewportType.forEach(({ viewport, isMobile, userId }) => {
    it(`can upload with an advance - ${viewport}`, () => {
      userId = 'c28b2eb1-975f-49f7-b8a3-c7377c0da908'; // readyToFinish2@ppm.approved
      if (isMobile) {
        setMobileViewport();
      }

      signInAndNavigateToAboutPageWithAdvance(userId);
    });
  });

  viewportType.forEach(({ viewport, isMobile, userId }) => {
    it(`can upload without an advance - ${viewport}`, () => {
      userId = 'cde987a1-a717-4a61-98b5-1f05e2e0844d'; // readyToFinish@ppm.approved
      if (isMobile) {
        setMobileViewport();
      }
      signInAndNavigateToAboutPageWithoutAdvance(userId);
    });
  });
});
