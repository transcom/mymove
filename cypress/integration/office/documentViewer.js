/* global cy, Cypress */
describe('The document viewer', function() {
  beforeEach(() => {
    cy.signIntoOffice();
  });
  it('redirects to sign in when not logged in', function() {
    cy.contains('Sign Out').click();
    cy.visit('/moves/foo/documents');
    cy.contains('Welcome');
    cy.contains('Sign In');
  });
  it('produces error when move cannot be found', () => {
    cy.visit('/moves/9bfa91d2-7a0c-4de0-ae02-b90988cf8b4b858b/documents');
    cy.contains('An error occurred'); //todo: we want better messages when we are making custom call
  });
  it('loads basic information about the move', () => {
    cy.visit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents');
    cy.contains('In Progress, PPM');
    cy.contains('GBXYUI');
    cy.contains('1617033988');
  });
  it('can upload a new document', () => {
    cy.visit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents/new');
    cy.contains('Upload a new document');
    cy.get('button.submit').should('be.disabled');
    cy.get('input[name="title"]').type('super secret info document');
    cy.get('select[name="move_document_type"]').select('Other document type');
    cy.get('input[name="notes"]').type('burn after reading');
    cy.get('button.submit').should('be.disabled');

    cy.upload_file('.filepond--root', 'top-secret.png');
    cy
      .get('button.submit')
      .should('not.be.disabled')
      .click();
    // TODO: add tests for uploaded document viewer
  });
  it('shows the newly uploaded document in the document list tab', () => {
    cy.visit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents');
    cy.contains('All Documents (1)');
    cy.contains('super secret info document');
    cy
      .get('.pad-ns')
      .find('a')
      .should('have.attr', 'href')
      .and('match', /^\/moves\/[^/]+\/documents\/[^/]+/);
  });
});

//F2AF74E2-61B0-40AB-9ABD-172A3863E258
