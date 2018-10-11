import { fileUploadTimeout } from '../../support/constants';
/* global cy */

describe('The document viewer', function() {
  beforeEach(() => {
    cy.signIntoTSP();
  });

  it('has a new document link', () => {
    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/queues\/new/);
    });

    // Find a shipment and open it
    cy
      .get('div')
      .contains('DOCVWR')
      .dblclick();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
    });

    cy.get('.documents-list-header').within(() => {
      cy
        .get('a')
        .should('have.attr', 'href')
        .and('match', /^\/shipments\/[^/]+\/documents\/new/);
    });
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
      .get('button.submit', { timeout: fileUploadTimeout })
      .should('not.be.disabled')
      .click();
  });
});
