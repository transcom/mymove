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
  FIELDS: 'fields',
  OVERRIDES: 'overrides',
  LAZY_OVERRIDES: 'lazyOverrides',
  POST_BUILD: 'postBuild',
  TRAITS: 'traits',
};

const applyOverrides = (object, overrides) => {
  const appliedFields = object;
  Object.entries(object).forEach(([field, value]) => {
    let appliedValue;
    switch (typeof value) {
      case 'function':
        if (overrides?.[field]) {
          appliedValue = value({ [BASE_FIELDS.OVERRIDES]: overrides[field] });
        } else {
          appliedValue = value();
        }
        break;
      case 'object':
        if (overrides?.[field]) {
          appliedValue = applyOverrides(value, overrides[field]);
        } else {
          appliedValue = value;
        }
        break;
      default:
        if (overrides && overrides[field]) {
          appliedValue = overrides[field];
        } else {
          appliedValue = value;
        }
        break;
    }
    appliedFields[field] = appliedValue;
  });
  return appliedFields;
};

const basePostBuild = (lazyOverrides, func = (o) => o) => {
  return (object) => {
    func(object);
    if (lazyOverrides) {
      applyOverrides(object, lazyOverrides);
    }
    return object;
  };
};

const camelCaseFields = (fields) => {
  const formattedFields = {};
  Object.entries(fields).forEach(([field, value]) => {
    const formattedField = _.camelCase(field);
    formattedFields[formattedField] = value;
  });
  return formattedFields;
};

const baseFactory = (params) => {
  const { fields, postBuild, lazyOverrides, overrides, traits, useTraits } = params;

  const appliedFields = applyOverrides(fields, overrides);
  const formattedFields = camelCaseFields(appliedFields);

  const builder = build({
    [BASE_FIELDS.FIELDS]: formattedFields,
    [BASE_FIELDS.POST_BUILD]: basePostBuild(lazyOverrides, postBuild),
    [BASE_FIELDS.TRAITS]: traits,
  });

  return builder({ [BASE_FIELDS.TRAITS]: useTraits });
};

export { BASE_FIELDS, baseFactory, basePostBuild, fake, getInternalSpec, getGHCSpec };
