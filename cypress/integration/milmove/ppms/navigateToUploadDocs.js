import { setMobileViewport } from '../../../support/ppmShared';

describe('PPM Request Payment - Begin providing documents flow', () => {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.intercept('GET', '**/internal/moves/**/mto_shipments').as('getShipment');
  });

  const viewportType = [
    { viewport: 'desktop', isMobile: false },
    { viewport: 'mobile', isMobile: true },
  ];

  viewportType.forEach(({ viewport, isMobile }) => {
    it(`has upload documents button enabled - ${viewport}`, () => {
      if (isMobile) {
        setMobileViewport();
      }

      // readyToFinish@ppm.approved
      const userId = 'cde987a1-a717-4a61-98b5-1f05e2e0844d';
      cy.apiSignInAsUser(userId);
      cy.wait('@getShipment');

      verifyHomePageAfterCounseling();
    });
  });
});

function verifyHomePageAfterCounseling() {
  cy.get('h3').should('contain', 'Your move is in progress.');
  cy.get('[data-testid="stepContainer5"]').within(() => {
    cy.get('p').contains('15 Apr 2022');
    cy.get('button').contains('Upload PPM Documents').should('be.enabled').click();
    cy.location().should((loc) => {
      expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/about/);
    });
  });
}
