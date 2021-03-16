import React from 'react';
import { action } from '@storybook/addon-actions';
import { Formik } from 'formik';

import { ContactInfoFieldset } from './index';

export default {
  title: 'Components/ContactInfoFieldset',
};

export const ContactInfoFieldsetBasic = () => (
  <Formik
    initialValues={{
      contact: {
        phone: '',
        alternatePhone: '',
        email: '',
      },
    }}
  >
    <ContactInfoFieldset
      name="contact"
      legend="Your contact info"
      onChangePreferPhone={action('clicked')}
      onChangePreferEmail={action('clicked')}
    />
  </Formik>
);
