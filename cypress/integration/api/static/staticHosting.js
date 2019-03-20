/* global cy, before */

describe('static file hosting', () => {
  before(() => {
    cy.setupBaseUrl('milmove');
  });

  it('returns the correct content type', () => {
    cy
      .request('/downloads/direct_deposit_form.pdf')
      .its('headers')
      .its('content-type')
      .should('include', 'application/pdf');
  });

  it('returns the correct status', () => {
    cy
      .request('/downloads/direct_deposit_form.pdf')
      .its('status')
      .should('equal', 200);
  });

  it('rejects POST requests', () => {
    let req = cy.request({
      method: 'POST',
      url: '/downloads/direct_deposit_form.pdf',
      failOnStatusCode: false,
    });

    req.its('status').should('equal', 405);
  });
});
