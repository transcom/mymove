import fs from 'fs';
import path from 'path';
import process from 'process';

import { faker } from '@faker-js/faker';
import { build, perBuild } from '@jackfranklin/test-data-bot';
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
  // RA Validator Status: Mitigated
  // RA Validator: leodis.f.scott.civ@mail.mil
  // RA Modified Severity: CAT III
  /* eslint-disable-next-line security/detect-non-literal-fs-filename */
  return yaml.load(fs.readFileSync(yamlPath, 'utf8'));
};

const getInternalSpec = () => {
  return loadSpec('internal.yaml');
};

const getGHCSpec = () => {
  return loadSpec('ghc.yaml');
};

const BASE_FIELDS = {
  FIELDS: 'fields',
  OVERRIDES: 'overrides',
  LAZY_OVERRIDES: 'lazyOverrides',
  POST_BUILD: 'postBuild',
  TRAITS: 'traits',
};

/**
 * applyOverrides takes an object of fields and an object of overrides and recursively applies the overrides where their structure matches the fields'.
 * e.g. fields.field.subfield will be overridden by fields.field.subfield if the latter exists.
 * If an override does not exist, and the field's value is a function, it's called.
 * Otherwise, the field's value is unaffected.
 * The return value is the fields object with the overrides applied.
 * @param {*} fields - An object containing fields and their values
 * @param {*} overrides - An object containing values to override for their corresponding fields
 */
const applyOverrides = (fields, overrides) => {
  const appliedFields = fields;
  Object.entries(fields).forEach(([field, value]) => {
    switch (typeof value) {
      case 'function':
        if (overrides && overrides[field]) {
          appliedFields[field] = value({ [BASE_FIELDS.OVERRIDES]: overrides[field] });
        } else {
          appliedFields[field] = value();
        }
        break;
      case 'object':
        if (overrides && overrides[field]) {
          if (typeof value.call === 'function' && value.generatorType === 'perBuild') {
            // we're in a perBuild function; just replace the value
            if (typeof overrides[field] === 'function') {
              appliedFields[field] = overrides[field]();
            } else {
              appliedFields[field] = overrides[field];
            }
          } else {
            // apply this function's logic to nested values
            appliedFields[field] = applyOverrides(value, overrides[field]);
          }
        }
        break;
      default:
        if (overrides && overrides[field]) {
          appliedFields[field] = overrides[field];
        }
        break;
    }
  });
  return appliedFields;
};

/**
 * basePostBuild takes an object of lazyOverrides and an optional function, and returns a function that:
 * - takes an object of fields
 * - applies the function to that object
 * - calls applyOverrides on that object, with lazyOverrides
 * - returns the object
 * This function is meant to serve as a wrapper for test-data-bot's postBuild(), with the additional ability to apply lazy overrides.
 * @param {*} lazyOverrides - An object containing overrides and their values
 * @param {*} postBuild - A function to apply to the fields before overrides are applied. No-op by default.
 */
const basePostBuild = (lazyOverrides, postBuild = (o) => o) => {
  return (fields) => {
    postBuild(fields);
    if (lazyOverrides) {
      applyOverrides(fields, lazyOverrides);
    }
    return fields;
  };
};

/**
 * camelCaseFields takes an object of fields and replaces its snake_cased or kebab-cased keys with camelCased keys.
 * @param {*} fields - An object of fields
 */
const camelCaseFields = (fields) => {
  const formattedFields = {};
  Object.entries(fields).forEach(([field, value]) => {
    const formattedField = field.replace(/([-_][a-z])/gi, ($1) => {
      return $1.toUpperCase().replace('-', '').replace('_', '');
    });
    formattedFields[formattedField] = value;
  });
  return formattedFields;
};

/**
 * baseFactory takes a set of params, applies overrides to fields, camelCases the fields, constructs a factory capable of applying lazy overrides, and then calls that factory, applying any traits.
 * This function is meant to serve as a wrapper for test-data-bot's build(), with the additional ability to camelCase fields and apply lazy overrides.
 * @param {*} params - An object optionally containing fields, postBuild, lazyOverrides, overrides, traits, and useTraits keys
 */
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
