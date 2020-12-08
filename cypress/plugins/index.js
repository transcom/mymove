// ***********************************************************
// This example plugins/index.js can be used to load plugins
//
// You can change the location of this file or turn off loading
// the plugins file with the 'pluginsFile' configuration option.
//
// You can read more here:
// https://on.cypress.io/plugins-guide
// ***********************************************************

// This function is called when a project is opened or re-opened (e.g. due to
// the project's config changing)

const { lighthouse, pa11y, prepareAudit } = require('cypress-audit');
const fs = require('fs');
const path = require('path');

const storeData = (data, filepath) => {
  try {
    fs.writeFile(filepath, JSON.stringify(data));
  } catch (err) {
    console.error(err);
  }
};

module.exports = (on, config) => {
  // `on` is used to hook into various events Cypress emits
  // `config` is the resolved Cypress config

  on('before:browser:launch', (browser = {}, launchOptions) => {
    prepareAudit(launchOptions);
  });

  on('task', {
    lighthouse: lighthouse((report) => {
      const filepath = path.resolve('cypress', 'results/lighthouse_report.json');
      storeData(report, filepath);
    }),
    pa11y: pa11y((report) => {
      const filepath = path.resolve('cypress', 'results/pa11y_report.json');
      storeData(report, filepath);
    }),
  });
};
