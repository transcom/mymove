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
  <div className="officeApp">
    <CustomerContactInfoForm initialValues={{}} onSubmit={argTypes.onSubmit} onBack={argTypes.onBack} />
  </div>
);

export const WithInitialValues = (argTypes) => (
  <div className="officeApp">
    <CustomerContactInfoForm
      initialValues={{
        firstName: 'Leo',
        middleName: 'Star',
        lastName: 'Spaceman',
        suffix: 'Mr.',
        customerTelephone: '555-555-5555',
        customerEmail: 'test@sample.com',
        customerAddress: {
          streetAddress1: '235 Prospect Valley Road SE',
          streetAddress2: 'Apt. 3B',
          city: 'El Paso',
          state: 'TX',
          postalCode: '79912',
        },
        name: 'Leo Spaceman',
        telephone: '555-555-5555',
        email: 'test@sample.com',
      }}
      onNextClick={argTypes.onSubmit}
      onCancelClick={argTypes.onBack}
    />
  </div>
);
