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
  <BackupMailingAddressForm
    formFieldsName="backup_mailing_address"
    initialValues={{
      backup_mailing_address: {
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
