import { includes } from 'lodash';
import { danger, warn, fail } from 'danger';

// No PR is too small to include a description of why you made a change
if (danger.github && danger.github.pr.body.length < 10) {
  warn('Please include a description of your PR changes.');
}

// load all modified and new files
const allFiles = danger.git.modified_files.concat(danger.git.created_files);

// Request changes to app code to also include changes to tests.
const hasAppChanges = allFiles.some((path) => !!path.match(/src\/.*\.jsx?/));
const hasTestChanges = allFiles.some((path) => !!path.match(/src\/.*\.test\.jsx?/));
if (hasAppChanges && !hasTestChanges) {
  warn('This PR does not include changes to unit tests, even though it affects app code.');
}

// Require new src/components files to include changes to storybook
const hasComponentChanges = danger.git.created_files.some((path) => includes(path, 'src/components'));
const hasStorybookChanges = allFiles.some((path) => includes(path, 'src/stories'));

if (hasComponentChanges && !hasStorybookChanges) {
  fail('This PR does not include changes to storybook, even though it affects component code.');
}

// Request update of yarn.lock if package.json changed but yarn.lock isn't
const packageChanged = includes(allFiles, 'package.json');
const lockfileChanged = includes(allFiles, 'yarn.lock');
if (packageChanged && !lockfileChanged) {
  const message = 'Changes were made to package.json, but not to yarn.lock';
  const idea = 'Perhaps you need to run `yarn install`?';
  warn(`${message} - <i>${idea}</i>`);
}
