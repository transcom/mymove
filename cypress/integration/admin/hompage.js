describe('Admin Home Page', function () {
  before(() => {
    cy.prepareAdminApp();
  });

  it('successfully loads when not logged in', function () {
    cy.logout();
    cy.contains('admin.move.mil');
    cy.contains('Sign In');
  });
});
