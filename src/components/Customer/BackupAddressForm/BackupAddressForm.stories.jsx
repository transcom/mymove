import React from 'react';

import BackupAddressForm from './BackupAddressForm';

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

export const DefaultState = (argTypes) => (
  <BackupAddressForm
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
  <BackupAddressForm
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
