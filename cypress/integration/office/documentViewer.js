/* global cy, Cypress */
describe('The document viewer', function() {
  beforeEach(() => {
    Cypress.config('baseUrl', 'http://officelocal:4000');
  });
  it('redirects to sign in when not logged in', function() {
    cy.visit('/moves/foo/documents');
    cy.contains('Welcome');
    cy.contains('Sign In');
  });
  it('produces error when move cannot be found', () => {
    cy.signIntoOffice();
    cy.visit('/moves/9bfa91d2-7a0c-4de0-ae02-b8cf8b4b858b/documents');
    cy.contains('An error occurred'); //todo: we want better messages when we are making custom call
  });
  it('loads basic information about the move', () => {
    cy.signIntoOffice();
    cy.visit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents');
    cy.contains('In Progress, PPM');
    cy.contains('GBXYUI');
    cy.contains('1617033988');
  });
});

//F2AF74E2-61B0-40AB-9ABD-172A3863E258
