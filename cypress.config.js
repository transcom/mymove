/* eslint-disable import/no-extraneous-dependencies */
/* eslint-disable global-require */
/* eslint-disable import/extensions */
const { defineConfig } = require('cypress');

module.exports = defineConfig({
  nodeVersion: 'system',
  lighthouse: {
    performance: 85,
    accessibility: 50,
    'best-practices': 85,
    seo: 85,
    pwa: 50,
  },
  e2e: {
    // We've imported your old cypress plugins here.
    // You may want to clean this up later by importing these.
    setupNodeEvents(on, config) {
      return require('./cypress/plugins/index.js')(on, config);
    },
  },
});
