import {
  deleteShipment,
  setMobileViewport,
  signInAndNavigateToFinalCloseoutPage,
  verifyFinalIncentiveAndTotals,
} from '../../../support/ppmCustomerShared';

describe('Final Closeout', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.intercept('GET', '**/internal/moves/**/mto_shipments').as('getShipment');
    cy.intercept('PATCH', '**/internal/mto-shipments/**').as('patchShipment');
  });

  const testCases = [
    { viewport: 'desktop', isMobile: false },
    { viewport: 'mobile', isMobile: true },
  ];

  testCases.forEach(({ viewport, isMobile }) => {
    it(`can see final closeout page with final estimated incentive and shipment totals - ${viewport}`, () => {
      if (isMobile) {
        setMobileViewport();
      }

      signInAndNavigateToFinalCloseoutPage('1c842b03-fc2d-4e92-ade8-bd3e579196e0'); // readyForFinalComplete@ppm.approved

      verifyFinalIncentiveAndTotals();
    });
  });
});
