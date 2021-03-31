import React from 'react';

import BackupContactForm from './index';

export default {
  title: 'Customer Components / Forms/ Backup Contact Form',
  component: BackupContactForm,
  argTypes: {
    onSubmit: { action: 'submitted on Next' },
    onBack: { action: 'to previous page on Back' },
  },
};

const emptyInitialValues = {
  name: '',
  telephone: '',
  email: '',
};

export const DefaultState = (argTypes) => (
  <BackupContactForm initialValues={emptyInitialValues} onSubmit={argTypes.onSubmit} onBack={argTypes.onBack} />
);

export const WithInitialValues = (argTypes) => (
  <BackupContactForm
    initialValues={{
      name: 'Leo Spaceman',
      telephone: '555-555-5555',
      email: 'test@sample.com',
    }}
    onSubmit={argTypes.onSubmit}
    onBack={argTypes.onBack}
  />
);
