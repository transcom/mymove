import React from 'react';
import { action } from '@storybook/addon-actions';
import { Formik } from 'formik';

import { ServiceMemberContactInfoFields } from './index';

export default {
  title: 'Components/Fieldsets/ServiceMemberContactInfoFields',
};

export const Basic = () => (
  <Formik
    initialValues={{
      contact: {
        phone: '',
        alternatePhone: '',
        email: '',
      },
    }}
  >
    <ServiceMemberContactInfoFields
      name="contact"
      legend="Your contact info"
      onChangePreferPhone={action('clicked')}
      onChangePreferEmail={action('clicked')}
    />
  </Formik>
);
