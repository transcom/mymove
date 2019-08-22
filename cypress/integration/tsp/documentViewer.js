import { fileUploadTimeout } from '../../support/constants';
/* global cy */

describe('The document viewer', function() {
  beforeEach(() => {
    // The document viewer is launched in a new tab, so prevent visiting home page first
    cy.signIntoTSP(false);
  });

  it('shows current shipment docs after viewing a shipment with no docs', () => {
    // Find a shipment with no docs
    cy.patientVisit('/shipments/65e00326-420e-436a-89fc-6aeb3f90b870', {
      log: true,
    });

    cy.get('[data-cy="document-upload-link"]')
      .find('a')
      .should('have.attr', 'href')
      .and('contain', '/shipments/65e00326-420e-436a-89fc-6aeb3f90b870/documents/new');

    cy.patientVisit('/queues/approved/', {
      log: true,
    });

    // Find a shipment with a doc
    cy.selectQueueItemMoveLocator('GOTDOC');

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
    });

    cy.get('[data-cy="document-upload-link"]')
      .find('a')
      .should('have.attr', 'href')
      .and('match', /^\/shipments\/[^/]+\/documents\/[^/]+/);
  });

  it('can upload a new document', () => {
    cy.patientVisit('/shipments/65e00326-420e-436a-89fc-6aeb3f90b870/documents', {
      log: true,
    });
    cy.get('[data-cy="document-upload-link"]')
      .find('a')
      .should('have.attr', 'href')
      .and('contain', '/shipments/65e00326-420e-436a-89fc-6aeb3f90b870/documents/new');

    cy.patientVisit('/shipments/65e00326-420e-436a-89fc-6aeb3f90b870/documents/new');

    cy.get('button.submit').should('be.disabled');
    cy.get('select[name="move_document_type"]').select('Other document type');
    cy.get('input[name="title"]').type('super secret info document');
    cy.get('input[name="notes"]').type('burn after reading');
    cy.get('button.submit').should('be.disabled');
    cy.upload_file('.filepond--root', 'top-secret.png');
    cy.get('button.submit', { timeout: fileUploadTimeout })
      .should('not.be.disabled')
      .click();
    cy.get('input[name="title"]').should('be.empty');
    cy.get('input[name="notes"]').should('be.empty');
  });
});
