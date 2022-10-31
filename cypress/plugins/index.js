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
const mkdirp = require('mkdirp');

const storeData = async (data, filepath) => {
  try {
    await mkdirp(path.dirname(filepath));
    fs.writeFile(filepath, JSON.stringify(data));
  } catch (err) {
    console.error(err);
  }
};

let reports = {};

const accumulateReports = (runner, report) => {
  reports[runner] = report;
};

module.exports = (on, config) => {
  // `on` is used to hook into various events Cypress emits
  // `config` is the resolved Cypress config

  on('before:browser:launch', (browser = {}, launchOptions) => {
    prepareAudit(launchOptions);
  });

  // this would default to true
  const doNotThrowA11yErrors = false;

  const runLighthouse = (report) => {
    if (doNotThrowA11yErrors) {
      return lighthouse((report) => {
        const filepath = path.resolve('cypress', `reports/lighthouse_report-${new Date()}.json`);
        storeData(report, filepath);
      });
    }
    return lighthouse();
  };

  const runPa11y = (report) => {
    if (doNotThrowA11yErrors) {
      return pa11y((report) => {
        console.log(report);
        const filepath = path.resolve('cypress', `reports/pa11y_report-${new Date()}.json`);
        storeData(report, filepath);
      });
    }
    return pa11y();
  };

  on('task', {
    lighthouse: runLighthouse(),
    pa11y: runPa11y(),
    a11yAudit: () => {
      lighthouse();
      pa11y();
    },
  });
};
