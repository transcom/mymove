import {
  deleteShipment,
  setMobileViewport,
  signInAndNavigateToFinalCloseoutPage,
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

function verifyFinalIncentiveAndTotals() {
  // TODO: Once we get back final incentive, set a value in the testdatagen func
  //  createMoveWithPPMShipmentReadyForFinalCloseout and check for it here.
  cy.get('h2').contains('Your final estimated incentive: $0.00');

  cy.get('li').contains('4,000 lbs total net weight');

  // TODO: Once we get moving expenses and pro gear back, check for those here as well.

  cy.get('li').contains('0 lbs of pro-gear');
  cy.get('li').contains('$450.00 in expenses claimed');
}
