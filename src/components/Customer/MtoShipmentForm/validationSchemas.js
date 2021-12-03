import * as Yup from 'yup';

import { requiredAddressSchema, ZIP_CODE_REGEX } from 'utils/validation';

export const AgentSchema = Yup.object().shape({
  firstName: Yup.string(),
  lastName: Yup.string(),
  phone: Yup.string().matches(/^[2-9]\d{2}-\d{3}-\d{4}$/, 'Must be valid phone number'),
  email: Yup.string().email('Must be valid email'),
});

export const OptionalAddressSchema = Yup.object().shape(
  {
    streetAddress1: Yup.string().when(
      ['streetAddress2', 'city', 'state', 'postalCode'],
      (street2, city, state, postalCode, schema) =>
        street2 || city || state || postalCode ? schema.required('Required') : schema,
    ),
    streetAddress2: Yup.string(),
    city: Yup.string().when(
      ['streetAddress1', 'streetAddress2', 'state', 'postalCode'],
      (street1, street2, state, postalCode, schema) =>
        street1 || street2 || state || postalCode ? schema.required('Required') : schema,
    ),
    state: Yup.string()
      .length(2, 'Must use state abbreviation')
      .when(['streetAddress1', 'streetAddress2', 'city', 'postalCode'], (street1, street2, city, postalCode, schema) =>
        street1 || street2 || city || postalCode ? schema.required('Required') : schema,
      ),
    postalCode: Yup.string()
      .matches(ZIP_CODE_REGEX, 'Must be valid zip code')
      .when(['streetAddress1', 'streetAddress2', 'city', 'state'], (street1, street2, city, state, schema) =>
        street1 || street2 || city || state ? schema.required('Required') : schema,
      ),
  },
  [
    ['streetAddress1', 'streetAddress2'],
    ['streetAddress1', 'city'],
    ['streetAddress1', 'state'],
    ['streetAddress1', 'postalCode'],
    ['city', 'state'],
    ['city', 'postalCode'],
    ['state', 'postalCode'],
  ],
);

export const RequiredPlaceSchema = Yup.object().shape({
  address: requiredAddressSchema,
  agent: AgentSchema,
});

export const OptionalPlaceSchema = Yup.object().shape({
  address: OptionalAddressSchema,
  agent: AgentSchema,
});

export const AdditionalAddressSchema = Yup.object().shape({
  address: OptionalAddressSchema,
});

export const StorageFacilityAddressSchema = Yup.object().shape({
  address: requiredAddressSchema,
  lotNumber: Yup.string(),
  facilityName: Yup.string().required('Required'),
  phone: Yup.string().matches(/^[2-9]\d{2}-\d{3}-\d{4}$/, 'Must be valid phone number'),
  email: Yup.string().email('Must be valid email'),
});
