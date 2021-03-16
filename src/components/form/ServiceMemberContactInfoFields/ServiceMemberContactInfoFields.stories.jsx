import React from 'react';
import { action } from '@storybook/addon-actions';
import { Formik } from 'formik';

import { ServiceMemberContactInfoFields } from './index';

import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';

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
    {() => (
      <Form className={formStyles.form}>
        <ServiceMemberContactInfoFields
          name="contact"
          legend="Your contact info"
          onChangePreferPhone={action('clicked')}
          onChangePreferEmail={action('clicked')}
        />
      </Form>
    )}
  </Formik>
);
