/* global cy */

// CSRF protection is turned on for all routes.
// We can test with the local dev login that uses POST
describe('testing CSRF protection', function() {
  const csrfForbiddenMsg = 'Forbidden - CSRF token invalid\n';
  const csrfForbiddenRespCode = 403;
  const userId = '9ceb8321-6a82-4f6d-8bb3-a1d85922a202';
  const signInAs = 'mymove';

  it('tests dev login with both unmasked and masked token', function() {
    cy.signInAsUserPostRequest(signInAs, userId);
    cy.contains('Move to be scheduled');
    cy.contains('Next Step: Finish setting up your move');
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
