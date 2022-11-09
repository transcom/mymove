import {
  navigateFromCloseoutReviewPageToProGearPage,
  setMobileViewport,
  signInAndNavigateToPPMReviewPage,
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
    { viewport: 'mobile', isMobile: true, userId: '7d4dbc69-2973-4c8b-bf75-6fb582d7a5f6' }, // progearWeightTicket2@ppm.approved
  ];
  viewportType.forEach(({ viewport, isMobile, userId }) => {
    it(`page loads - ${viewport}`, () => {
      if (isMobile) {
        setMobileViewport();
      }

      signInAndNavigateToPPMReviewPage(userId);
      navigateFromCloseoutReviewPageToProGearPage();
      cy.get('[data-testid="selfProGear"]').should('not.be.checked');
      cy.get('[data-testid="spouseProGear"]').should('not.be.checked');
      cy.get('label[for="ownerOfProGearSelf"').click();
      cy.get('[data-testid="selfProGear"]').should('be.checked');
      cy.get('[data-testid="spouseProGear"]').should('not.be.checked');
      cy.get('label[for="ownerOfProGearSpouse"').click();
      cy.get('[data-testid="selfProGear"]').should('not.be.checked');
      cy.get('[data-testid="spouseProGear"]').should('be.checked');
      cy.a11y();
    });
  });
});
