import React from 'react';

import EditContactInfoForm from './EditContactInfoForm';

export default {
  title: 'Customer Components / Forms / EditContactInfoForm',
  component: EditContactInfoForm,
  argTypes: {
    onSubmit: {
      action: 'save form',
    },
    onCancel: {
      action: 'cancel',
    },
  },
};

const fakeData = {
  telephone: '915-555-2942',
  secondary_telephone: '',
  personal_email: 'mm@example.com',
  phone_is_preferred: false,
  email_is_preferred: true,
  residential_address: {
    streetAddress1: '235 Prospect Valley Road SE',
    streetAddress2: '',
    city: 'El Paso',
    state: 'TX',
    postalCode: '79912',
  },
  backup_mailing_address: {
    streetAddress1: '9 W 2nd Ave',
    streetAddress2: '',
    city: 'El Paso',
    state: 'TX',
    postalCode: '79936',
  },
  backup_contact: {
    name: 'Peyton Wing',
    email: 'pw@example.com',
    telephone: '915-555-8761',
  },
};

export const DefaultState = (argTypes) => (
  <EditContactInfoForm
    initialValues={{
      telephone: '',
      secondary_telephone: '',
      personal_email: '',
      residential_address: {
        streetAddress1: '',
        streetAddress2: '',
        city: '',
        state: '',
        postalCode: '',
      },
      backup_mailing_address: {
        streetAddress1: '',
        streetAddress2: '',
        city: '',
        state: '',
        postalCode: '',
      },
      backup_contact: {
        name: '',
        email: '',
        telephone: '',
      },
    }}
    onCancel={argTypes.onCancel}
    onSubmit={argTypes.onSubmit}
  />
);

export const WithInitialValues = (argTypes) => (
  <EditContactInfoForm initialValues={fakeData} onCancel={argTypes.onCancel} onSubmit={argTypes.onSubmit} />
);
