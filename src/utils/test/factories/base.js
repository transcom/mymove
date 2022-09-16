import fs from 'fs';
import path from 'path';
import process from 'process';

import { faker } from '@faker-js/faker';
import { build, perBuild } from '@jackfranklin/test-data-bot';
import yaml from 'js-yaml';
import { _ } from 'lodash';

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

const BASE_FIELDS = {
  OVERRIDES: 'overrides',
  LAZY_OVERRIDES: 'lazyOverrides',
  LAZY_OVERRIDES_FIELD_PATH: 'fieldPath',
  LAZY_OVERRIDES_VALUE: 'value',
  TRAITS: 'traits',
};

const applyLazyOverrides = (object, lazyOverrides = []) => {
  lazyOverrides.forEach(({ fieldPath, value }) => {
    _.set(object, fieldPath, value);
  });
};

const basePostBuild = (lazyOverrides, func = (o) => o) => {
  return (object) => {
    func(object);
    applyLazyOverrides(object, lazyOverrides);
    return object;
  };
};

const baseFactory = (params) => {
  const { fields, postBuild, lazyOverrides, overrides, traits, useTraits } = params;
  // these will be initial overrides. you'll somehow need to iterate over the fields and pass these to subfactories' fields if there are any
  // const allOverrides = {
  //   ...overrides,
  // };

  const builder = build({
    fields,
    postBuild: basePostBuild(lazyOverrides, postBuild),
    traits,
  });

  return builder({ overrides, traits: useTraits });
};

export { BASE_FIELDS, baseFactory, basePostBuild, fake, getInternalSpec, getGHCSpec, applyLazyOverrides };
