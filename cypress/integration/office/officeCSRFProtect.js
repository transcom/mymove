import { ServicesCounselorOfficeUserType, TOOOfficeUserType } from '../../support/constants';

// CSRF protection is turned on for all routes.
// We can test with the local dev login that uses POST
describe('testing CSRF protection', function () {
  before(() => {
    cy.prepareOfficeApp();
  });

  const csrfForbiddenMsg = 'Forbidden - CSRF token invalid\n';
  const csrfForbiddenRespCode = 403;
  const userId = '9bfa91d2-7a0c-4de0-ae02-b8cf8b4b858b';
  const requestParams = {
    url: '/devlocal-auth/login',
    method: 'POST',
    body: {
      id: userId,
      userType: ServicesCounselorOfficeUserType,
    },
    form: true,
    failOnStatusCode: false,
  };

  it('can successfully dev login with both unmasked and masked token', function () {
    cy.apiSignInAsUser(userId, TOOOfficeUserType);
    cy.contains('All moves');
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

describe('testing CSRF protection updating move info', function () {
  before(() => {
    cy.prepareOfficeApp();
  });

  it('tests updating user profile with proper tokens', function () {
    cy.signIntoOffice();

    // update info
    cy.get('tr[data-uuid]').first().click();

    // save info
    cy.get('[data-testid="edit-orders"]').first().click();

    cy.get('input[name="ordersNumber"]').clear().type('CSRF Test').blur();

    cy.get('select[name="ordersTypeDetail"]').select('HHG_PERMITTED');

    cy.get('button[class="usa-button"]').contains('Save').should('be.enabled').click();

    cy.get('dd[data-testid="ordersNumber"]').contains('CSRF Test');

    cy.patientReload();

    cy.contains('CSRF Test');
  });

  it('tests updating user profile without masked token', function () {
    cy.signIntoOffice();

    // update info
    cy.get('tr[data-uuid]').first().click();

    // save info
    cy.get('[data-testid="edit-orders"]').first().click();

    cy.get('input[name="ordersNumber"]').clear().type('CSRF Protection Failed').blur();

    cy.get('select[name="ordersTypeDetail"]').select('HHG_PERMITTED');

    // clear cookie
    cy.clearCookie('masked_gorilla_csrf');
    cy.getCookie('masked_gorilla_csrf').should('not.exist');

    cy.get('button[class="usa-button"]').contains('Save').should('be.enabled').click();

    // Save fails, remains on update form with button deactivated
    cy.get('button[class="usa-button"]').contains('Save').should('be.disabled');
  });
});
