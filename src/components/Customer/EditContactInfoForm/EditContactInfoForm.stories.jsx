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
    street_address_1: '235 Prospect Valley Road SE',
    street_address_2: '',
    city: 'El Paso',
    state: 'TX',
    postal_code: '79912',
  },
  backup_mailing_address: {
    street_address_1: '9 W 2nd Ave',
    street_address_2: '',
    city: 'El Paso',
    state: 'TX',
    postal_code: '79936',
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
        street_address_1: '',
        street_address_2: '',
        city: '',
        state: '',
        postal_code: '',
      },
      backup_mailing_address: {
        street_address_1: '',
        street_address_2: '',
        city: '',
        state: '',
        postal_code: '',
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
