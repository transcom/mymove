import React from 'react';

import BackupMailingAddressForm from './BackupMailingAddressForm';

export default {
  title: 'Customer Components / Forms / BackupMailingAddressForm',
  component: BackupMailingAddressForm,
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
  <BackupMailingAddressForm
    formFieldsName="backup_mailing_address"
    initialValues={{
      backup_mailing_address: {
        street_address_1: '',
        street_address_2: '',
        city: '',
        state: '',
        postal_code: '',
      },
    }}
    onBack={argTypes.onBack}
    onSubmit={argTypes.onSubmit}
  />
);

export const WithInitialValues = (argTypes) => (
  <BackupMailingAddressForm
    formFieldsName="backup_mailing_address"
    initialValues={{
      backup_mailing_address: {
        street_address_1: '235 Prospect Valley Road SE',
        street_address_2: '',
        city: 'El Paso',
        state: 'TX',
        postal_code: '79912',
      },
    }}
    onBack={argTypes.onBack}
    onSubmit={argTypes.onSubmit}
  />
);
