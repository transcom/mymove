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

Cypress.Commands.add('signInAsNewUser', () => {
  cy.request('POST', 'devlocal-auth/new').then(() => cy.visit('/'));
  //  cy.contains('Local Sign In').click();
  //  cy.contains('Login as New User').click();
});

Cypress.Commands.add('signIntoOffice', () => {
  Cypress.config('baseUrl', 'http://officelocal:4000');
  cy.signInAsUser('9bfa91d2-7a0c-4de0-ae02-b8cf8b4b858b');
});
Cypress.Commands.add('signIntoTSP', () => {
  Cypress.config('baseUrl', 'http://tsplocal:4000');
  cy.signInAsUser('6cd03e5b-bee8-4e97-a340-fecb8f3d5465');
});
Cypress.Commands.add('signInAsUser', userId => {
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

Cypress.Commands.add('testPremoveSurvey', () => {
  // Click on edit Premove Survey
  cy
    .get('.editable-panel-header')
    .contains('Premove')
    .siblings()
    .click();

  // Enter details in form and save orders
  cy
    .get('input[name="survey.pm_survey_pack_date"]')
    .first()
    .type('8/1/2018')
    .blur();
  cy
    .get('input[name="survey.pm_survey_pickup_date"]')
    .first()
    .type('8/2/2018')
    .blur();
  cy
    .get('input[name="survey.pm_survey_latest_pickup_date"]')
    .first()
    .type('8/3/2018')
    .blur();
  cy
    .get('input[name="survey.pm_survey_earliest_delivery_date"]')
    .first()
    .type('8/4/2018')
    .blur();
  cy
    .get('input[name="survey.pm_survey_latest_delivery_date"]')
    .first()
    .type('8/5/2018')
    .blur();
  cy
    .get('input[name="survey.pm_survey_weight_estimate"]')
    .first()
    .type('6000')
    .blur();
  cy
    .get('input[name="survey.pm_survey_progear_weight_estimate"]')
    .first()
    .type('7000')
    .blur();
  cy
    .get('input[name="survey.pm_survey_spouse_progear_weight_estimate"]')
    .first()
    .type('8000')
    .blur();
  cy
    .get('input[name="survey.pm_survey_notes"]')
    .first()
    .type('Notes notes notes')
    .blur();

  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save')
    .click();

  // Verify data has been saved in the UI
  cy.get('span').contains('7,000 lbs');

  // Refresh browser and make sure changes persist
  cy.reload();

  cy
    .get('div.pm_survey_earliest_delivery_date')
    .get('span')
    .contains('7,000 lbs');
  cy.get('div.pm_survey_notes').contains('Notes notes notes');
});
