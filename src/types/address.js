/* eslint-disable import/prefer-default-export */
import { shape, string } from 'prop-types';

export const SimpleAddressShape = shape({
  city: string,
  state: string,
  postalCode: string,
  usprcID: string,
});

export const MandatorySimpleAddressShape = shape({
  city: string.isRequired,
  state: string.isRequired,
  postalCode: string.isRequired,
  usprcID: string,
});

export const AddressShape = shape({
  city: string,
  state: string,
  postalCode: string,
  streetAddress1: string,
  streetAddress2: string,
  streetAddress3: string,
  country: string,
  usprcID: string,
});

export const ResidentialAddressShape = shape({
  streetAddress1: string,
  streetAddress2: string,
  city: string,
  state: string,
  postalCode: string,
  usprcID: string,
});

export const W2AddressShape = shape({
  streetAddress1: string,
  streetAddress2: string,
  city: string,
  state: string,
  postalCode: string,
  usprcID: string,
});
