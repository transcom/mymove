/* global cy */

// CSRF protection is turned on for all routes.
// We can test with the local dev login that uses POST
describe('testing CSRF protection', function() {
  const csrfForbiddenMsg = 'Forbidden - CSRF token invalid\n';
  const csrfForbiddenRespCode = 403;
  const userId = '6cd03e5b-bee8-4e97-a340-fecb8f3d5465';
  const signInAs = 'tsp';

  it('tests dev login with both unmasked and masked token', function() {
    cy.signInAsUserPostRequest(signInAs, userId);
    cy.contains('Queue: New Shipments');
  });

  it('tests dev login with masked token only', function() {
    cy.signInAsUserPostRequest(signInAs, userId, csrfForbiddenRespCode, csrfForbiddenMsg, false, true);
  });

  it('tests dev login with unmasked token only', function() {
    cy.signInAsUserPostRequest(signInAs, userId, csrfForbiddenRespCode, csrfForbiddenMsg, true, false);
  });

  it('tests dev login without unmasked and masked token', function() {
    cy.signInAsUserPostRequest(signInAs, userId, csrfForbiddenRespCode, csrfForbiddenMsg, false, false);
  });
});
