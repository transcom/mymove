import { deleteShipment, setMobileViewport, signInAndNavigateToAboutPage } from '../../../support/ppmCustomerShared';

describe('About Your PPM', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.intercept('GET', '**/internal/moves/**/mto_shipments').as('getShipment');
    cy.intercept('PATCH', '**/internal/mto-shipments/**').as('patchShipment');
  });

  const viewportType = [
    { viewport: 'desktop', isMobile: false, userId: 'c28b2eb1-975f-49f7-b8a3-c7377c0da908', selectAdvance: true }, // readyToFinish2@ppm.approved
    { viewport: 'mobile', isMobile: true, userId: '539af373-9474-49f3-b06b-bc4b4d4111de', selectAdvance: true }, // readyToFinish3@ppm.approved
    { viewport: 'desktop', isMobile: false, userId: 'c48998dc-8f93-437a-bd0c-2c0b187b12cb', selectAdvance: false }, // readyToFinish4@ppm.approved
    { viewport: 'mobile', isMobile: true, userId: '62e20f62-638f-4390-bbc0-c672cd7fd2e3', selectAdvance: false }, // readyToFinish5@ppm.approved
  ];

  viewportType.forEach(({ viewport, isMobile, userId, selectAdvance }) => {
    it(`can submit actual PPM shipment info ${selectAdvance ? 'with' : 'without'} an advance - ${viewport}`, () => {
      if (isMobile) {
        setMobileViewport();
      }

      signInAndNavigateToAboutPage(userId, selectAdvance);
    });
  });
});
