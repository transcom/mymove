import * as child from 'child_process';

/* eslint-disable import/no-extraneous-dependencies */
import { includes, replace } from 'lodash';
import { danger, warn, fail } from 'danger';
import jiraIssue from 'danger-plugin-jira-issue';

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
(src/shared or src/scenes). Please relocate them according to the file structure described [here](https://github.com/transcom/mymove/wiki/frontend#file-layout--naming).

View the [frontend file org ADR](https://github.com/transcom/mymove/blob/master/docs/adr/0048-frontend-file-org.md) for more information`);
  }

  if (legacyFiles.modified) {
    warn(`Files located in legacy directories (src/shared or src/scenes) have
been edited. Are you sure you don’t want to also relocate them to the new [file structure](https://github.com/transcom/mymove/wiki/frontend#file-layout--naming)?

View the [frontend file org ADR](https://github.com/transcom/mymove/blob/master/docs/adr/0048-frontend-file-org.md) for more information`);
  }

  // Request changes to app code to also include changes to tests.
  const hasAppChanges = allFiles.some((path) => !!path.match(/src\/.*\.jsx?/));
  const hasTestChanges = allFiles.some((path) => !!path.match(/src\/.*\.test\.jsx?/));
  if (hasAppChanges && !hasTestChanges) {
    warn('This PR does not include changes to unit tests, even though it affects app code.');
  }

  // Require new src/components files to include changes to storybook
  const hasComponentChanges = danger.git.created_files.some((path) => includes(path, 'src/components'));
  const hasStorybookChanges = allFiles.some(
    (path) => includes(path, 'src/stories') || !!path.match(/src\/.*\.stories.jsx?/),
  );

  if (hasComponentChanges && !hasStorybookChanges) {
    warn('This PR does not include changes to storybook, even though it affects component code.');
  }

  // Request update of yarn.lock if package.json changed but yarn.lock isn't
  const packageChanged = includes(allFiles, 'package.json');
  const lockfileChanged = includes(allFiles, 'yarn.lock');
  // eslint-disable-next-line no-constant-condition
  if (false && packageChanged && !lockfileChanged) {
    const message = 'Changes were made to package.json, but not to yarn.lock';
    const idea = 'Perhaps you need to run `yarn install`?';
    warn(`${message} - <i>${idea}</i>`);
  }
};
// eslint-disable jsx-no-target-blank test
const bypassingLinterChecks = async () => {
  // load all modified and new files
  const allFiles = danger.git.modified_files.concat(danger.git.created_files);
  const bypassCodes = ['eslint-disable', 'eslint-disable-next-line', '#nosec'];
  const okBypassRules = [
    'no-underscore-dangle',
    'prefer-object-spread',
    'object-shorthand',
    'camelcase',
    'jsx-props-no-spreading',
    'destructuring-assignment',
    'forbid-prop-types',
    'prefer-stateless-function',
    'sort-comp',
    'import/no-extraneous-dependencies',
    'order',
    'prefer-default-export',
    'no-named-as-default',
  ];
  let addedByPassCode = false;
  const diffs = await Promise.all(allFiles.map((f) => danger.git.diffForFile(f)));

  for (let i = 0; i < diffs.length; i += 1) {
    const diff = diffs[Number(i)];
    const diffsWithbypassCodes = bypassCodes.find((b) => diff.diff.includes(b));
    if (diffsWithbypassCodes) {
      for (const rule in okBypassRules) {
        if (diff.diff.includes(okBypassRules[`${rule}`]) === false) {
          addedByPassCode = true;
          break;
        }
      }
    }
  }

  if (addedByPassCode === true) {
    warn(
      `It looks like you are attempting to bypass a linter rule, which is not a sustainable solution to meet
      security compliance rules. Please remove the bypass code and address the underlying issue. cc: @transcom/Truss-Pamplemoose`,
    );
  }
};

const cypressUpdateChecks = async () => {
  // load all modified and new files
  const allFiles = danger.git.modified_files.concat(danger.git.created_files);

  // check if relevant package.jsons have changed
  const rootPackageFile = 'package.json';
  const cypressPackageFile = 'cypress/package.json';
  const rootPackageChanged = includes(allFiles, rootPackageFile);
  const cypressPackageChanged = includes(allFiles, 'cypress/package.json');
  const cypressPackageName = '"cypress":';
  const versionRegex = /(~|\^|)\d+.\d+.\d+/;

  let hasRootCypressDepChanged = false;
  let hasCypressPackageCypressDepChanged = false;
  let rootVersion;
  let cypressPackageVersion;

  // if root changed, check for cypress in diff
  if (rootPackageChanged) {
    const rootPackageDiff = await danger.git.diffForFile(rootPackageFile);
    if (rootPackageDiff && rootPackageDiff.diff.includes(cypressPackageName)) {
      hasRootCypressDepChanged = true;

      const diff = rootPackageDiff.diff.split(cypressPackageName)[1]; // we don't care about diff before cypress
      [rootVersion] = diff.match(versionRegex); // the first version # will be cypress's
    }
  }

  // if cypress package changed, check for cypress in diff
  if (cypressPackageChanged) {
    const cypressPackageDiff = await danger.git.diffForFile(cypressPackageFile);
    if (cypressPackageDiff && cypressPackageDiff.diff.includes(cypressPackageName)) {
      hasCypressPackageCypressDepChanged = true;

      const diff = cypressPackageDiff.diff.split(cypressPackageName)[1]; // we don't care about diff before cypress
      [cypressPackageVersion] = diff.match(versionRegex); // the first version # will be cypress's
    }
  }

  if (hasRootCypressDepChanged !== hasCypressPackageCypressDepChanged) {
    warn(
      `It looks like you updated the Cypress package dependency in one of two
required places. Please update it in both the root package.json and the cypress/
folder's separate package.json`,
    );
  } else if (rootVersion !== cypressPackageVersion) {
    warn(
      `It looks like there is a Cypress version mismatch between the root
package.json and the cypress/ folder's separate package.json. Please double
check they have the same version.`,
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
              `Path ${replace(audit.data.resolution.path, />/g, ' > ')}\n` +
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
if (!danger.github || (danger.github && danger.github.pr.user.login !== 'dependabot-preview[bot]')) {
  githubChecks();
  fileChecks();
  checkYarnAudit();
  cypressUpdateChecks();
  bypassingLinterChecks();
}
