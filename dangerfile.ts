// eslint-disable-next-line import/no-extraneous-dependencies
import { danger, warn, fail } from 'danger';
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

function diffContainsNosec(diffForFile) {
  return !!diffForFile.includes('#nosec');
}

function diffContainsEslint(diffForFile) {
  return !!diffForFile.includes('eslint-disable');
}

function doesLineHaveProhibitedOverride(disablingString) {
  const okBypassRules = [
    'no-underscore-dangle',
    'prefer-object-spread',
    'object-shorthand',
    'camelcase',
    'react/jsx-props-no-spreading',
    'react/destructuring-assignment',
    'react/forbid-prop-types',
    'react/prefer-stateless-function',
    'react/sort-comp',
    'import/no-extraneous-dependencies',
    'import/order',
    'import/prefer-default-export',
    'import/no-named-as-default',
  ];
  let prohibitedOverrideMsg = '';
  // disablingStringParts format: 'eslint-disable-next-line no-jsx, no-default'
  // split along commas and/or spaces and remove surrounding spaces
  const disablingStringParts = disablingString
    .trim()
    .split(/[\s,]+/)
    .map((item) => item.trim())
    .filter((str) => !str.includes('*/')); // edgecase where string has a dangling */ or */}
  // disablingStringParts format: ['eslint-disable-next-line', 'no-jsx', 'no-default']

  if (disablingStringParts.length === 1) {
    // fail because rule should be specified
    prohibitedOverrideMsg = 'Must specify the rule you are disabling';
    return prohibitedOverrideMsg;
  }

  // rules format: ['no-jsx', 'no-default']
  const rules = disablingStringParts.slice(1);
  for (let r = 0; r < rules.length; r += 1) {
    const rule = rules[r];
    if (!okBypassRules.includes(rule)) {
      prohibitedOverrideMsg = `Contains a rule that is not in the permitted eslint list. You are free to disable only: (\n${okBypassRules.map(
        (q) => `${q}\n`,
      )})`;
      break;
    }
  }
  return prohibitedOverrideMsg;
}

function checkPRHasProhibitedLinterOverride(dangerJSDiffCollection) {
  let badOverrideMsg = '';
  Object.keys(dangerJSDiffCollection).forEach((d) => {
    const diffFile = dangerJSDiffCollection[`${d}`];
    const diff = diffFile.added;
    if (diffContainsNosec(diff)) {
      badOverrideMsg = 'Contains prohibited linter override "#nosec".';
      return;
    }

    if (!diffContainsEslint(diff)) {
      return;
    }

    // split file diffs into lines
    const lines = diff.split('\n');
    for (let l = 0; l < lines.length; l += 1) {
      const line = lines[l];
      if (diffContainsEslint(line)) {
        // check for comment marker (// or /*)
        // eg line: 'const whatever = something() // eslint-disable-line'
        let lineParts = line.split('//');
        if (lineParts.length === 1) {
          lineParts = line.split('/*');
          if (lineParts.length === 1) {
            throw new Error('uhhhh, how did we find eslint disable but no // or /*');
          }
        }

        // eg lineParts: ['const whatever = something()', 'eslint-disable-line']
        const stringAfterCommentMarker = lineParts[1];
        badOverrideMsg = doesLineHaveProhibitedOverride(stringAfterCommentMarker);
      }
    }
  });

  return badOverrideMsg;
}

const bypassingLinterChecks = async () => {
  const allFiles = danger.git.modified_files.concat(danger.git.created_files).filter(file => file.includes('src/') || file.includes('pkg/'));
  const diffsByFile = await Promise.all(allFiles.map((f) => danger.git.diffForFile(f)));
  const dangerMsgSegment = checkPRHasProhibitedLinterOverride(diffsByFile);
  if (dangerMsgSegment) {
    warn(
      `It looks like you are attempting to bypass a linter rule, which is not within
      security compliance rules.\n** ${dangerMsgSegment} **\n Please remove the bypass code and address the underlying issue. cc: @transcom/Truss-Pamplemoose`,
    );
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
Please update it in both the root package.json and the [cirlcleci-docker/milmove-cypress/](https://github.com/transcom/circleci-docker) folder's separate package.json`,
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
if (!danger.github || (danger.github && danger.github.pr.user.login !== 'dependabot-preview[bot]')) {
  githubChecks();
  fileChecks();
  checkYarnAudit();
  cypressUpdateChecks();
  bypassingLinterChecks();
}
