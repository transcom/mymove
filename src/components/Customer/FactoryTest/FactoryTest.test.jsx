import { BASE_FIELDS } from 'utils/test/factories/base';
import testFactory, { TEST_FACTORY_FIELDS } from 'utils/test/factories/discovery/factory';
import { TEST_SUBFACTORY_FIELDS } from 'utils/test/factories/discovery/subfactory';
import dutyLocationFactory, { DUTY_LOCATION_FIELDS } from 'utils/test/factories/dutyLocation';
import { ADDRESS_FIELDS } from 'utils/test/factories/address';

const test = testFactory();

const testWithOverriddenFieldViaOverride = testFactory({
  [BASE_FIELDS.OVERRIDES]: {
    [TEST_FACTORY_FIELDS.DEFAULT]: 'overriddenByOverrides',
    [TEST_FACTORY_FIELDS.SUBFACTORY]: {
      [TEST_SUBFACTORY_FIELDS.FIELD_TO_OVERRIDE]: 'overriddenBySuboverrides',
    },
  },
});

const testWithOverridenFieldByLazyOverrides = testFactory({
  [BASE_FIELDS.LAZY_OVERRIDES]: {
    [TEST_FACTORY_FIELDS.POST_BUILD_TOUCHED_FIELD]: 'overriddenByLazyOverrides',
    [TEST_FACTORY_FIELDS.SUBFACTORY]: {
      [TEST_SUBFACTORY_FIELDS.POST_BUILD_TOUCHED_FIELD]: 'overriddenByLazyOverrides',
    },
  },
});

const testWithTrait = testFactory({ useTraits: ['testTrait'] });

const dutyLocationStreetAddressOverrideValue = '123 Main St.';

const dutyLocation = dutyLocationFactory({
  [BASE_FIELDS.OVERRIDES]: {
    [DUTY_LOCATION_FIELDS.ADDRESS]: {
      [ADDRESS_FIELDS.STREET_ADDRESS_1]: dutyLocationStreetAddressOverrideValue,
    },
  },
});

describe('testFactory', () => {
  it('has correct default values', () => {
    expect(test[TEST_FACTORY_FIELDS.DEFAULT]).toBe('default');
    expect(test[TEST_FACTORY_FIELDS.SUBFACTORY][TEST_SUBFACTORY_FIELDS.DEFAULT]).toBe('default');
    expect(test[TEST_FACTORY_FIELDS.SUBFACTORY][TEST_SUBFACTORY_FIELDS.FIELD_TO_OVERRIDE]).toBe('default');
  });
  it('has an overridden field', () => {
    expect(testWithOverriddenFieldViaOverride[TEST_FACTORY_FIELDS.DEFAULT]).toBe('overriddenByOverrides');
  });
  it('has overridden fields by postBuild', () => {
    expect(test[TEST_FACTORY_FIELDS.POST_BUILD_TOUCHED_FIELD]).toBe('overridden');
    expect(test[TEST_FACTORY_FIELDS.SUBFACTORY][TEST_SUBFACTORY_FIELDS.POST_BUILD_TOUCHED_FIELD]).toBe('overridden');
  });
  it('has an overridden field by its lazy overrides', () => {
    expect(testWithOverridenFieldByLazyOverrides[TEST_FACTORY_FIELDS.POST_BUILD_TOUCHED_FIELD]).toBe(
      'overriddenByLazyOverrides',
    );
  });
  it("has an overridden field by a subfactory's lazy overrides", () => {
    expect(
      testWithOverridenFieldByLazyOverrides[TEST_FACTORY_FIELDS.SUBFACTORY][
        TEST_SUBFACTORY_FIELDS.POST_BUILD_TOUCHED_FIELD
      ],
    ).toBe('overriddenByLazyOverrides');
  });
  it('has an overridden field by trait', () => {
    expect(testWithTrait[TEST_FACTORY_FIELDS.DEFAULT]).toBe('overriddenByTrait');
  });
  it("has an overridden field by a subfactory's overrides", () => {
    expect(
      testWithOverriddenFieldViaOverride[TEST_FACTORY_FIELDS.SUBFACTORY][TEST_SUBFACTORY_FIELDS.FIELD_TO_OVERRIDE],
    ).toBe('overriddenBySuboverrides');
  });
  it('maintains subfactory values that were not overridden', () => {
    expect(testWithOverriddenFieldViaOverride[TEST_FACTORY_FIELDS.SUBFACTORY][TEST_SUBFACTORY_FIELDS.DEFAULT]).toBe(
      'default',
    );
  });
  it('converts a field set with snake_case to a field with camelCase', () => {
    expect(test[TEST_FACTORY_FIELDS.CAMEL_CASE_FIELD]).toBe('caseFieldValue');
  });
  it("overrides a dutyLocation's address's street address", () => {
    expect(dutyLocation[DUTY_LOCATION_FIELDS.ADDRESS][ADDRESS_FIELDS.STREET_ADDRESS_1]).toBe(
      dutyLocationStreetAddressOverrideValue,
    );
  });
});
