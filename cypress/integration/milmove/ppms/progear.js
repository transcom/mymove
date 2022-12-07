import {
  setMobileViewport,
  signInAndNavigateToProgearPage,
  submitProgearPage,
} from '../../../support/ppmCustomerShared';

describe('Progear', function () {
  before(() => {
    cy.prepareCustomerApp();
  });
  beforeEach(() => {
    cy.intercept('GET', '**/internal/moves/**/mto_shipments').as('getShipment');
    cy.intercept('PATCH', '**/internal/ppm-shipments/**/pro-gear-weight-tickets/**').as('patchProgearWeightTicket');
    cy.intercept('POST', '**/internal/ppm-shipments/**/uploads**').as('uploadFile');
  });
  const viewportType = [
    { viewport: 'desktop', isMobile: false, userId: '33eabbb6-416d-4d91-ba5b-bfd7d35e3037' }, // progearWeightTicket@ppm.approved
    { viewport: 'mobile', isMobile: true, userId: '7d4dbc69-2973-4c8b-bf75-6fb582d7a5f6' }, // progearWeightTicket2@ppm.approved
  ];
  viewportType.forEach(({ viewport, isMobile, userId }) => {
    it(`page loads - ${viewport}`, () => {
      if (isMobile) {
        setMobileViewport();
      }
      signInAndNavigateToProgearPage(userId);
      submitProgearPage({ belongsToSelf: true });

      cy.wait('@patchProgearWeightTicket');

      cy.get('.progearSection h4').contains('Set 2');
      cy.get('.progearSection dl').eq(-1).as('progearSet').contains('Pro-gear');
      cy.get('@progearSet').contains('Radio equipment');
      // contains won't work when text content is divided between multiple semantic html elements
      cy.get('@progearSet').contains('Weight:');
      cy.get('@progearSet').contains('2,000 lbs');
    });
  });
});
