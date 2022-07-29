import {
  setMobileViewport,
  submitWeightTicketPage,
  signInAndNavigateToWeightTicketPage,
} from '../../../support/ppmCustomerShared';

describe('Weight Tickets', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.intercept('GET', '**/internal/moves/**/mto_shipments').as('getShipment');
    cy.intercept('PATCH', '**/internal/mto-shipments/**').as('patchShipment');
  });

  const viewportType = [
    { viewport: 'desktop', isMobile: false, userId: '88007896-6ae7-4600-866a-873d3bc67fd3' }, // actualPPMDateZIPAdvanceDone@ppm.approved
    { viewport: 'mobile', isMobile: true, userId: '22dba194-3d9a-49c6-8328-718dd945292f' }, // actualPPMDateZIPAdvanceDone2@ppm.approved
  ];

  viewportType.forEach(({ viewport, isMobile, userId }) => {
    it(`proceed with weight ticket documents - ${viewport}`, () => {
      if (isMobile) {
        setMobileViewport();
      }

      signInAndNavigateToWeightTicketPage(userId);
      submitWeightTicketPage();
    });
  });

  const viewportType2 = [
    { viewport: 'desktop', isMobile: false, userId: '9ec731d8-f347-4d34-8b54-4ce9e6ea3282' }, // actualPPMDateZIPAdvanceDone3@ppm.approved
    { viewport: 'mobile', isMobile: true, userId: '2a0146c4-ec9a-4efc-a94c-6c2849c3e167' }, // actualPPMDateZIPAdvanceDone4@ppm.approved
  ];

  viewportType2.forEach(({ viewport, isMobile, userId }) => {
    it(`proceed with claiming trailer - ${viewport}`, () => {
      if (isMobile) {
        setMobileViewport();
      }
      signInAndNavigateToWeightTicketPage(userId);
      submitWeightTicketPage({ hasTrailer: true, ownTrailer: true });
    });
  });

  const viewportType3 = [
    { viewport: 'desktop', isMobile: false, userId: 'bab42ae8-fe0d-4165-87be-dc1317ae0099' }, // actualPPMDateZIPAdvanceDone5@ppm.approved
    { viewport: 'mobile', isMobile: true, userId: '2c4eaae3-5226-456a-94d5-177c679b0656' }, // actualPPMDateZIPAdvanceDone6@ppm.approved
  ];

  viewportType3.forEach(({ viewport, isMobile, userId }) => {
    it(`proceed without claiming trailer - ${viewport}`, () => {
      if (isMobile) {
        setMobileViewport();
      }
      signInAndNavigateToWeightTicketPage(userId);
      submitWeightTicketPage({ hasTrailer: true, ownTrailer: false });
    });
  });
});
