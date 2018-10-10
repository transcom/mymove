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
});
