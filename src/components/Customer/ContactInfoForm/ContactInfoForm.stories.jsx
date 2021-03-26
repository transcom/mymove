/* eslint-disable no-console */
import React from 'react';

import ContactInfoForm from './index';

export default {
  title: 'Customer Components / Forms/ Contact Info Form',
  component: ContactInfoForm,
};

const initialValues = {
  telephone: '',
  secondary_telephone: '',
  personal_email: '',
};
const onSubmit = () => console.log('submitted on Next');
const onBack = () => console.log('Back to previous page');
export const DefaultState = () => (
  <ContactInfoForm initialValues={initialValues} onSubmit={() => {}} onBack={() => {}} />
);

export const WithInitialValues = () => (
  <ContactInfoForm
    initialValues={{
      telephone: '555-555-5555',
      secondary_telephone: '555-444-5555',
      personal_email: 'test@sample.com',
      phone_is_preferred: false,
      email_is_preferred: true,
    }}
    onSubmit={onSubmit}
    onBack={onBack}
  />
);
