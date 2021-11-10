import React from 'react';

import ResidentialAddressForm from './ResidentialAddressForm';

export default {
  title: 'Customer Components / Forms / ResidentialAddressForm',
  component: ResidentialAddressForm,
  argTypes: {
    onSubmit: {
      action: 'submit form',
    },
    onBack: {
      action: 'go back',
    },
  },
};

export const DefaultState = (argTypes) => (
  <ResidentialAddressForm
    formFieldsName="residential_address"
    initialValues={{
      residential_address: {
        streetAddress1: '',
        streetAddress2: '',
        city: '',
        state: '',
        postalCode: '',
      },
    }}
    onBack={argTypes.onBack}
    onSubmit={argTypes.onSubmit}
  />
);

export const WithInitialValues = (argTypes) => (
  <ResidentialAddressForm
    formFieldsName="residential_address"
    initialValues={{
      residential_address: {
        streetAddress1: '235 Prospect Valley Road SE',
        streetAddress2: '',
        city: 'El Paso',
        state: 'TX',
        postalCode: '79912',
      },
    }}
    onBack={argTypes.onBack}
    onSubmit={argTypes.onSubmit}
  />
);

export const WithCustomValidators = (argTypes) => (
  <ResidentialAddressForm
    formFieldsName="residential_address"
    initialValues={{
      residential_address: {
        streetAddress1: '',
        streetAddress2: '',
        city: '',
        state: '',
        postalCode: '',
      },
    }}
    onBack={argTypes.onBack}
    onSubmit={argTypes.onSubmit}
    validators={{
      city: (value) => (value === 'Nowhere' ? 'No one lives there' : ''),
      postalCode: (value) => (value !== '99999' ? 'ZIP code must be 99999' : ''),
    }}
  />
);
