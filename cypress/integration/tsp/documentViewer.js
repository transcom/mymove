import { fileUploadTimeout } from '../../support/constants';
/* global cy, Cypress */

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
    // cy.visit('/shipments/3ec119bd-1aba-433f-8ad2-c8939f4a8676/documents/new');
    // cy.contains('Upload a new document');
    // cy.get('input[name="title"]').type('super secret info document');
    // cy.get('button.submit').should('be.disabled');
    // cy.get('select[name="move_document_type"]').select('Other document type');
    // cy.get('input[name="notes"]').type('burn after reading');
    // cy.get('button.submit').should('be.disabled');
    // cy.upload_file('.filepond--root', 'top-secret.png');
    // cy
    //   .get('button.submit', { timeout: fileUploadTimeout })
    //   .should('not.be.disabled')
    //   .click();
    // cy.get('input[name="title"]').should('be.empty');
    // cy.get('select[name="move_document_type"]').should('be.empty');
    // cy.get('input[name="notes"]').should('be.empty');
  });
});
