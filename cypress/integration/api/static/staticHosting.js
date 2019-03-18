/* global cy, before */

describe('static file hosting', () => {
  before(() => {
    cy.setupBaseUrl('milmove');
  });

  it('returns the correct content type', () => {
    cy
      .request('/swagger-ui/internal.html')
      .its('headers')
      .its('content-type')
      .should('include', 'text/html');
  });

  it('returns the correct status', () => {
    cy
      .request('/swagger-ui/internal.html')
      .its('status')
      .should('equal', 200);
  });

  it('rejects POST requests', () => {
    let req = cy.request({
      method: 'POST',
      url: '/swagger-ui/internal.html',
      failOnStatusCode: false,
    });

    req.its('status').should('equal', 405);
  });
});
