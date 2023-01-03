import { fileUploadTimeout } from '../../support/constants';

describe('The document viewer', function () {
  describe('When not logged in', function () {
    beforeEach(() => {
      cy.prepareOfficeApp();
      cy.logout();
    });

    it('shows page not found', function () {
      cy.patientVisit('/moves/foo/documents');
      cy.contains('Welcome');
      cy.contains('Sign in');
    });
  });

  // describe('When user is logged in', function () {
  //   beforeEach(() => {
  //     cy.clearAllCookies();
  //     cy.signIntoOffice();
  //   });

  //   it('produces error when move cannot be found', () => {
  //     cy.patientVisit('/moves/9bfa91d2-7a0c-4de0-ae02-b90988cf8b4b858b/documents');
  //     cy.contains('An error occurred'); //todo: we want better messages when we are making custom call
  //   });
  //   it('loads basic information about the move', () => {
  //     cy.patientVisit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents');
  //     cy.contains('In Progress, PPM');
  //     cy.contains('GBXYUI');
  //     cy.contains('1617033988');
  //   });
  //   it('can upload a new document', () => {
  //     cy.patientVisit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents');
  //     cy.get('[data-testid="document-upload-link"]')
  //       .find('a')
  //       .should('have.attr', 'href')
  //       .and('contain', '/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents/new');
  //     cy.get('[data-testid="document-upload-link"]');
  //     cy.patientVisit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents/new');

  //     cy.contains('Upload a new document');
  //     cy.get('button.submit').should('be.disabled');
  //     cy.get('input[name="title"]').type('super secret info document');
  //     cy.get('select[name="move_document_type"]').select('Other document type');
  //     cy.get('input[name="notes"]').type('burn after reading');
  //     cy.get('button.submit').should('be.disabled');

  //     cy.intercept('/internal/uploads').as('uploadFile');
  //     cy.upload_file('.filepond--root', 'sample-orders.png');
  //     cy.wait('@uploadFile');

  //     cy.intercept('/internal/moves/*/move_documents').as('postDocument');
  //     cy.get('button.submit', { timeout: fileUploadTimeout }).should('not.be.disabled').click();
  //     cy.wait('@postDocument');
  //   });
  //   it('can upload a weight ticket set', () => {
  //     cy.patientVisit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents/new');
  //     cy.contains('Upload a new document');
  //     cy.get('button.submit').should('be.disabled');
  //     cy.get('input[name="title"]').type('Weight ticket document');
  //     cy.get('select[name="move_document_type"]').select('Weight ticket set');
  //     cy.get('input[name="notes"]').type('burn after reading');
  //     cy.get('button.submit').should('be.disabled');

  //     cy.intercept('/internal/uploads').as('uploadFile');
  //     cy.upload_file('.filepond--root', 'sample-orders.png');
  //     cy.wait('@uploadFile');

  //     cy.intercept('/internal/moves/*/move_documents').as('postDocument');
  //     cy.get('button.submit', { timeout: fileUploadTimeout }).should('not.be.disabled').click();
  //     cy.wait('@postDocument');
  //   });
  // });
});
