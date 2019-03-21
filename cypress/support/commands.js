import * as mime from 'mime-types';
import { milmoveAppName, officeAppName, tspAppName, longPageLoadTimeout } from './constants';

/* global Cypress, cy */
// ***********************************************
// This example commands.js shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************
//
//
// -- This is a parent command --
// Cypress.Commands.add("login", (email, password) => { ... })
//
//
// -- This is a child command --
// Cypress.Commands.add("drag", { prevSubject: 'element'}, (subject, options) => { ... })
//
//
// -- This is a dual command --
// Cypress.Commands.add("dismiss", { prevSubject: 'optional'}, (subject, options) => { ... })
//
//
// -- This is will overwrite an existing command --
// Cypress.Commands.overwrite("visit", (originalFn, url, options) => { ... })

Cypress.Commands.add('signInAsNewUser', () => {
  // make sure we log out first before sign in
  cy.logout();

  cy.visit('/devlocal-auth/login');
  // should have both our csrf cookie tokens now
  cy.getCookie('_gorilla_csrf').should('exist');
  cy.getCookie('masked_gorilla_csrf').should('exist');
  cy.get('button[data-hook="new-user-login"]').click();
});

Cypress.Commands.add('signIntoMyMoveAsUser', userId => {
  cy.setupBaseUrl(milmoveAppName);
  cy.signInAsUser(userId);
});
Cypress.Commands.add('signIntoOffice', () => {
  cy.setupBaseUrl(officeAppName);
  cy.signInAsUser('9bfa91d2-7a0c-4de0-ae02-b8cf8b4b858b');
  cy.waitForReactTableLoad();
});
Cypress.Commands.add('signIntoTSP', () => {
  cy.setupBaseUrl(tspAppName);
  cy.signInAsUser('6cd03e5b-bee8-4e97-a340-fecb8f3d5465');
  cy.waitForReactTableLoad();
});
Cypress.Commands.add('signInAsUser', userId => {
  // make sure we log out first before sign in
  cy.logout();

  cy.visit('/devlocal-auth/login');
  // should have both our csrf cookie tokens now
  cy.getCookie('_gorilla_csrf').should('exist');
  cy.getCookie('masked_gorilla_csrf').should('exist');
  cy.get('button[value="' + userId + '"]').click();
});

// Reloads the page but makes an attempt to wait for the loading screen to disappear
Cypress.Commands.add('patientReload', () => {
  cy.reload();
  cy.waitForLoadingScreen();
});

// Visits a given URL but makes an attempt to wait for the loading screen to disappear
Cypress.Commands.add('patientVisit', url => {
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

// Attempts to double-click a given move locator in a shipment queue list
Cypress.Commands.add('selectQueueItemMoveLocator', moveLocator => {
  cy.waitForReactTableLoad();

  cy
    .get('div')
    .contains(moveLocator)
    .dblclick();

  cy.waitForLoadingScreen();
});

Cypress.Commands.add('setFeatureFlag', (flagVal, url = '/queues/new') => {
  cy.visit(`${url}?flag:${flagVal}`);
});

Cypress.Commands.add(
  'signInAsUserPostRequest',
  (
    signInAs,
    userId,
    expectedStatusCode = 200,
    expectedRespBody = null,
    sendGorillaCSRF = true,
    sendMaskedGorillaCSRF = true,
  ) => {
    // setup baseurl
    cy.setupBaseUrl(signInAs);

    // request use to log in
    let sendRequest = maskedCSRFToken => {
      cy
        .request({
          url: '/devlocal-auth/login',
          method: 'POST',
          headers: {
            'X-CSRF-TOKEN': maskedCSRFToken,
          },
          body: {
            id: userId,
          },
          form: true,
          failOnStatusCode: false,
        })
        .then(resp => {
          cy.visit('/');
          // Default status code to check is 200
          expect(resp.status).to.eq(expectedStatusCode);
          // check response body if needed
          if (expectedRespBody) {
            expect(resp.body).to.eq(expectedRespBody);
          }
        });
    };

    // make sure we log out first before sign in
    cy.logout();
    // GET landing page to get csrf cookies
    cy.request('/');

    // Clear out cookies if we don't want to send in request
    if (!sendGorillaCSRF) {
      cy.clearCookie('_gorilla_csrf');
      // Don't include cookie in request header
      cy.getCookie('_gorilla_csrf').should('not.exist');
    } else {
      // Include cookie in request header
      cy.getCookie('_gorilla_csrf').should('exist');
    }

    if (!sendMaskedGorillaCSRF) {
      // Clear out the masked CSRF token
      cy.clearCookie('masked_gorilla_csrf');
      // Send request without masked token
      cy
        .getCookie('masked_gorilla_csrf')
        .should('not.exist')
        .then(() => {
          // null token will omit the 'X-CSRF-HEADER' from request
          sendRequest();
        });
    } else {
      // Send request with masked token
      cy
        .getCookie('masked_gorilla_csrf')
        .should('exist')
        .then(cookie => {
          sendRequest(cookie.value);
        });
    }
  },
);

Cypress.Commands.add('logout', () => {
  cy.clearCookies();
  cy.patientVisit('/');
  cy.getCookie('masked_gorilla_csrf').then(cookie => {
    cy.request({
      url: '/auth/logout',
      method: 'POST',
      headers: { 'x-csrf-token': cookie.value },
    });

    // In case of login redirect we once more go to the homepage
    cy.patientVisit('/');
  });
});

Cypress.Commands.add('nextPage', () => {
  cy
    .get('button.next')
    .should('be.enabled')
    .click();
});

Cypress.Commands.add(
  'resetDb',
  () => {},
  /*
   * Resetting the DB in this manner is slow and should be avoided.
   * Instead of adding this to a test please create a new data record for your test in pkg/testdatagen/scenario/e2ebasic.go
   * For development you can issue `make db_e2e_reset` if you need to clean up your data.
   *
   * cy
   *   .exec('make db_e2e_reset')
   *   .its('code')
   *   .should('eq', 0),
   */
);

//from https://github.com/cypress-io/cypress/issues/669
//Cypress doesn't give the right File constructor, so we grab the window's File
Cypress.Commands.add('upload_file', (selector, fileUrl) => {
  const nameSegments = fileUrl.split('/');
  const name = nameSegments[nameSegments.length - 1];
  const rawType = mime.lookup(name);
  // mime returns false if lookup fails
  const type = rawType ? rawType : '';
  return cy.window().then(win => {
    return cy
      .fixture(fileUrl, 'base64')
      .then(Cypress.Blob.base64StringToBlob)
      .then(blob => {
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
  cy
    .get(`${classSelector} input[type="text"]`)
    .first()
    .type(`{selectall}{backspace}${inputData}`, { force: true, delay: 75 });

  // Click on the first presented option
  cy
    .get(classSelector)
    .find('div[role="option"]')
    .first()
    .click();
}

Cypress.Commands.add('typeInInput', ({ name, value }) => {
  cy
    .get(`input[name="${name}"]`)
    .clear()
    .type(value)
    .blur();
});

// function typeInTextArea({ name, value }) {
Cypress.Commands.add('typeInTextarea', ({ name, value }) => {
  cy
    .get(`textarea[name="${name}"]`)
    .clear()
    .type(value)
    .blur();
});

Cypress.Commands.add('selectDutyStation', (stationName, fieldName) => {
  let classSelector = '.duty-input-box';
  genericSelect(stationName, fieldName, classSelector);
});

Cypress.Commands.add('selectTariff400ngItem', itemName => {
  let classSelector = '.tariff400-select';
  let fieldName = 'tariff400ng_item';
  genericSelect(itemName, fieldName, classSelector);
});

Cypress.Commands.add('setupBaseUrl', appname => {
  // setup baseurl
  switch (appname) {
    case milmoveAppName:
      Cypress.config('baseUrl', 'http://milmovelocal:4000');
      break;
    case officeAppName:
      Cypress.config('baseUrl', 'http://officelocal:4000');
      break;
    case tspAppName:
      Cypress.config('baseUrl', 'http://tsplocal:4000');
      break;
    default:
      break;
  }
});

Cypress.Commands.add('removeFetch', () => {
  // cypress server/route/wait currently does not support window.fetch api
  // https://github.com/cypress-io/cypress/issues/95#issuecomment-347607198
  // delete window.fetch to force fallback to supported xhr.
  // https://github.com/cypress-io/cypress-example-recipes/tree/master/examples/stubbing-spying__window-fetch
  cy.on('window:before:load', win => {
    delete win.fetch;
  });
});
