import {
  deleteShipment,
  setMobileViewport,
  signInAndNavigateToAboutPageWithAdvance,
  signInAndNavigateToAboutPageWithoutAdvance,
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
    { viewport: 'desktop', isMobile: false, userId: 'c28b2eb1-975f-49f7-b8a3-c7377c0da908' }, // readyToFinish2@ppm.approved
    { viewport: 'mobile', isMobile: true, userId: 'c28b2eb1-975f-49f7-b8a3-c7377c0da908' }, // readyToFinish2@ppm.approved
  ];

  viewportType.forEach(({ viewport, isMobile, userId }) => {
    it(`can upload with an advance - ${viewport}`, () => {
      if (isMobile) {
        setMobileViewport();
      }

      signInAndNavigateToAboutPageWithAdvance(userId);
    });
  });
});

describe('About Your PPM - Without Advances', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.intercept('GET', '**/internal/moves/**/mto_shipments').as('getShipment');
    cy.intercept('PATCH', '**/internal/mto-shipments/**').as('patchShipment');
    cy.intercept('DELETE', '**/internal/mto-shipments/**').as('deleteShipment');
    cy.intercept('GET', '**/internal/moves/**/signed_certifications').as('signedCertifications');
  });
  const viewportType = [
    { viewport: 'desktop', isMobile: false, userId: 'cde987a1-a717-4a61-98b5-1f05e2e0844d' }, // readyToFinish@ppm.approved
    { viewport: 'mobile', isMobile: true, userId: 'cde987a1-a717-4a61-98b5-1f05e2e0844d' }, // readyToFinish@ppm.approved
  ];

  viewportType.forEach(({ viewport, isMobile, userId }) => {
    it(`can upload without an advance - ${viewport}`, () => {
      if (isMobile) {
        setMobileViewport();
      }
      signInAndNavigateToAboutPageWithoutAdvance(userId);
    });
  });
});
