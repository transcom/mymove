import { perBuild } from '@jackfranklin/test-data-bot';

import testSubfactory, { TEST_SUBFACTORY_FIELDS } from './subfactory';

import { baseFactory, BASE_FIELDS } from 'utils/test/factories/base';

export const TEST_FACTORY_FIELDS = {
  DEFAULT: 'default',
  SUBFACTORY: 'subfactory',
  POST_BUILD_TOUCHED_FIELD: 'postBuildTouchedField',
};

export const TEST_FACTORY_TRAITS = {
  TEST_TRAIT: 'testTrait',
};

const testFactory = (params) => {
  return baseFactory({
    fields: {
      [TEST_FACTORY_FIELDS.DEFAULT]: 'default',
      [TEST_FACTORY_FIELDS.SUBFACTORY]: (subparams) => testSubfactory(subparams),
      // this also works of course:
      // [TEST_FACTORY_FIELDS.SUBFACTORY]: testSubfactory,
      [TEST_FACTORY_FIELDS.POST_BUILD_TOUCHED_FIELD]: 'default',
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
