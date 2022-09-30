// Please do not delete. This file and the others in this directory are meant to preserve the integrity of our factory logic.
import { baseFactory } from 'utils/test/factories/base';

export const TEST_SUBFACTORY_FIELDS = {
  DEFAULT: 'default',
  FIELD_TO_OVERRIDE: 'fieldToOverride',
  POST_BUILD_TOUCHED_FIELD: 'postBuildTouchedField',
};

const testSubfactory = (params) => {
  return baseFactory({
    fields: {
      [TEST_SUBFACTORY_FIELDS.DEFAULT]: 'default',
      [TEST_SUBFACTORY_FIELDS.FIELD_TO_OVERRIDE]: 'default',
    },
    postBuild: (object) => {
      object[TEST_SUBFACTORY_FIELDS.POST_BUILD_TOUCHED_FIELD] = 'overridden';
    },
    ...params,
  });
};

export default testSubfactory;
