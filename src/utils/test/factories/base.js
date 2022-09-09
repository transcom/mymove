import fs from 'fs';
import path from 'path';
import process from 'process';

import { faker } from '@faker-js/faker';
import { perBuild } from '@jackfranklin/test-data-bot';
import yaml from 'js-yaml';

faker.setLocale('en_US');

const fake = (callback) => {
  return perBuild(() => callback(faker));
};

const loadSpec = (fileName) => {
  const yamlPath = path.join(process.cwd(), `./swagger/${fileName}`);

  // RA Summary: eslint - security/detect-non-literal-fs-filename - OWASP Path Traversal
  // RA: Detects variable in filename argument of fs calls, which might allow an attacker to access anything on your
  // RA: system.
  // RA: In this case, this doesn't contain user input and serves to be able to grab the local swagger files.
  // RA Developer Status: Mitigated
  // RA Validator Status: {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}
  // RA Modified Severity:
  /* eslint-disable-next-line security/detect-non-literal-fs-filename */
  return yaml.load(fs.readFileSync(yamlPath, 'utf8'));
};

let internalSpec = null;

const getInternalSpec = () => {
  if (!internalSpec) {
    internalSpec = loadSpec('internal.yaml');
  }

  return internalSpec;
};

let ghcSpec = null;

const getGHCSpec = () => {
  if (!ghcSpec) {
    ghcSpec = loadSpec('ghc.yaml');
  }

  return ghcSpec;
};

export { fake, getInternalSpec, getGHCSpec };
