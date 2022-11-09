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
const a11yReportPath = `cypress/reports/a11y_report-${moment().format('MM-DD-YYYY')}.json`;
const lighthouseReportsPath = `cypress/reports/`;

const getCurrentIssues = () => {
  try {
    return require(a11yReportPath);
  } catch (e) {
    return [];
  }
};

let currentTest = null;

let currentIssues = getCurrentIssues();

const ensureReportPath = async (filepath) => {
  if (!fs.existsSync(filepath)) {
    try {
      await mkdirp(path.dirname(filepath));
    } catch (err) {
      console.error(err);
    }
  }
};

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
    ensureReportPath(a11yReportPath);
  });

  const issuesMatch = (issue1, issue2) => {
    return issue1.selector === issue2.selector && issue1.context === issue2.context && issue1.code === issue2.code;
  };

  const accumulateIssues = (report) => {
    report.issues.forEach((issue) => {
      const matchingIndex = currentIssues.findIndex((existingIssue) => issuesMatch(existingIssue, issue));
      if (matchingIndex >= 0) {
        if (!currentIssues[matchingIndex].tests.find((existingIssue) => issuesMatch(existingIssue, issue))) {
          currentIssues[matchingIndex].tests.push(currentTest);
        }
      } else {
        issue.tests = [currentTest];
        currentIssues.push(issue);
      }
    });
    return currentIssues;
  };

  const pa11yReport = pa11y((report) => {
    const currentIssues = accumulateIssues(report);
    storeData(currentIssues, a11yReportPath);
  });

  const lighthouseReport = lighthouse(async (lighthouseReport) => {
    await mkdirp(path.dirname(lighthouseReportsPath));

    fs.writeFile(
      `${lighthouseReportsPath}lighthouse-${moment().format('MM-DD-YYYY')}.json`,
      lighthouseReport.report,
      (error) => {
        error ? console.log(error) : console.log('Report created successfully');
      },
    );
  });

  on('task', {
    lighthouse: lighthouseReport,
    pa11y: pa11yReport,
    log: (message) => {
      console.log(message);
      return null;
    },
    setCurrentTest: (test) => {
      currentTest = test;
      return null;
    },
  });
};
