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
  it('can upload an expense document', () => {
    cy.visit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents/new');
    cy.contains('Upload a new document');
    cy.get('button.submit').should('be.disabled');
    cy.get('select[name="move_document_type"]').select('Expense');
    cy.get('input[name="title"]').type('expense document');
    cy.get('select[name="moving_expense_type"]').select('Contracted Expense');
    cy.get('input[name="reimbursement.requested_amount"]').type('12');
    cy.get('select[name="reimbursement.method_of_receipt"]').select('MilPay');

    cy.get('button.submit').should('be.disabled');

    cy.upload_file('.filepond--root', 'top-secret.png');
    cy
      .get('button.submit')
      .should('not.be.disabled')
      .click();
  });
  it('can select and update newly-uploaded expense document', () => {
    cy.visit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents');
    cy.contains('expense document');
    cy
      .get('.pad-ns')
      .find('a')
      .should('have.attr', 'href')
      .and('match', /^\/moves\/[^/]+\/documents\/[^/]+/);

    cy.contains('expense document').click();
    cy.contains('Details').click();
    cy.contains('MilPay');
    cy.contains('Edit').click();

    cy
      .get('select[name="moveDocument.reimbursement.method_of_receipt"]')
      .select('GTCC');
    cy.get('select[name="moveDocument.status"]').select('OK');

    cy
      .get('button')
      .contains('Save')
      .should('not.be.disabled')
      .click();

    cy.contains('GTCC');
  });
  it('can update expense document to other doc type', () => {
    cy.visit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents');
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

    cy
      .get('select[name="moveDocument.move_document_type"]')
      .select('Other document type');

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
    cy.visit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents');
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
    cy
      .get('select[name="moveDocument.moving_expense_type"]')
      .select('Contracted Expense');
    cy
      .get('input[name="moveDocument.reimbursement.requested_amount"]')
      .type('12');
    cy
      .get('select[name="moveDocument.reimbursement.method_of_receipt"]')
      .select('GTCC');
    cy.get('select[name="moveDocument.status"]').select('OK');

    cy
      .get('button')
      .contains('Save')
      .should('not.be.disabled')
      .click();

    cy.contains('Expense');
    cy.contains('GTCC');
  });
});

//F2AF74E2-61B0-40AB-9ABD-172A3863E258
