/* global cy, Cypress */

// ***********************************************************
// This example support/index.js is processed and
// loaded automatically before your test files.
//
// This is a great place to put global configuration and
// behavior that modifies Cypress.
//
// You can change the location of this file or turn off
// automatically serving support files with the
// 'supportFile' configuration option.
//
// You can read more here:
// https://on.cypress.io/configuration
// ***********************************************************

// Import commands.js using ES2015 syntax:
import './commands';

// Capture the full screen by default
// These options are overridden for auto-screenshot on fail
Cypress.Screenshot.defaults({
  capture: 'fullPage',
});

cy.on('before:browser:launch', (browser = {}, args) => {
  // Disable shared memory when running headless since running out of memory can cause cypress to hang indefinitely
  // https://github.com/cypress-io/cypress/issues/8206
  if (config.env.headless && browser.name === 'chrome') {
    args.push('--disable-dev-shm-usage');
  } else if (config.env.headless && browser.name === 'electron') {
    args['disable-dev-shm-usage'] = true;
  }
});

afterEach(function () {
  if (this.currentTest.state === 'failed') {
    // Take another screenshot so we get the full page
    cy.screenshot();
  }
});

// Alternatively you can use CommonJS syntax:
// require('./commands')
