import { tspAppName } from '../../support/constants';

/* global cy */

// CSRF protection is turned on for all routes.
// We can test with the local dev login that uses POST
describe('testing CSRF protection', function() {
  const csrfForbiddenMsg = 'Forbidden - CSRF token invalid\n';
  const csrfForbiddenRespCode = 403;
  const userId = '6cd03e5b-bee8-4e97-a340-fecb8f3d5465';

  it('tests dev login with both unmasked and masked token', function() {
    cy.signInAsUserPostRequest(tspAppName, userId);
    cy.contains('Queue: New Shipments');
  });

  it('tests dev login with masked token only', function() {
    cy.signInAsUserPostRequest(tspAppName, userId, csrfForbiddenRespCode, csrfForbiddenMsg, false, true);
  });

  it('tests dev login with unmasked token only', function() {
    cy.signInAsUserPostRequest(tspAppName, userId, csrfForbiddenRespCode, csrfForbiddenMsg, true, false);
  });

  it('tests dev login without unmasked and masked token', function() {
    cy.signInAsUserPostRequest(tspAppName, userId, csrfForbiddenRespCode, csrfForbiddenMsg, false, false);
  });
});

describe('testing CSRF protection updating shipment info', function() {
  it('tests updating user profile with proper tokens', function() {
    cy.signIntoTSP();

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
      .get('textarea[name="dates.pm_survey_notes"]')
      .clear()
      .type('CSRF Test')
      .blur();

    cy
      .get('button[class="usa-button editable-panel-save"]')
      .should('be.enabled')
      .click();

    cy.patientReload();

    cy.contains('CSRF Test');
  });

  it('tests updating user profile without masked tokens', function() {
    cy.signIntoTSP();

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
      .get('textarea[name="dates.pm_survey_notes"]')
      .clear()
      .type('CSRF failed!')
      .blur();

    // clear cookie
    cy.clearCookie('masked_gorilla_csrf');
    cy.getCookie('masked_gorilla_csrf').should('not.exist');

    cy
      .get('button[class="usa-button editable-panel-save"]')
      .should('be.enabled')
      .click();

    cy.patientReload();

    // No error pops up so we check the value
    cy
      .get('div[class="panel-field pm_survey_notes notes"]')
      .should('exist')
      .should('not.contain', 'CSRF failed!');
  });
});
