import React from 'react';

import W2AddressForm from './W2AddressForm';

export default {
  title: 'Customer Components / PPM Closeout /  W2AddressForm',
  component: W2AddressForm,
};

export const DefaultState = () => (
  <W2AddressForm
    formFieldsName="w2_address"
    initialValues={{
      w2_address: {
        streetAddress1: '',
        streetAddress2: '',
        city: '',
        state: '',
        postalCode: '',
      },
    }}
  />
);

export const WithInitialValues = () => (
  <W2AddressForm
    formFieldsName="w2_address"
    initialValues={{
      w2_address: {
        streetAddress1: '123 Anywhere St',
        streetAddress2: '',
        city: 'Santa Monica',
        state: 'CA',
        postalCode: '90402',
      },
    }}
  />
);

export const WithCustomValidators = () => (
  <W2AddressForm
    formFieldsName="w2_address"
    initialValues={{
      w2_address: {
        streetAddress1: '',
        streetAddress2: '',
        city: '',
        state: '',
        postalCode: '',
      },
    }}
    validators={{
      city: (value) => (value === 'Nowhere' ? 'No one lives there' : ''),
      postalCode: (value) => (value !== '99999' ? 'ZIP code must be 99999' : ''),
    }}
  />
);
