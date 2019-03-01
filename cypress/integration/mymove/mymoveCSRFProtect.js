import { milmoveAppName } from '../../support/constants';

/* global cy */

// CSRF protection is turned on for all routes.
// We can test with the local dev login that uses POST
describe('testing CSRF protection for dev login', function() {
  const csrfForbiddenMsg = 'Forbidden - CSRF token invalid\n';
  const csrfForbiddenRespCode = 403;
  const userId = '9ceb8321-6a82-4f6d-8bb3-a1d85922a202';

  it('tests dev login with both unmasked and masked token', function() {
    cy.signInAsUserPostRequest(milmoveAppName, userId);
    cy.contains('Move to be scheduled');
    cy.contains('Next Step: Finish setting up your move');
  });

  it('tests dev login with masked token only', function() {
    cy.signInAsUserPostRequest(milmoveAppName, userId, csrfForbiddenRespCode, csrfForbiddenMsg, false, true);
  });

  it('tests dev login with unmasked token only', function() {
    cy.signInAsUserPostRequest(milmoveAppName, userId, csrfForbiddenRespCode, csrfForbiddenMsg, true, false);
  });

  it('tests dev login without unmasked and masked token', function() {
    cy.signInAsUserPostRequest(milmoveAppName, userId, csrfForbiddenRespCode, csrfForbiddenMsg, false, false);
  });
});

describe('testing CSRF protection updating user profile', function() {
  const userId = '9ceb8321-6a82-4f6d-8bb3-a1d85922a202';

  it('tests updating user profile with proper tokens', function() {
    cy.signIntoMyMoveAsUser(userId);

    cy.visit('/moves/review/edit-profile');

    // update info
    cy
      .get('select[name=affiliation]')
      .should('exist')
      .should('have.value', 'ARMY');

    cy
      .get('input[name="middle_name"]')
      .clear()
      .type('CSRF Test')
      .blur();

    // save info
    cy.get('button[type="submit"]').click();

    cy.location().should(loc => {
      expect(loc.pathname).to.match(/^\/$/);
    });

    // reload page
    cy.visit('/moves/review/edit-profile');

    cy
      .get('input[name="middle_name"]')
      .should('exist')
      .should('have.value', 'CSRF Test');
  });

  it('tests updating user profile without masked token', function() {
    cy.signIntoMyMoveAsUser(userId);

    cy.visit('/moves/review/edit-profile');

    // update info
    cy
      .get('input[name="middle_name"]')
      .clear()
      .type('CSRF failed!')
      .blur();

    //clear out masked token
    cy.clearCookie('masked_gorilla_csrf');

    // save info
    cy.get('button[type="submit"]').click();

    cy.get('div[class="usa-alert-text"]').contains('Forbidden');

    // reload page
    cy.visit('/moves/review/edit-profile');

    cy
      .get('input[name="middle_name"]')
      .should('exist')
      .should('not.have.value', 'CSRF failed!');
  });
});
