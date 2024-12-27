import React from 'react';
import { Provider } from 'react-redux';

import ResidentialAddressForm from './ResidentialAddressForm';

import { configureStore } from 'shared/store';

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

const mockStore = configureStore({});

export const DefaultState = (argTypes) => (
  <Provider store={mockStore.store}>
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
  </Provider>
);

export const WithInitialValues = (argTypes) => (
  <Provider store={mockStore.store}>
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
  </Provider>
);

export const WithCustomValidators = (argTypes) => (
  <Provider store={mockStore.store}>
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
  </Provider>
);
