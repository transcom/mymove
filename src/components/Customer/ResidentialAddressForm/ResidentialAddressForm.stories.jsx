import React from 'react';

import ResidentialAddressForm from './ResidentialAddressForm';

export default {
  title: 'Customer Components / Forms / ResidentialAddressForm',
  component: ResidentialAddressForm,
  argTypes: {
    onSubmit: {
      action: 'submit form',
    },
  },
};

export const DefaultState = (argTypes) => (
  <ResidentialAddressForm
    initialValues={{
      residential_address: {
        street_address_1: '',
        street_address_2: '',
        city: '',
        state: '',
        postal_code: '',
      },
    }}
    onSubmit={argTypes.onSubmit}
  />
);

export const WithCustomValidators = (argTypes) => (
  <ResidentialAddressForm
    initialValues={{
      residential_address: {
        street_address_1: '',
        street_address_2: '',
        city: '',
        state: '',
        postal_code: '',
      },
    }}
    onSubmit={argTypes.onSubmit}
    validators={{
      city: (value) => (value === 'Nowhere' ? 'No one lives there' : ''),
      postalCode: (value) => (value !== '99999' ? 'ZIP code must be 99999' : ''),
    }}
  />
);
