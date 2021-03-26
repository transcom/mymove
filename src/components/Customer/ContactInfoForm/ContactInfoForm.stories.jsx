import React from 'react';

import ContactInfoForm from './index';

export default {
  title: 'Customer Components / Forms/ Contact Info Form',
  component: ContactInfoForm,
  argTypes: {
    onSubmit: { action: 'submitted on Next' },
    onBack: { action: 'to previous page on Back' },
  },
};

const emptyInitialValues = {
  telephone: '',
  secondary_telephone: '',
  personal_email: '',
};

export const DefaultState = (argTypes) => (
  <ContactInfoForm initialValues={emptyInitialValues} onSubmit={argTypes.onSubmit} onBack={argTypes.onBack} />
);

export const WithInitialValues = (argTypes) => (
  <ContactInfoForm
    initialValues={{
      telephone: '555-555-5555',
      secondary_telephone: '555-444-5555',
      personal_email: 'test@sample.com',
      phone_is_preferred: false,
      email_is_preferred: true,
    }}
    onSubmit={argTypes.onSubmit}
    onBack={argTypes.onBack}
  />
);
