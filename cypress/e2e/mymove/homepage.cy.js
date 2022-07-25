import { deleteShipment, setMobileViewport } from '../../support/ppmCustomerShared';

describe('The Home Page', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.intercept('GET', '**/internal/moves/**/mto_shipments').as('getShipment');
    cy.intercept('GET', '**/internal/moves/**/mto_shipments').as('getShipment');
    cy.intercept('DELETE', '**/internal/mto-shipments/**').as('deleteShipment');
  });

  it('passes a pa11y audit', function () {
    cy.visit('/');
    cy.pa11y();
  });

  it('creates new devlocal user', function () {
    cy.signInAsNewMilMoveUser();
  });

  it('successfully loads when not logged in', function () {
    cy.logout();
    cy.contains('Welcome');
    cy.contains('Sign in');
  });

  it('contains the link to customer service', function () {
    cy.get('[data-testid=contact-footer]').contains('Contact Us');
    cy.get('address').within(() => {
      cy.get('a').should('have.attr', 'href', 'https://move.mil/customer-service');
    });
  });

  const editTestCases = [
    { canEditOrDelete: true, moveSubmitted: false, userID: '1b16773e-995b-4efe-ad1c-bef2ae1253f8' }, // full@ppm.unsubmitted
    { canEditOrDelete: false, moveSubmitted: true, userID: '2d6a16ec-c031-42e2-aa55-90a1e29b961a' }, // new@ppm.submitted
  ];

  editTestCases.forEach(({ canEditOrDelete, moveSubmitted, userID }) => {
    const testTitle = `${canEditOrDelete ? 'can' : "can't"} edit/delete the shipment when move ${
      moveSubmitted ? 'is' : "isn't"
    } submitted`;

    it(testTitle, () => {
      cy.apiSignInAsUser(userID);

      cy.wait('@getShipment');

      if (moveSubmitted) {
        cy.get('h3').should('contain', 'Next step: Your move gets approved');
      } else {
        cy.get('h3').should('contain', 'Time to submit your move');
      }

      if (canEditOrDelete) {
        cy.get('[data-testid="shipment-list-item-container"] button').contains('Edit');
        cy.get('[data-testid="shipment-list-item-container"] button').contains('Delete');
      } else {
        cy.get('[data-testid="shipment-list-item-container"] button').should('not.exist');
        cy.get('[data-testid="shipment-list-item-container"] button').should('not.exist');
      }
    });
  });

  const viewportType = [
    { viewport: 'desktop', isMobile: false, userId: '57d58062-93ac-4eb7-b1da-21dd137e4f65' }, // deleteShipment@ppm.unsubmitted
    { viewport: 'mobile', isMobile: true, userId: '781cf194-4eb2-4def-9da6-01abdc62333d' }, // deleteShipmentMobile@ppm.unsubmitted
  ];

  viewportType.forEach(({ viewport, isMobile, userId }) => {
    it(`deletes shipment - ${viewport}`, () => {
      if (isMobile) {
        setMobileViewport();
      }

      navigateDeletingShipment(userId, isMobile);
    });
  });
});

function navigateDeletingShipment(userId, isMobile = false) {
  cy.apiSignInAsUser(userId);
  cy.wait('@getShipment');
  cy.get('[data-testid="shipment-list-item-container"]').as('shipmentListContainer');
  deleteShipment('@shipmentListContainer', 0);
}
