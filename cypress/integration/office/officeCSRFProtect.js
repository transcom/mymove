import { officeAppName } from '../../support/constants';

/* global cy */

// CSRF protection is turned on for all routes.
// We can test with the local dev login that uses POST
describe('testing CSRF protection', function() {
  const csrfForbiddenMsg = 'Forbidden - CSRF token invalid\n';
  const csrfForbiddenRespCode = 403;
  const userId = '9bfa91d2-7a0c-4de0-ae02-b8cf8b4b858b';

  it('tests dev login with both unmasked and masked token', function() {
    cy.signInAsUserPostRequest(officeAppName, userId);
    cy.contains('Queue: New Moves');
  });

  it('tests dev login with masked token only', function() {
    cy.signInAsUserPostRequest(officeAppName, userId, csrfForbiddenRespCode, csrfForbiddenMsg, false, true);
  });

  it('tests dev login with unmasked token only', function() {
    cy.signInAsUserPostRequest(officeAppName, userId, csrfForbiddenRespCode, csrfForbiddenMsg, true, false);
  });

  it('tests dev login without unmasked and masked token', function() {
    cy.signInAsUserPostRequest(officeAppName, userId, csrfForbiddenRespCode, csrfForbiddenMsg, false, false);
  });
});

describe('testing CSRF protection updating move info', function() {
  it('tests updating user profile with proper tokens', function() {
    cy.signIntoOffice();

    // update info
    cy
      .get('div[class="rt-tr -odd"]')
      .first()
      .dblclick();

    // save info
    cy
      .get('a[class="editable-panel-edit"]')
      .first()
      .click();

    cy
      .get('input[name="orders.orders_number"]')
      .clear()
      .type('CSRF Test')
      .blur();

    cy.get('select[name="orders.orders_type_detail"]').select('HHG_PERMITTED');

    cy
      .get('button[class="usa-button editable-panel-save"]')
      .should('be.enabled')
      .click();

    cy.get('div.orders_number').contains('CSRF Test');

    cy.patientReload();

    cy.contains('CSRF Test');
  });

  it('tests updating user profile without masked token', function() {
    cy.signIntoOffice();

    // update info
    cy
      .get('div[class="rt-tr -odd"]')
      .first()
      .dblclick();

    // save info
    cy
      .get('a[class="editable-panel-edit"]')
      .first()
      .click();

    cy
      .get('input[name="orders.orders_number"]')
      .clear()
      .type('CSRF Protection Failed')
      .blur();

    cy.get('select[name="orders.orders_type_detail"]').select('HHG_PERMITTED');

    // clear cookie
    cy.clearCookie('masked_gorilla_csrf');
    cy.getCookie('masked_gorilla_csrf').should('not.exist');

    cy
      .get('button[class="usa-button editable-panel-save"]')
      .should('be.enabled')
      .click();

    cy.contains('There was an error: Forbidden.');
  });
});
