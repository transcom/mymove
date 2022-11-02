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
const moment = require('moment');
const a11yReportPath = `reports/a11y_report-${moment().format('MM-DD-YYYY')}.json`;
let currentIssues;
try {
  currentIssues = require(a11yReportPath);
} catch (e) {
  currentIssues = [];
}

const storeData = async (data, filepath) => {
  try {
    await mkdirp(path.dirname(filepath));
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

  // this would default to true
  const doNotThrowA11yErrors = true;

  const accumulateReports = (report) => {
    report.issues.forEach((issue) => {
      if (
        !currentIssues.find(
          (existingIssue) =>
            existingIssue.selector === issue.selector &&
            existingIssue.context === issue.context &&
            existingIssue.code === issue.code,
        )
      ) {
        console.log('new issue', issue);
        console.log('new report', report);
        currentIssues.push(issue);
      }
    });
  };

  on('task', {
    lighthouse: lighthouse((report) => {
      accumulateReports(report);
      const filepath = path.resolve('cypress', a11yReportPath);
      storeData(currentIssues, filepath);
    }),
    pa11y: pa11y((report) => {
      accumulateReports(report);
      const filepath = path.resolve('cypress', a11yReportPath);
      storeData(currentIssues, filepath);
    }),
  });
};
