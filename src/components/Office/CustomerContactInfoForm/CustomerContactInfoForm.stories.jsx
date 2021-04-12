import React from 'react';

import CustomerContactInfoForm from './CustomerContactInfoForm';

export default {
  title: 'Office Components / Forms/ Customer Contact Info Form',
  component: CustomerContactInfoForm,
  argTypes: {
    onSubmit: { action: 'submitted on Next' },
    onBack: { action: 'to previous page on Back' },
  },
};

export const DefaultState = (argTypes) => (
  <CustomerContactInfoForm initialValues={{}} onSubmit={argTypes.onSubmit} onBack={argTypes.onBack} />
);

export const WithInitialValues = (argTypes) => (
  <CustomerContactInfoForm
    initialValues={{
      firstName: 'Leo',
      middleName: 'Star',
      lastName: 'Spaceman',
      suffix: 'Mr.',
      customerTelephone: '555-555-5555',
      customerEmail: 'test@sample.com',
      customerAddress: {
        street_address_1: '235 Prospect Valley Road SE',
        street_address_2: 'Apt. 3B',
        city: 'El Paso',
        state: 'TX',
        postal_code: '79912',
      },
      name: 'Leo Spaceman',
      telephone: '555-555-5555',
      email: 'test@sample.com',
    }}
    onNextClick={argTypes.onSubmit}
    onCancelClick={argTypes.onBack}
  />
);
