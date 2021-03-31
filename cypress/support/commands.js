import * as mime from 'mime-types';
import 'cypress-wait-until';
import 'cypress-audit/commands';

import {
  milmoveBaseURL,
  officeBaseURL,
  adminBaseURL,
  milmoveUserType,
  PPMOfficeUserType,
  TOOOfficeUserType,
  TIOOfficeUserType,
  longPageLoadTimeout,
} from './constants';

/**
 * Use this file to define custom commands for Cypress
 * Follow best practices: https://docs.cypress.io/api/cypress-api/custom-commands.html#Best-Practices
 * Prefix commands that make direct calls to the API with 'api' (ie, apiSignInAsUser)
 */

/** Global */
Cypress.Commands.add('prepareCustomerApp', () => {
  Cypress.config('baseUrl', milmoveBaseURL);
});

Cypress.Commands.add('prepareOfficeApp', () => {
  Cypress.config('baseUrl', officeBaseURL);
});

Cypress.Commands.add('prepareAdminApp', () => {
  Cypress.config('baseUrl', adminBaseURL);
});

Cypress.Commands.add('setFeatureFlag', (flagVal, url = '/queues/new') => {
  cy.visit(`${url}?flag:${flagVal}`);
});

// Persist session cookies across multiple tests (use in beforeEach)
Cypress.Commands.add('persistSessionCookies', () => {
  Cypress.Cookies.preserveOnce('masked_gorilla_csrf', 'office_session_token', '_gorilla_csrf');
});

// Use this for issue where Cypress is not clearing cookies between tests
// Delete ALL cookies across domains (milmove, office)
// https://github.com/cypress-io/cypress/issues/781
Cypress.Commands.add('clearAllCookies', () => {
  cy.clearCookies({ domain: null });
});

// Reloads the page but makes an attempt to wait for the loading screen to disappear
Cypress.Commands.add('patientReload', () => {
  cy.reload();
  cy.waitForLoadingScreen();
});

// Visits a given URL but makes an attempt to wait for the loading screen to disappear
Cypress.Commands.add('patientVisit', (url) => {
  cy.visit(url);
  cy.waitForLoadingScreen();
});

// Waits for the loading screen to disappear for a given amount of milliseconds
Cypress.Commands.add('waitForLoadingScreen', (ms = longPageLoadTimeout) => {
  cy.get('h2[data-name="loading-placeholder"]', { timeout: ms }).should('not.exist');
});

// Attempts to double-click a given move locator in a shipment queue list
Cypress.Commands.add('waitForReactTableLoad', () => {
  // Wait for ReactTable loading to be completed
  cy.get('.ReactTable').within(() => {
    cy.get('.-loading.-active', { timeout: longPageLoadTimeout }).should('not.exist');
  });
});

/** Log in */
Cypress.Commands.add('signInAsNewUser', (userType) => {
  cy.visit('/devlocal-auth/login');
  // select the user type and then login as new user
  cy.get('button[data-hook="new-user-login-' + userType + '"]').click();
});

Cypress.Commands.add('signInAsNewMilMoveUser', () => {
  cy.signInAsNewUser(milmoveUserType);
});

Cypress.Commands.add('signInAsNewPPMOfficeUser', () => {
  cy.signInAsNewUser(PPMOfficeUserType);
});

Cypress.Commands.add('signInAsNewTOOUser', () => {
  cy.signInAsNewUser(TOOOfficeUserType);
});

Cypress.Commands.add('signInAsNewTIOUser', () => {
  cy.signInAsNewUser(TIOOfficeUserType);
});

Cypress.Commands.add('signInAsNewAdminUser', () => {
  cy.visit('/devlocal-auth/login');
  // select the user type and then login as new user
  cy.get('button[data-hook="new-user-login-admin"]').click();
});

Cypress.Commands.add('signInAsMultiRoleOfficeUser', () => {
  cy.apiSignInAsUser('9bda91d2-7a0c-4de1-ae02-b8cf8b4b858b', PPMOfficeUserType);
});

Cypress.Commands.add('signIntoOffice', () => {
  cy.apiSignInAsUser('9bfa91d2-7a0c-4de0-ae02-b8cf8b4b858b', PPMOfficeUserType);
  cy.waitForReactTableLoad();
});

// Log in via a direct API request, not the devlocal UI
// Defaults to service member user type, pass in param if signing into Office app
Cypress.Commands.add('apiSignInAsUser', (userId, userType = milmoveUserType) => {
  // This API call is what sets CSRF cookies from the server
  // cy.request('/internal/users/is_logged_in');

  // TODO: Above is not working, I believe because of handling cross-domain cookies/setting baseUrl in between tests
  // https://github.com/cypress-io/cypress/issues/781
  cy.visit('/');

  cy.waitUntil(() => cy.getCookie('masked_gorilla_csrf').then((cookie) => cookie && cookie.value)).then((csrfToken) => {
    cy.request({
      url: '/devlocal-auth/login',
      method: 'POST',
      headers: {
        'X-CSRF-TOKEN': csrfToken,
      },
      body: {
        id: userId,
        userType,
      },
      form: true,
      failOnStatusCode: false,
    }).then((response) => {
      cy.visit('/');
    });
  });
});

// TODO: this is a temporary command for transition of home pages.  Remove when all paths in place
// Log in via a direct API request, not the devlocal UI
// Defaults to service member user type, pass in param if signing into Office app
Cypress.Commands.add('apiSignInAsPpmUser', (userId, userType = milmoveUserType) => {
  // This API call is what sets CSRF cookies from the server
  // cy.request('/internal/users/is_logged_in');

  // TODO: Above is not working, I believe because of handling cross-domain cookies/setting baseUrl in between tests
  // https://github.com/cypress-io/cypress/issues/781
  cy.visit('/ppm');

  cy.waitUntil(() => cy.getCookie('masked_gorilla_csrf').then((cookie) => cookie && cookie.value)).then((csrfToken) => {
    cy.request({
      url: '/devlocal-auth/login',
      method: 'POST',
      headers: {
        'X-CSRF-TOKEN': csrfToken,
      },
      body: {
        id: userId,
        userType,
      },
      form: true,
      failOnStatusCode: false,
    }).then((response) => {
      cy.visit('/ppm');
    });
  });
});

Cypress.Commands.add('logout', () => {
  cy.patientVisit('/');

  cy.getCookie('masked_gorilla_csrf').then((cookie) => {
    cy.request({
      url: '/auth/logout',
      method: 'POST',
      headers: { 'x-csrf-token': cookie.value },
    }).then((resp) => {
      expect(resp.status).to.equal(200);
    });

    // In case of login redirect we once more go to the homepage
    cy.patientVisit('/');
  });
});

/** UI Shortcuts */
// Attempts to double-click a given move locator in a shipment queue list
Cypress.Commands.add('selectQueueItemMoveLocator', (moveLocator) => {
  cy.waitForReactTableLoad();
  cy.get('div').contains(moveLocator).dblclick();
  cy.waitForLoadingScreen();
});

Cypress.Commands.add('nextPage', () => {
  cy.get('[data-testid="wizardNextButton"]').should('be.enabled').click();
});

Cypress.Commands.add('completeFlow', () => {
  cy.get('[data-testid="wizardCompleteButton"]').should('be.enabled').click();
});

Cypress.Commands.add('nextPageAndCheckLocation', (dataCyValue, pageTitle, locationMatch) => {
  const locationRegex = new RegExp(locationMatch); // eslint-disable-line security/detect-non-literal-regexp

  cy.nextPage();
  cy.get(`[data-testid="${dataCyValue}"]`).contains(pageTitle);
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(locationRegex);
  });
});

//from https://github.com/cypress-io/cypress/issues/669
//Cypress doesn't give the right File constructor, so we grab the window's File
Cypress.Commands.add('upload_file', (selector, fileUrl) => {
  const nameSegments = fileUrl.split('/');
  const name = nameSegments[nameSegments.length - 1];
  const rawType = mime.lookup(name);
  // mime returns false if lookup fails
  const type = rawType ? rawType : '';
  return cy.window().then((win) => {
    return cy.fixture(fileUrl, 'base64').then((file) => {
      const blob = Cypress.Blob.base64StringToBlob(file, type);
      const testFile = new win.File([blob], name, { type });
      const event = {};
      event.dataTransfer = new win.DataTransfer();
      event.dataTransfer.items.add(testFile);
      return cy.get(selector).trigger('drop', event);
    });
  });
});

function genericSelect(inputData, fieldName, classSelector) {
  // fieldName is passed as a classname to the react-select component, so select for it if provided
  if (fieldName) {
    classSelector = `${classSelector}.${fieldName}`;
  }
  cy.get(`${classSelector} input[type="text"]`)
    .first()
    .type(`{selectall}{backspace}${inputData}`, { force: true, delay: 75 });

  // Click on the first presented option
  cy.get(classSelector).find('div[class*="option"]').first().click();
}

Cypress.Commands.add('selectDutyStation', (stationName, fieldName) => {
  let classSelector = '.duty-input-box';
  genericSelect(stationName, fieldName, classSelector);
});
