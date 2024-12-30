import React from 'react';
import { Provider } from 'react-redux';

import BackupAddressForm from './BackupAddressForm';

import { configureStore } from 'shared/store';

export default {
  title: 'Customer Components / Forms / BackupAddressForm',
  component: BackupAddressForm,
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
    <BackupAddressForm
      formFieldsName="backup_mailing_address"
      initialValues={{
        backup_mailing_address: {
          streetAddress1: '',
          streetAddress2: '',
          city: '',
          state: '',
          postalCode: '',
          county: '',
        },
      }}
      onBack={argTypes.onBack}
      onSubmit={argTypes.onSubmit}
    />
  </Provider>
);

export const WithInitialValues = (argTypes) => (
  <Provider store={mockStore.store}>
    <BackupAddressForm
      formFieldsName="backup_mailing_address"
      initialValues={{
        backup_mailing_address: {
          streetAddress1: '235 Prospect Valley Road SE',
          streetAddress2: '',
          city: 'El Paso',
          state: 'TX',
          postalCode: '79912',
          county: 'EL PASO',
        },
      }}
      onBack={argTypes.onBack}
      onSubmit={argTypes.onSubmit}
    />
  </Provider>
);
