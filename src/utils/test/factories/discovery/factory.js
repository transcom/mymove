// Please do not delete. This file and the others in this directory are meant to preserve the integrity of our factory logic.
import { perBuild } from '@jackfranklin/test-data-bot';

import testSubfactory, { TEST_SUBFACTORY_FIELDS } from './subfactory';

import { baseFactory, BASE_FIELDS } from 'utils/test/factories/base';

export const TEST_FACTORY_FIELDS = {
  DEFAULT: 'default',
  SUBFACTORY: 'subfactory',
  POST_BUILD_TOUCHED_FIELD: 'postBuildTouchedField',
  SNAKE_CASE_FIELD: 'case_field',
  CAMEL_CASE_FIELD: 'caseField',
};

export const TEST_FACTORY_TRAITS = {
  TEST_TRAIT: 'testTrait',
};

const testFactory = (params) => {
  return baseFactory({
    fields: {
      [TEST_FACTORY_FIELDS.DEFAULT]: 'default',
      [TEST_FACTORY_FIELDS.SUBFACTORY]: (subparams) => testSubfactory(subparams),
      [TEST_FACTORY_FIELDS.POST_BUILD_TOUCHED_FIELD]: 'default',
      [TEST_FACTORY_FIELDS.SNAKE_CASE_FIELD]: 'caseFieldValue',
    },
    postBuild: (object) => {
      object[TEST_SUBFACTORY_FIELDS.POST_BUILD_TOUCHED_FIELD] = 'overridden';
    },
    [BASE_FIELDS.TRAITS]: {
      [TEST_FACTORY_TRAITS.TEST_TRAIT]: {
        [BASE_FIELDS.OVERRIDES]: {
          [TEST_FACTORY_FIELDS.DEFAULT]: perBuild(() => 'overriddenByTrait'),
        },
      },
    },
    ...params,
  });
};

export default testFactory;
