const { defineConfig } = require('cypress');
const setupNodeEvents = require('./cypress/plugins/index.js');

module.exports = defineConfig({
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
