import * as mime from 'mime-types';
import 'cypress-wait-until';

import {
  milmoveBaseURL,
  officeBaseURL,
  milmoveAppName,
  officeAppName,
  milmoveUserType,
  PPMOfficeUserType,
  TOOOfficeUserType,
  TIOOfficeUserType,
  userTypeToBaseURL,
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

// Call this in your before or beforeEach hook when using cy.route / cy.wait
// https://github.com/cypress-io/cypress/issues/95#issuecomment-347607198
// deletes window.fetch to force fallback to supported XHR
// https://github.com/cypress-io/cypress-example-recipes/tree/master/examples/stubbing-spying__window-fetch
Cypress.Commands.add('removeFetch', () => {
  cy.on('window:before:load', (win) => {
    delete win.fetch;
  });
});

Cypress.Commands.add('setFeatureFlag', (flagVal, url = '/queues/new') => {
  cy.visit(`${url}?flag:${flagVal}`);
});

// Persist session cookies across multiple tests (use in beforeEach)
Cypress.Commands.add('persistSessionCookies', () => {
  Cypress.Cookies.preserveOnce('masked_gorilla_csrf', 'office_session_token', '_gorilla_csrf');
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
  // Visiting the root URL will trigger an is_logged_in API call and set CSRF cookies
  cy.visit('/');

  cy.waitUntil(() => cy.getCookie('masked_gorilla_csrf').then((cookie) => cookie && cookie.value)).then((csrfToken) => {
    cy.request({
      url: '/devlocal-auth/login',
      method: 'POST',
      headers: { 'X-CSRF-TOKEN': csrfToken },
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

// TODO - see if we can remove this (cookies are cleared between tests by default)
Cypress.Commands.add('setBaseUrlAndClearAllCookies', (userType) => {
  [milmoveBaseURL, officeBaseURL].forEach((url) => {
    Cypress.config('baseUrl', url);
    cy.visit('/');
    cy.clearCookies();
  });
  const baseUrl = userTypeToBaseURL[userType]; // eslint-disable-line security/detect-object-injection
  Cypress.config('baseUrl', baseUrl);
  cy.visit('/');
});

Cypress.Commands.add(
  'signInAsUserPostRequest',
  (
    userType,
    userId,
    expectedStatusCode = 200,
    expectedRespBody = null,
    sendGorillaCSRF = true,
    sendMaskedGorillaCSRF = true,
    checkSessionToken = true,
  ) => {
    // setup baseurl
    cy.setBaseUrlAndClearAllCookies(userType);

    // request use to log in
    let sendRequest = (sendRequestUserType, maskedCSRFToken) => {
      cy.request({
        url: '/devlocal-auth/login',
        method: 'POST',
        headers: {
          'X-CSRF-TOKEN': maskedCSRFToken,
        },
        body: {
          id: userId,
          userType: sendRequestUserType,
        },
        form: true,
        failOnStatusCode: false,
      }).then((resp) => {
        cy.visit('/');
        // Default status code to check is 200
        expect(resp.status).to.eq(expectedStatusCode);
        // check response body if needed
        if (expectedRespBody) {
          expect(resp.body).to.eq(expectedRespBody);
        }

        // Login should provide named session tokens
        if (checkSessionToken) {
          // Check that two CSRF cookies and one session cookie exists
          cy.getCookies().should('have.length', 3);
          if (sendRequestUserType === milmoveAppName) {
            cy.getCookie('mil_session_token').should('exist');
            cy.getCookie('office_session_token').should('not.exist');
          } else if (sendRequestUserType === officeAppName) {
            cy.getCookie('mil_session_token').should('not.exist');
            cy.getCookie('office_session_token').should('exist');
          }
        }
      });
    };

    // make sure we log out first before sign in
    cy.logout();
    // GET landing page to get csrf cookies
    cy.request('/');

    // Wait for cookies to be present to make sure the page is fully loaded
    // Otherwise we delete cookies before they exist
    cy.getCookie('_gorilla_csrf').should('exist');
    // Clear out cookies if we don't want to send in request
    if (!sendGorillaCSRF) {
      // Don't include cookie in request header
      cy.clearCookie('_gorilla_csrf');
    }

    if (!sendMaskedGorillaCSRF) {
      // Clear out the masked CSRF token
      cy.clearCookie('masked_gorilla_csrf');
      // Send request without masked token
      sendRequest(userType);
    } else {
      // Send request with masked token
      cy.getCookie('masked_gorilla_csrf').then((cookie) => {
        sendRequest(userType, cookie.value);
      });
    }
  },
);

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
  cy.get('button.next').should('be.enabled').click();
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
    return cy
      .fixture(fileUrl, 'base64')
      .then(Cypress.Blob.base64StringToBlob)
      .then((blob) => {
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
