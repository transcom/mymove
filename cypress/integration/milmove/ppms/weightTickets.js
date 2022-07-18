import { deleteShipment, setMobileViewport, signInAndNavigateToAboutPage } from '../../../support/ppmCustomerShared';

// TODOS
// - Create another user for desktop
describe('Weight Tickets', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.intercept('GET', '**/internal/moves/**/mto_shipments').as('getShipment');
    cy.intercept('PATCH', '**/internal/mto-shipments/**').as('patchShipment');
  });

  const viewportType = [
    { viewport: 'desktop', isMobile: false },
    { viewport: 'mobile', isMobile: true },
  ];

  viewportType.forEach(({ viewport, isMobile, userId }) => {
    it(`proceed with weight ticket documents - ${viewport}`, () => {
      userId = '88007896-6ae7-4600-866a-873d3bc67fd3'; // actualPPMDateZIPAdvanceDone@ppm.approved
      if (isMobile) {
        setMobileViewport();
      }

      signInAndNavigateToWeightTicketPage(userId, true);
    });
  });

  viewportType.forEach(({ viewport, isMobile, userId }) => {
    it(`proceed without weight tickets (constructed weight) - ${viewport}`, () => {
      userId = '88007896-6ae7-4600-866a-873d3bc67fd3'; // actualPPMDateZIPAdvanceDone@ppm.approved
      if (isMobile) {
        setMobileViewport();
      }
      signInAndNavigateToAboutPage(userId, false);
    });
  });

  viewportType.forEach(({ viewport, isMobile, userId }) => {
    it(`proceed with claiming trailer - ${viewport}`, () => {
      userId = '88007896-6ae7-4600-866a-873d3bc67fd3'; // actualPPMDateZIPAdvanceDone@ppm.approved
      if (isMobile) {
        setMobileViewport();
      }
      signInAndNavigateToAboutPage(userId, false);
    });
  });

  viewportType.forEach(({ viewport, isMobile, userId }) => {
    it(`proceed without claiming trailer - ${viewport}`, () => {
      userId = '88007896-6ae7-4600-866a-873d3bc67fd3'; // actualPPMDateZIPAdvanceDone@ppm.approved
      if (isMobile) {
        setMobileViewport();
      }
      signInAndNavigateToAboutPage(userId, false);
    });
  });
});
