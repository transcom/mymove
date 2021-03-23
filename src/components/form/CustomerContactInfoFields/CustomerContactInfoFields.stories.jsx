import React from 'react';
import { Formik } from 'formik';

import { CustomerContactInfoFields } from './index';

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
        <CustomerContactInfoFields name="contact" legend="Your contact info" />
      </Form>
    )}
  </Formik>
);
