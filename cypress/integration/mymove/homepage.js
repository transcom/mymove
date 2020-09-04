describe('The Home Page', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  it('creates new devlocal user', function () {
    cy.signInAsNewMilMoveUser();
  });

  it('successfully loads when not logged in', function () {
    cy.logout();
    cy.contains('Welcome');
    cy.contains('Sign In');
  });

  it('contains the link to customer service', function () {
    cy.visit('/');
    cy.get('[data-testid=contact-footer]').contains('Contact Us');
    cy.get('address').within(() => {
      cy.get('a').should('have.attr', 'href', 'https://move.mil/customer-service');
    });
  });
});
