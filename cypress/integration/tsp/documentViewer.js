import { fileUploadTimeout } from '../../support/constants';
/* global cy */

describe('The document viewer', function() {
  beforeEach(() => {
    // The document viewer is launched in a new tab, so prevent visiting home page first
    cy.signIntoTSP(false);
  });

  it('has a new document links', () => {
    cy.patientVisit('/');

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

    cy
      .get('.usa-heading')
      .contains('Documents')
      .within(() => {
        cy
          .get('a')
          .should('have.attr', 'href')
          .and('match', /^\/shipments\/[^/]+\/documents\/new/);
      });

    cy
      .get('.documents > .status')
      .contains('Upload new document')
      .should('have.attr', 'href')
      .and('match', /^\/shipments\/[^/]+\/documents\/new/);
  });

  it('shows current shipment docs after viewing a shipment with no docs', () => {
    // Find a shipment with no docs
    cy.patientVisit('/shipments/65e00326-420e-436a-89fc-6aeb3f90b870', {
      log: true,
    });

    cy
      .get('.documents > .status')
      .contains('Upload new document')
      .should('have.attr', 'href')
      .and('match', /^\/shipments\/[^/]+\/documents\/new/);

    cy.patientVisit('/queues/approved/', {
      log: true,
    });

    // Find a shipment with a doc
    cy
      .get('div')
      .contains('GOTDOC')
      .dblclick();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/shipments\/[^/]+/);
    });

    cy
      .get('.documents > .status')
      .should('have.attr', 'href')
      .and('match', /^\/shipments\/[^/]+\/documents\/[^/]+/);
  });

  it('can upload a new document', () => {
    cy.patientVisit('/shipments/65e00326-420e-436a-89fc-6aeb3f90b870/documents/new', {
      log: true,
    });

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
    cy.get('input[name="title"]').should('be.empty');
    cy.get('input[name="notes"]').should('be.empty');
  });

  it('can navigate to the shipment info page and show line item info', () => {
    cy.patientVisit('/shipments/67a3cbe7-4ae3-4f6a-9f9a-4f312e7458b9/documents/new', {
      log: true,
    });

    cy.get('button.submit').should('be.disabled');

    cy.get('a[title="Home"]').click();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/queues\/new/);
    });

    cy
      .get('div')
      .contains('Delivered Shipments')
      .click();

    cy
      .get('div')
      .contains('DOOB')
      .dblclick();

    cy
      .get('.invoice-panel')
      .get('.invoice-panel-table-cont')
      .find('tbody tr')
      .should(rows => {
        expect(rows).to.have.length.of.at.least(3);
      });
  });
});
