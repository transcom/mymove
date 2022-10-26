import { milmoveAppName } from '../../support/constants';

// CSRF protection is turned on for all routes.
// We can test with the local dev login that uses POST
describe('testing CSRF protection for dev login', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  const csrfForbiddenMsg = 'Forbidden - CSRF token invalid\n';
  const csrfForbiddenRespCode = 403;
  const userId = '9ceb8321-6a82-4f6d-8bb3-a1d85922a202';
  const requestParams = {
    url: '/devlocal-auth/login',
    method: 'POST',
    body: {
      id: userId,
      userType: milmoveAppName,
    },
    form: true,
    failOnStatusCode: false,
  };

  it('tests dev login with both unmasked and masked token', function () {
    // sm_no_move_type@example.com
    cy.apiSignInAsUser(userId);
    cy.contains('set up your shipments');

    cy.contains("Share where and when you're moving");
  });

  it('cannot dev login with masked token only', function () {
    cy.request('/internal/users/is_logged_in');
    cy.getCookie('_gorilla_csrf').should('exist');

    // Remove unmasked token
    cy.clearCookie('_gorilla_csrf');

    cy.getCookie('masked_gorilla_csrf')
      .then((cookie) => cookie && cookie.value)
      .then((csrfToken) => {
        cy.request({
          ...requestParams,
          headers: { 'X-CSRF-TOKEN': csrfToken },
        }).then((response) => {
          cy.visit('/');
          expect(response.status).to.eq(csrfForbiddenRespCode);
          expect(response.body).to.eq(csrfForbiddenMsg);
        });
      });
  });

  it('cannot dev login with unmasked token only', function () {
    cy.request('/internal/users/is_logged_in');
    cy.getCookie('_gorilla_csrf').should('exist');

    // Remove masked CSRF token
    cy.clearCookie('masked_gorilla_csrf');

    // Attempt to log in with no X-CSRF-TOKEN
    cy.request({
      ...requestParams,
      headers: { 'X-CSRF-TOKEN': null },
    }).then((response) => {
      cy.visit('/');
      expect(response.status).to.eq(csrfForbiddenRespCode);
      expect(response.body).to.eq(csrfForbiddenMsg);
    });
  });

  it('cannot dev login without unmasked and masked token', function () {
    cy.request('/internal/users/is_logged_in');
    cy.getCookie('_gorilla_csrf').should('exist');

    // Remove both CSRF tokens
    cy.clearCookie('_gorilla_csrf');
    cy.clearCookie('masked_gorilla_csrf');

    // Attempt to log in with no X-CSRF-TOKEN
    cy.request({
      ...requestParams,
      headers: { 'X-CSRF-TOKEN': null },
    }).then((response) => {
      cy.visit('/');
      expect(response.status).to.eq(csrfForbiddenRespCode);
      expect(response.body).to.eq(csrfForbiddenMsg);
    });
  });
});

describe('testing CSRF protection updating user profile', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.intercept('PATCH', '**/internal/service_members/**').as('patchServiceMember');
  });

  it('tests updating user profile with proper tokens', function () {
    // sm_no_move_type@example.com
    const userId = '9ceb8321-6a82-4f6d-8bb3-a1d85922a202';
    cy.apiSignInAsUser(userId);

    cy.visit('/moves/review/edit-profile');

    // update info
    cy.get('select[name=affiliation]').should('exist').should('have.value', 'ARMY');

    cy.get('input[name="middle_name"]').clear().type('CSRF Test').blur();

    // save info
    cy.get('button[type="submit"]').click();
    cy.wait('@patchServiceMember');

    cy.getCookie('_gorilla_csrf').should('exist');
    cy.getCookie('masked_gorilla_csrf').should('exist');

    /*
    cy.location().should((loc) => {
      expect(loc.pathname).to.match(/^\/ppm$/);
    });
    */

    // reload page
    cy.getCookie('_gorilla_csrf').should('exist');
    cy.getCookie('masked_gorilla_csrf').should('exist');

    cy.visit('/moves/review/edit-profile');

    cy.get('input[name="middle_name"]').should('exist').should('have.value', 'CSRF Test');
  });

  it('tests updating user profile without masked token', function () {
    // sm_no_move_type@example.com
    const userId = '9ceb8321-6a82-4f6d-8bb3-a1d85922a202';
    cy.apiSignInAsUser(userId);

    cy.visit('/moves/review/edit-profile');

    // update info
    cy.get('input[name="middle_name"]').clear().type('CSRF failed!').blur();

    //clear out masked token
    cy.clearCookie('masked_gorilla_csrf');

    // save info
    cy.get('button[type="submit"]').click();

    cy.get('div[class="usa-alert__text"]').contains('Forbidden');

    // reload page
    cy.visit('/moves/review/edit-profile');

    cy.get('input[name="middle_name"]').should('exist').should('not.have.value', 'CSRF failed!');
  });
});
