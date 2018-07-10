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
  cy.visit('/');
  cy.request('POST', 'devlocal-auth/new').then(() => cy.visit('/'));
  //  cy.contains('Local Sign In').click();
  //  cy.contains('Login as New User').click();
});

Cypress.Commands.add('signIntoOffice', () => {
  Cypress.config('baseUrl', 'http://officelocal:4000');
  cy.signInAsUser('9bfa91d2-7a0c-4de0-ae02-b8cf8b4b858b');
});
Cypress.Commands.add('signInAsUser', userId => {
  cy.visit('/');
  cy
    .request({
      method: 'POST',
      url: '/devlocal-auth/login',
      form: true,
      body: { id: userId },
    })
    .then(() => cy.visit('/'));
});

Cypress.Commands.add('nextPage', () => {
  cy
    .get('button.next')
    .should('be.enabled')
    .click();
});

Cypress.Commands.add('resetDb', () =>
  cy
    .exec('make db_e2e_reset')
    .its('code')
    .should('eq', 0),
);
//from https://github.com/cypress-io/cypress/issues/669
//not quite working yet
Cypress.Commands.add('upload_file', (selector, fileUrl, type = '') => {
  return cy
    .fixture(fileUrl, 'base64')
    .then(Cypress.Blob.base64StringToBlob)
    .then(blob => {
      const nameSegments = fileUrl.split('/');
      const name = nameSegments[nameSegments.length - 1];
      const testFile = new File([blob], name, { type });
      const event = { dataTransfer: { files: [testFile] } };
      return cy.get(selector).trigger('drop', event);
    });
});
