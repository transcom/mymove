/* eslint-disable camelcase */
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
    street_address_1: Yup.string().when(
      ['street_address_2', 'city', 'state', 'postal_code'],
      (street2, city, state, postalCode, schema) =>
        street2 || city || state || postalCode ? schema.required('Required') : schema,
    ),
    street_address_2: Yup.string(),
    city: Yup.string().when(
      ['street_address_1', 'street_address_2', 'state', 'postal_code'],
      (street1, street2, state, postalCode, schema) =>
        street1 || street2 || state || postalCode ? schema.required('Required') : schema,
    ),
    state: Yup.string()
      .length(2, 'Must use state abbreviation')
      .when(
        ['street_address_1', 'street_address_2', 'city', 'postal_code'],
        (street1, street2, city, postalCode, schema) =>
          street1 || street2 || city || postalCode ? schema.required('Required') : schema,
      ),
    postal_code: Yup.string()
      .matches(ZIP_CODE_REGEX, 'Must be valid zip code')
      .when(['street_address_1', 'street_address_2', 'city', 'state'], (street1, street2, city, state, schema) =>
        street1 || street2 || city || state ? schema.required('Required') : schema,
      ),
  },
  [
    ['street_address_1', 'street_address_2'],
    ['street_address_1', 'city'],
    ['street_address_1', 'state'],
    ['street_address_1', 'postal_code'],
    ['city', 'state'],
    ['city', 'postal_code'],
    ['state', 'postal_code'],
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
