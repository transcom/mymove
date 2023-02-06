// eslint-disable-next-line import/no-extraneous-dependencies
import { danger, warn, fail } from 'danger';
//RA Summary: eslint-plugin-security - detect-child-process -
//RA Executing commands from an untrusted source or in an untrusted environment can cause an application to execute malicious commands on behalf of an attacker.
//RA: Locates usages of child process.
//RA: This usage checks for any critical or high vulnerabilities and upgrades in our dependencies to alert the Github user.
//RA: This usage does not utilize any user input and there is no opening for command injection.
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Validator: leodis.f.scott.civ@mail.mil
//RA Modified Severity: CAT III
// eslint-disable-next-line security/detect-child-process
const child = require('child_process');
// eslint-disable-next-line import/no-extraneous-dependencies
const jiraIssue = require('danger-plugin-jira-issue').default;

const githubChecks = () => {
  if (danger.github) {
    // No PR is too small to include a description of why you made a change
    if (danger.github.pr.body.length < 10) {
      warn('Please include a description of your PR changes.');
    }
    // PRs should have a Jira ID in the title
    jiraIssue({
      key: 'MB',
      url: 'https://dp3.atlassian.net/browse',
      location: 'title',
    });
  }
};

const fileChecks = () => {
  // load all modified and new files
  const allFiles = danger.git.modified_files.concat(danger.git.created_files);

  const legacyFiles = danger.git.fileMatch('src/shared/**/*', 'src/scenes/**/*');

  if (legacyFiles.created) {
    fail(`New files have been created under one of the legacy directories
(src/shared or src/scenes). Please relocate them according to the file structure described [here](https://transcom.github.io/mymove-docs/docs/dev/contributing/frontend/frontend#file-layout--naming).

View the [frontend file org ADR](https://github.com/transcom/mymove/blob/main/docs/adr/0048-frontend-file-org.md) for more information`);
  }

  if (legacyFiles.modified) {
    warn(`Files located in legacy directories (src/shared or src/scenes) have
been edited. Are you sure you don’t want to also relocate them to the new [file structure](https://transcom.github.io/mymove-docs/docs/dev/contributing/frontend/frontend#file-layout--naming)?

View the [frontend file org ADR](https://github.com/transcom/mymove/blob/main/docs/adr/0048-frontend-file-org.md) for more information`);
  }

  // Request changes to app code to also include changes to tests.
  const hasAppChanges = allFiles.some((path) => !!path.match(/src\/.*\.jsx?/));
  const hasTestChanges = allFiles.some((path) => !!path.match(/src\/.*\.test\.jsx?/));
  if (hasAppChanges && !hasTestChanges) {
    warn('This PR does not include changes to unit tests, even though it affects app code.');
  }

  // Require new src/components files to include changes to storybook
  const hasComponentChanges = danger.git.created_files.some((path) => path.includes('src/components'));
  const hasStorybookChanges = allFiles.some(
    (path) => path.includes('src/stories') || !!path.match(/src\/.*\.stories.jsx?/),
  );

  if (hasComponentChanges && !hasStorybookChanges) {
    warn('This PR does not include changes to storybook, even though it affects component code.');
  }

  // Request update of yarn.lock if package.json changed but yarn.lock isn't
  const packageChanged = allFiles.includes('package.json');
  const lockfileChanged = allFiles.includes('yarn.lock');
  if (packageChanged && !lockfileChanged) {
    const message = 'Changes were made to package.json, but not to yarn.lock';
    const idea = 'Perhaps you need to run `yarn install`?';
    warn(`${message} - <i>${idea}</i>`);
  }
};

const cypressUpdateChecks = async () => {
  // load all modified and new files
  const allFiles = danger.git.modified_files.concat(danger.git.created_files);

  // check if relevant package.jsons have changed
  const rootPackageFile = 'package.json';
  const rootPackageChanged = allFiles.includes(rootPackageFile);
  const cypressPackageNames = [
    '"cypress":',
    '"cypress-audit":',
    '"cypress-multi-reporters":',
    '"cypress-wait-until":',
    '"mocha":',
    '"mocha-junit-reporter":',
    '"moment":',
  ];

  let hasRootCypressDepChanged = false;

  // if root changed, check for cypress in diff
  if (rootPackageChanged) {
    const rootPackageDiff = await danger.git.diffForFile(rootPackageFile);
    cypressPackageNames.forEach((cypressPackageName) => {
      if (hasRootCypressDepChanged || (rootPackageDiff && rootPackageDiff.diff.includes(cypressPackageName))) {
        hasRootCypressDepChanged = true;
      }
    });
  }

  if (hasRootCypressDepChanged) {
    warn(
      `It looks like you updated the Cypress package dependency in one of two required places.
Please update it in both the root package.json and the [circleci-docker/milmove-cypress/](https://github.com/transcom/circleci-docker) folder's separate package.json`,
    );
  }
};

const checkYarnAudit = () => {
  const result = child.spawnSync('yarn', ['audit', '--groups=dependencies', '--level=high', '--json']);
  const output = result.stdout.toString().split('\n');
  const summary = JSON.parse(output[output.length - 2]);
  if (
    'data' in summary &&
    'vulnerabilities' in summary.data &&
    'high' in summary.data.vulnerabilities &&
    'critical' in summary.data.vulnerabilities
  ) {
    if (summary.data.vulnerabilities.high > 0 || summary.data.vulnerabilities.critical > 0) {
      let issuesFound = 'Yarn Audit Issues Found:\n';
      output.forEach((rawAudit) => {
        try {
          const audit = JSON.parse(rawAudit);
          if (audit.type === 'auditAdvisory') {
            issuesFound +=
              `${audit.data.advisory.severity} - ${audit.data.advisory.title}\n` +
              `Package ${audit.data.advisory.module_name}\n` +
              `Patched in ${audit.data.advisory.patched_versions}\n` +
              `Dependency of ${audit.data.resolution.path.split('>')[0]}\n` +
              `Path ${audit.data.resolution.path.replace(/>/g, ' > ')}\n` +
              `More info ${audit.data.advisory.url}\n\n`;
          }
        } catch {
          // not all outputs maybe json and that's okay
        }
      });
      fail(
        `${issuesFound}${summary.data.vulnerabilities.high} high vulnerabilities and ` +
          `${summary.data.vulnerabilities.critical} critical vulnerabilities found`,
      );
    }
  } else {
    warn(`Couldn't find summary of vulnerabilities from yarn audit`);
  }
};

// skip these checks if PR is by dependabot, if we don't have a github object let it run also since we are local
if (!danger.github || (danger.github && danger.github.pr.user.login !== 'dependabot[bot]')) {
  githubChecks();
  fileChecks();
  checkYarnAudit();
  cypressUpdateChecks();
}
