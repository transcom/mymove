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

Cypress.Commands.add('nextPage', () => {
  cy
    .get('button.next')
    .should('be.enabled')
    .click();
});
