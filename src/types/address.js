/* eslint-disable import/prefer-default-export */
import { shape, string } from 'prop-types';

export const SimpleAddressShape = shape({
  city: string,
  state: string,
  postal_code: string,
});

export const AddressShape = shape({
  city: string,
  state: string,
  postal_code: string,
  street_address_1: string,
  street_address_2: string,
  street_address_3: string,
  country: string,
});

export const ResidentialAddressShape = shape({
  street_address_1: string,
  street_address_2: string,
  city: string,
  state: string,
  postal_code: string,
});
