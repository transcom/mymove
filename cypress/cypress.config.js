const { defineConfig } = require('cypress');
const mochaConfig = require('./mocha-reporter-config.json');
const setupNodeEvents = require('./plugins');
// NOTE: THIS FILE IS A WORK IN PROGRESS
module.exports = defineConfig({
  // setupNodeEvents can be defined in either
  // the e2e or component configuration
  e2e: {
    setupNodeEvents,
    baseUrl: 'http://milmovelocal:4000',
  },
  lighthouse: {
    thresholds: {
      performance: 85,
      accessibility: 50,
      'best-practices': 85,
      seo: 85,
      pwa: 50,
    },
  },
  reporter: 'cypress-multi-reporters',
  reporterOptions: {
    configFile: 'mocha-reporter-config.json',
  },
  viewportWidth: 1440,
  viewportHeight: 900,
  videoUploadOnPasses: false,
});
