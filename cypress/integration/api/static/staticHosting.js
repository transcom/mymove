/* global cy, before */

describe('static file hosting', () => {
  before(() => {
    cy.setupBaseUrl('milmove');
  });

  it('returns the correct content type', () => {
    let req = cy.request('/favicon.ico');

    req
      .its('headers')
      .its('content-type')
      .should('include', 'image/x-icon');
  });

  it('returns the correct status', () => {
    let req = cy.request('/favicon.ico');

    req.its('status').should('equal', 200);
  });

  it('rejects POST requests', () => {
    let req = cy.request({
      method: 'POST',
      url: '/favicon.ico',
      failOnStatusCode: false,
    });

    req.its('status').should('equal', 405);
  });
});
