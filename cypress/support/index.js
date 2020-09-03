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

// Disable smooth-scrolling on every page load (https://github.com/cypress-io/cypress/issues/3200)
Cypress.on('window:load', (win) => {
  const { document } = win;
  const node = document.createElement('style');
  node.innerHTML = 'html { scroll-behavior: inherit !important; }';
  document.body.appendChild(node);
});
