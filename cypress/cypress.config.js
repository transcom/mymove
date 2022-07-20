const { defineConfig } = require('cypress');
const mochaConfig = require('mocha-reporter-config.json');
const setupNodeEvents = require('./cypress/plugins/index.js');
// NOTE: THIS FILE IS A WORK IN PROGRESS
module.exports = defineConfig({
  // setupNodeEvents can be defined in either
  // the e2e or component configuration
  e2e: {
    setupNodeEvents,
    baseUrl: 'http://milmovelocal:4000',
  },
  component: {
    devServer: {
      framework: 'react', // or vue
      bundler: 'mocha',
      mochaConfig,
    },
  },
});
