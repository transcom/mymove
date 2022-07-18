import {
  setMobileViewport,
  signInAndNavigateToAboutPage,
  signInAndNavigateToWeightTicketPage,
} from '../../../support/ppmCustomerShared';

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
    it(`validation errors - ${viewport}`, () => {
      userId = '88007896-6ae7-4600-866a-873d3bc67fd3'; // actualPPMDateZIPAdvanceDone@ppm.approved
      if (isMobile) {
        setMobileViewport();
      }

      invalidInputs();
    });
  });

  viewportType.forEach(({ viewport, isMobile, userId }) => {
    it(`proceed with weight ticket documents - ${viewport}`, () => {
      userId = '88007896-6ae7-4600-866a-873d3bc67fd3'; // actualPPMDateZIPAdvanceDone@ppm.approved
      if (isMobile) {
        setMobileViewport();
      }

      signInAndNavigateToWeightTicketPage(userId);
    });
  });

  viewportType.forEach(({ viewport, isMobile, userId }) => {
    it(`proceed without weight tickets (constructed weight) - ${viewport}`, () => {
      userId = '88007896-6ae7-4600-866a-873d3bc67fd3'; // actualPPMDateZIPAdvanceDone@ppm.approved
      if (isMobile) {
        setMobileViewport();
      }
      signInAndNavigateToWeightTicketPage(userId, { useConstructedWeight: true });
    });
  });

  viewportType.forEach(({ viewport, isMobile, userId }) => {
    it(`proceed with claiming trailer - ${viewport}`, () => {
      userId = '88007896-6ae7-4600-866a-873d3bc67fd3'; // actualPPMDateZIPAdvanceDone@ppm.approved
      if (isMobile) {
        setMobileViewport();
      }
      signInAndNavigateToWeightTicketPage(userId, { hasTrailer: true, ownTrailer: true });
    });
  });

  viewportType.forEach(({ viewport, isMobile, userId }) => {
    it(`proceed without claiming trailer - ${viewport}`, () => {
      userId = '88007896-6ae7-4600-866a-873d3bc67fd3'; // actualPPMDateZIPAdvanceDone@ppm.approved
      if (isMobile) {
        setMobileViewport();
      }
      signInAndNavigateToWeightTicketPage(userId, { hasTrailer: true, ownTrailer: false });
    });
  });
});

function invalidInputs() {
  // missing required vehicle description
  cy.get('input[name="vehicleDescription"]').clear().blur();
  cy.get('[class="usa-error-message"]').as('errorMessage');
  cy.get('@errorMessage').contains('Required');
  cy.get('@errorMessage').next('div').get('input').should('have.id', 'vehicleDescription');
  cy.get('input[name="vehicleDescription"]').clear().type('kia forte').blur();
  cy.get('@errorMessage').should('not.exist');

  // missing required empty weight
  cy.get('input[name="emptyWeight"]').clear().blur();
  cy.get('@errorMessage').contains('Required');
  cy.get('@errorMessage').next('div').get('input').should('have.id', 'emptyWeight');
  cy.get('input[name="emptyWeight"]').clear().type(500).blur();
  cy.get('@errorMessage').should('not.exist');

  // missing required full weight
  cy.get('input[name="fullWeight"]').clear().blur();
  cy.get('@errorMessage').contains('Required');
  cy.get('@errorMessage').next('div').get('input').should('have.id', 'fullWeight');
  cy.get('input[name="fullWeight"]').clear().type(700).blur();
  cy.get('@errorMessage').should('not.exist');

  // full weight must be greater than empty weight
  cy.get('input[name="emptyWeight"]').clear().type(1000).blur();
  cy.get('input[name="fullWeight"]').clear().type(500).blur();
  cy.get('@errorMessage').contains('The full weight must be greater than the empty weight');
  cy.get('@errorMessage').next('div').get('input').should('have.id', 'fullWeight');
  cy.get('input[name="fullWeight"]').clear().type(5000).blur();
  cy.get('@errorMessage').should('not.exist');
}
