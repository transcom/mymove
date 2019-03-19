import { fileUploadTimeout, officeAppName } from '../../support/constants';

/* global cy */
describe('The document viewer', function() {
  describe('When not logged in', function() {
    beforeEach(() => {
      cy.setupBaseUrl(officeAppName);
      cy.logout();
    });
    it('shows page not found', function() {
      cy.patientVisit('/moves/foo/documents');
      cy.contains('Welcome');
      cy.contains('Sign In');
    });
  });

  describe('When user is logged in', function() {
    beforeEach(() => {
      // The document viewer is launched in a new tab, so pass in false to prevent visiting home page first
      cy.signIntoOffice(false);
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
      cy.patientVisit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents/new');
      cy.contains('Upload a new document');
      cy.get('button.submit').should('be.disabled');
      cy.get('input[name="title"]').type('super secret info document');
      cy.get('select[name="move_document_type"]').select('Other document type');
      cy.get('input[name="notes"]').type('burn after reading');
      cy.get('button.submit').should('be.disabled');

      cy.upload_file('.filepond--root', 'top-secret.png');
      cy
        .get('button.submit', { timeout: fileUploadTimeout })
        .should('not.be.disabled')
        .click();
    });
    it('shows the newly uploaded document in the document list tab', () => {
      cy.patientVisit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents');
      cy.contains('All Documents (1)');
      cy.contains('super secret info document');
      cy
        .get('.pad-ns')
        .find('a')
        .should('have.attr', 'href')
        .and('match', /^\/moves\/[^/]+\/documents\/[^/]+/);
    });
    it('can upload an expense document', () => {
      cy.patientVisit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents/new');
      cy.contains('Upload a new document');
      cy.get('button.submit').should('be.disabled');
      cy.get('select[name="move_document_type"]').select('Expense');
      cy.get('input[name="title"]').type('expense document');
      cy.get('select[name="moving_expense_type"]').select('Contracted Expense');
      cy.get('input[name="requested_amount_cents"]').type('4,000.92');
      cy.get('select[name="payment_method"]').select('Other account');

      cy.get('button.submit').should('be.disabled');

      cy.upload_file('.filepond--root', 'top-secret.png');
      cy
        .get('button.submit', { timeout: fileUploadTimeout })
        .should('not.be.disabled')
        .click();
    });
    it('can select and update newly-uploaded expense document', () => {
      cy.patientVisit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents');
      cy.contains('expense document');
      cy
        .get('.pad-ns')
        .find('a')
        .should('have.attr', 'href')
        .and('match', /^\/moves\/[^/]+\/documents\/[^/]+/);

      cy.contains('expense document').click();
      cy.contains('Details').click();

      // Verify values have been stored correctly
      cy.contains('4,000.92');
      cy.contains('Contracted Expense');
      cy.contains('Other account');

      // Edit payment method
      cy.contains('Edit').click();
      cy.get('select[name="moveDocument.payment_method"]').select('GTCC');
      cy.get('select[name="moveDocument.status"]').select('OK');

      cy
        .get('button')
        .contains('Save')
        .should('not.be.disabled')
        .click();

      cy.contains('GTCC');
    });
    it('can update expense document to other doc type', () => {
      cy.patientVisit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents');
      cy.contains('expense document');
      cy
        .get('.pad-ns')
        .find('a')
        .should('have.attr', 'href')
        .and('match', /^\/moves\/[^/]+\/documents\/[^/]+/);

      cy.contains('expense document').click();
      cy.contains('Details').click();
      cy.contains('GTCC');
      cy.contains('Edit').click();

      cy.get('select[name="moveDocument.move_document_type"]').select('Other document type');

      cy
        .get('button')
        .contains('Save')
        .should('not.be.disabled')
        .click();

      cy.contains('Other document type');

      cy
        .get('.field-title')
        .contains('Expense Type')
        .should('not.exist');
    });
    it('can update other document type back to expense type', () => {
      cy.patientVisit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents');
      cy.contains('expense document');
      cy
        .get('.pad-ns')
        .find('a')
        .should('have.attr', 'href')
        .and('match', /^\/moves\/[^/]+\/documents\/[^/]+/);

      cy.contains('expense document').click();
      cy.contains('Details').click();
      cy.contains('OK');
      cy.contains('Edit').click();

      cy.get('select[name="moveDocument.move_document_type"]').select('Expense');
      cy.get('select[name="moveDocument.moving_expense_type"]').select('Contracted Expense');
      cy.get('input[name="moveDocument.requested_amount_cents"]').type('4,999.92');
      cy.get('select[name="moveDocument.payment_method"]').select('GTCC');
      cy.get('select[name="moveDocument.status"]').select('OK');

      cy
        .get('button')
        .contains('Save')
        .should('not.be.disabled')
        .click();

      cy.contains('Expense');
      cy.contains('4,999.92');
      cy.contains('GTCC');
    });
    it('can upload documents to an HHG move', () => {
      cy.patientVisit('moves/533d176f-0bab-4c51-88cd-c899f6855b9d/documents/new');

      cy.contains('Upload a new document');
      cy.get('button.submit').should('be.disabled');
      cy.get('select[name="move_document_type"]').select('Expense');
      cy.get('input[name="title"]').type('expense document');
      cy.get('select[name="moving_expense_type"]').select('Contracted Expense');
      cy.get('input[name="requested_amount_cents"]').type('4,000.92');
      cy.get('select[name="payment_method"]').select('Other account');

      cy.get('button.submit').should('be.disabled');

      cy.upload_file('.filepond--root', 'top-secret.png');
      cy
        .get('button.submit', { timeout: fileUploadTimeout })
        .should('not.be.disabled')
        .click();

      cy.contains('All Documents (1)');

      cy.get('select[name="move_document_type"]').select('Weight ticket');
      cy.get('input[name="title"]').type('Wait ticket');
      cy.get('input[name="notes"]').type('wait for this document');
      cy.upload_file('.filepond--root', 'top-secret.png');
      cy
        .get('button.submit', { timeout: fileUploadTimeout })
        .should('not.be.disabled')
        .click();

      cy.contains('All Documents (2)');
    });
    it('can navigate to the shipment info page and show line item info', () => {
      cy.patientVisit('/moves/fb4105cf-f5a5-43be-845e-d59fdb34f31c/documents/new', {
        log: true,
      });

      cy.patientVisit('/');

      cy.location().should(loc => {
        expect(loc.pathname).to.match(/^\/queues\/new/);
      });

      cy
        .get('div')
        .contains('Delivered HHGs')
        .click();

      cy.selectQueueItemMoveLocator('DOOB');

      cy.location().should(loc => {
        expect(loc.pathname).to.match(/^\/queues\/new\/moves\/[^/]+\/basics/);
      });

      cy
        .get('a')
        .contains('HHG')
        .click(); // navtab

      cy
        .get('.invoice-panel')
        .get('.invoice-panel-table-cont')
        .find('tbody tr')
        .should(rows => {
          expect(rows).to.have.length.of.at.least(3);
        });
    });
  });
});
