const { defineConfig } = require('cypress');
// NOTE: THIS FILE IS A WORK IN PROGRESS
module.exports = defineConfig({
  // setupNodeEvents can be defined in either
  // the e2e or component configuration
  e2e: {
    setupNodeEvents(on, config) {
      // bind to the event we care about
      on('<event>', (arg1, arg2) => {
        // plugin stuff here
      });
    },
    baseUrl: 'http://milmovelocal:4000',
  },
});
