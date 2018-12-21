import * as mime from 'mime-types';

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

Cypress.Commands.add('signInAsNewUser', (navigateHome = true) => {
  cy.request('POST', 'devlocal-auth/new').then(() => {
    if (navigateHome) {
      cy.visit('/');
    }
  });
  //  cy.contains('Local Sign In').click();
  //  cy.contains('Login as New User').click();
});

Cypress.Commands.add('signIntoMyMoveAsUser', (userId, navigateHome = true) => {
  Cypress.config('baseUrl', 'http://milmovelocal:4000');
  cy.signInAsUser(userId, navigateHome);
});
Cypress.Commands.add('signIntoOffice', (navigateHome = true) => {
  Cypress.config('baseUrl', 'http://officelocal:4000');
  cy.signInAsUser('9bfa91d2-7a0c-4de0-ae02-b8cf8b4b858b', navigateHome);
});
Cypress.Commands.add('signIntoTSP', (navigateHome = true) => {
  Cypress.config('baseUrl', 'http://tsplocal:4000');
  cy.signInAsUser('6cd03e5b-bee8-4e97-a340-fecb8f3d5465', navigateHome);
});
Cypress.Commands.add('signInAsUser', (userId, navigateHome = true) => {
  cy
    .request({
      method: 'POST',
      url: '/devlocal-auth/login',
      form: true,
      body: { id: userId },
    })
    .then(() => {
      if (navigateHome) {
        cy.visit('/');
      }
    });
});

// Reloads the page but makes an attempt to wait for the loading screen to disappear
Cypress.Commands.add('patientReload', () => {
  cy.reload();
  cy.get('h2[data-name="loading-placeholder"]').should('not.exist', { timeout: 10000 });
});

Cypress.Commands.add('logout', () => {
  cy
    .request('/auth/logout')
    .its('status')
    .should('equal', 200);
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
