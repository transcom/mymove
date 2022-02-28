import { fileUploadTimeout, longPageLoadTimeout } from '../../support/constants';

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

  describe('When user is logged in', function () {
    beforeEach(() => {
      cy.clearAllCookies();
      cy.signIntoOffice();
    });

    it('produces error when move cannot be found', () => {
      cy.patientVisit('/moves/9bfa91d2-7a0c-4de0-ae02-b90988cf8b4b858b/documents');
      cy.contains('An error occurred'); //todo: we want better messages when we are making custom call
    });
    it('loads basic information about the move', () => {
      cy.patientVisit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents');
      cy.contains('In Progress, PPM');
      cy.contains('GBXYUI');
      cy.contains('1617033988');
    });
    it('can upload a new document', () => {
      cy.patientVisit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents');
      cy.get('[data-testid="document-upload-link"]')
        .find('a')
        .should('have.attr', 'href')
        .and('contain', '/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents/new');
      cy.get('[data-testid="document-upload-link"]');
      cy.patientVisit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents/new');

      cy.contains('Upload a new document');
      cy.get('button.submit').should('be.disabled');
      cy.get('input[name="title"]').type('super secret info document');
      cy.get('select[name="move_document_type"]').select('Other document type');
      cy.get('input[name="notes"]').type('burn after reading');
      cy.get('button.submit').should('be.disabled');

      cy.upload_file('.filepond--root', 'sample-orders.png');
      cy.get('button.submit', { timeout: fileUploadTimeout }).should('not.be.disabled').click();
      cy.contains('super secret info document', { timeout: longPageLoadTimeout }).should('have.attr', 'href');
    });
    it('can upload a weight ticket set', () => {
      cy.patientVisit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents/new');
      cy.contains('Upload a new document');
      cy.get('button.submit').should('be.disabled');
      cy.get('input[name="title"]').type('Weight ticket document');
      cy.get('select[name="move_document_type"]').select('Weight ticket set');
      cy.get('input[name="notes"]').type('burn after reading');
      cy.get('button.submit').should('be.disabled');

      cy.upload_file('.filepond--root', 'sample-orders.png');
      cy.get('button.submit', { timeout: fileUploadTimeout }).should('not.be.disabled').click();
      cy.contains('Weight ticket document', { timeout: 2 * longPageLoadTimeout }).should('have.attr', 'href');
    });

    it('can edit an uploaded weight ticket set', () => {
      cy.patientVisit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents');
      cy.get('[data-testid="doc-link"]')
        .find('a')
        .contains('Weight ticket document')
        .should('have.attr', 'href')
        .then((href) => {
          cy.patientVisit(href);
        });
      cy.contains('Details').click();

      cy.contains('Edit').click();

      cy.get('select[name="moveDocument.weight_ticket_set_type"]').select('CAR');
      cy.get('input[name="moveDocument.vehicle_make"]').type('Herbie');
      cy.get('input[name="moveDocument.vehicle_model"]').type('Hotrod');
      cy.get('input[name="moveDocument.empty_weight"]').type('1000');
      cy.get('input[name="moveDocument.full_weight"]').type('2000');

      cy.get('select[name="moveDocument.status"]').select('OK');

      cy.get('button').contains('Save').should('not.be.disabled').click();

      cy.contains('Car');
      cy.contains('Herbie');
      cy.contains('Hotrod');
      cy.contains('1,000');
      cy.contains('2,000');
    });

    it('shows the newly uploaded document in the document list tab', () => {
      cy.patientVisit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents');
      cy.contains('All Documents (2)');
      cy.get('.panel-field')
        .find('a')
        .contains('super secret info document')
        .should('have.attr', 'href')
        .and('match', /^\/moves\/[^/]+\/documents\/[^/]+/);
    });
    it('can upload an expense document', () => {
      cy.patientVisit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents/new');
      cy.contains('Upload a new document');
      cy.get('button.submit').should('be.disabled');
      cy.get('select[name="move_document_type"]').select('Expense');
      cy.get('input[name="title"]').type('expense document');
      cy.get('select[name="moving_expense_type"]').select('Contracted expense');
      cy.get('input[name="requested_amount_cents"]').type('4,000.92');
      cy.get('select[name="payment_method"]').select('Other account');

      cy.get('button.submit').should('be.disabled');

      cy.upload_file('.filepond--root', 'sample-orders.png');
      cy.get('button.submit', { timeout: fileUploadTimeout }).should('not.be.disabled').click();
      cy.contains('expense document', { timeout: longPageLoadTimeout }).should('have.attr', 'href');
    });
    it('can select and update newly-uploaded expense document', () => {
      cy.patientVisit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents');
      cy.get('[data-testid="doc-link"]', { timeout: 2 * longPageLoadTimeout })
        .find('a')
        .contains('expense document')
        .should('have.attr', 'href')
        .and('match', /^\/moves\/[^/]+\/documents\/[^/]+/)
        .then((href) => {
          cy.patientVisit(href);
        });

      cy.contains('Details').click();

      // Verify values have been stored correctly
      cy.contains('4,000.92');
      cy.contains('Contracted expense');
      cy.contains('Other account');

      // Edit payment method
      cy.contains('Edit').click();
      cy.get('select[name="moveDocument.payment_method"]').select('GTCC');
      cy.get('select[name="moveDocument.status"]').select('OK');

      cy.get('button').contains('Save').should('not.be.disabled').click();

      cy.contains('GTCC');
    });
    it('can update expense document to other doc type', () => {
      cy.patientVisit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents');
      cy.get('[data-testid="doc-link"]')
        .find('a')
        .contains('expense document')
        .should('have.attr', 'href')
        .and('match', /^\/moves\/[^/]+\/documents\/[^/]+/)
        .then((href) => {
          cy.patientVisit(href);
        });
      cy.contains('Details').click();
      cy.contains('GTCC');
      cy.contains('Edit').click();

      cy.get('select[name="moveDocument.move_document_type"]').select('Other document type');

      cy.get('button').contains('Save').should('not.be.disabled').click();

      cy.contains('Other document type');

      cy.get('.field-title').contains('Expense Type').should('not.exist');
    });
    it('can update other document type back to expense type', () => {
      cy.patientVisit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents');
      cy.get('[data-testid="doc-link"]')
        .find('a')
        .contains('expense document')
        .should('have.attr', 'href')
        .and('match', /^\/moves\/[^/]+\/documents\/[^/]+/)
        .then((href) => {
          cy.patientVisit(href);
        });

      cy.contains('Details').click();
      cy.contains('OK');
      cy.contains('Edit').click();

      cy.get('select[name="moveDocument.move_document_type"]').select('Expense');
      cy.get('select[name="moveDocument.moving_expense_type"]').select('Contracted expense');
      cy.get('input[name="moveDocument.requested_amount_cents"]').clear().type('4,999.92');
      cy.get('select[name="moveDocument.payment_method"]').select('GTCC');
      cy.get('select[name="moveDocument.status"]').select('OK');

      cy.get('button').contains('Save').should('not.be.disabled').click();

      cy.contains('Expense');
      cy.contains('4,999.92');
      cy.contains('GTCC');
    });
  });
});
