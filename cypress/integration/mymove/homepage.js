describe('The Home Page', function () {
  before(() => {
    cy.prepareCustomerApp();
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
    // cy.visit('/ppm');
    cy.get('[data-testid=contact-footer]').contains('Contact Us');
    cy.get('address').within(() => {
      cy.get('a').should('have.attr', 'href', 'https://move.mil/customer-service');
    });
  });
});
