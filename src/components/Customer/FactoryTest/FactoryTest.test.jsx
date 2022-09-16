import { BASE_FIELDS } from 'utils/test/factories/base';
import testFactory, { TEST_FACTORY_FIELDS } from 'utils/test/factories/discovery/factory';
import { TEST_SUBFACTORY_FIELDS } from 'utils/test/factories/discovery/subfactory';

const test = testFactory();

const testWithOverriddenFieldViaOverride = testFactory({
  [BASE_FIELDS.OVERRIDES]: {
    [TEST_FACTORY_FIELDS.DEFAULT]: 'overridden',
  },
});

const testWithOverridenFieldByLazyOverrides = testFactory({
  [BASE_FIELDS.LAZY_OVERRIDES]: [
    {
      [BASE_FIELDS.LAZY_OVERRIDES_FIELD_PATH]: [TEST_FACTORY_FIELDS.POST_BUILD_TOUCHED_FIELD],
      [BASE_FIELDS.LAZY_OVERRIDES_VALUE]: 'overriddenByLazyOverrides',
    },
    {
      [BASE_FIELDS.LAZY_OVERRIDES_FIELD_PATH]: [
        TEST_FACTORY_FIELDS.SUBFACTORY,
        TEST_SUBFACTORY_FIELDS.POST_BUILD_TOUCHED_FIELD,
      ],
      [BASE_FIELDS.LAZY_OVERRIDES_VALUE]: 'overriddenByLazyOverrides',
    },
  ],
});

const testWithTrait = testFactory({ useTraits: ['testTrait'] });

describe('testFactory', () => {
  it('has correct default values', () => {
    expect(test[TEST_FACTORY_FIELDS.DEFAULT]).toBe('default');
    expect(test[TEST_FACTORY_FIELDS.SUBFACTORY][TEST_SUBFACTORY_FIELDS.DEFAULT]).toBe('default');
    expect(test[TEST_FACTORY_FIELDS.SUBFACTORY][TEST_SUBFACTORY_FIELDS.FIELD_TO_OVERRIDE]).toBe('default');
  });
  it('has an overridden field', () => {
    expect(testWithOverriddenFieldViaOverride[TEST_FACTORY_FIELDS.DEFAULT]).toBe('overridden');
  });
  it('has an overridden fields by postBuild', () => {
    expect(test[TEST_FACTORY_FIELDS.POST_BUILD_TOUCHED_FIELD]).toBe('overridden');
    expect(test[TEST_FACTORY_FIELDS.SUBFACTORY][TEST_SUBFACTORY_FIELDS.POST_BUILD_TOUCHED_FIELD]).toBe('overridden');
  });
  it('has overridden fields by its lazy overrides', () => {
    expect(testWithOverridenFieldByLazyOverrides[TEST_FACTORY_FIELDS.POST_BUILD_TOUCHED_FIELD]).toBe(
      'overriddenByLazyOverrides',
    );
    expect(
      testWithOverridenFieldByLazyOverrides[TEST_FACTORY_FIELDS.SUBFACTORY][
        TEST_SUBFACTORY_FIELDS.POST_BUILD_TOUCHED_FIELD
      ],
    ).toBe('overriddenByLazyOverrides');
  });
  it('has an overridden field by trait', () => {
    expect(testWithTrait[TEST_FACTORY_FIELDS.DEFAULT]).toBe('overriddenByTrait');
  });
});
