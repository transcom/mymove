/* eslint-disable import/prefer-default-export */

import * as Yup from 'yup';

import { ZIP_CODE_REGEX, IsSupportedState, UnsupportedStateErrorMsg } from 'utils/validation';

export const OptionalAddressSchema = Yup.object().shape(
  {
    streetAddress1: Yup.string().when(
      ['streetAddress2', 'city', 'state', 'postalCode'],
      ([street2, city, state, postalCode], schema) =>
        street2 || city || state || postalCode ? schema.required('Required') : schema,
    ),
    streetAddress2: Yup.string(),
    city: Yup.string().when(
      ['streetAddress1', 'streetAddress2', 'state', 'postalCode'],
      ([street1, street2, state, postalCode], schema) =>
        street1 || street2 || state || postalCode ? schema.required('Required') : schema,
    ),
    state: Yup.string()
      .test('', UnsupportedStateErrorMsg, IsSupportedState)
      .length(2, 'Must use state abbreviation')
      .when(
        ['streetAddress1', 'streetAddress2', 'city', 'postalCode'],
        ([street1, street2, city, postalCode], schema) =>
          street1 || street2 || city || postalCode ? schema.required('Required') : schema,
      ),
    postalCode: Yup.string()
      .matches(ZIP_CODE_REGEX, 'Must be valid zip code')
      .when(['streetAddress1', 'streetAddress2', 'city', 'state'], ([street1, street2, city, state], schema) =>
        street1 || street2 || city || state ? schema.required('Required') : schema,
      ),
    countryID: Yup.string().when('[streetAddress1, city]', ([street1, city], schema) =>
      street1 || city ? schema.required('Required') : schema,
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
