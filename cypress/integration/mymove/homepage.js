describe('The Home Page', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.intercept('GET', '**/internal/moves/**/mto_shipments').as('getShipment');
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
    { canEditOrDelete: true, moveSubmitted: false, userID: '1b16773e-995b-4efe-ad1c-bef2ae1253f8' }, // finished@ppm.unsubmitted
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
        cy.get('h3').should('contain', 'Next: Talk to a move counselor');
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
});
