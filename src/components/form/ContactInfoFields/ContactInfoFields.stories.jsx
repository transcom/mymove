import React from 'react';
import { Formik } from 'formik';

import { ContactInfoFields } from './ContactInfoFields';

import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';

export default {
  title: 'Components/Fieldsets/ContactInfoFields',
};

export const Basic = () => (
  <Formik
    initialValues={{
      contact: {
        firstName: '',
        lastName: '',
        phone: '',
        email: '',
      },
    }}
  >
    {() => (
      <Form className={formStyles.form}>
        <ContactInfoFields name="contact" legend="Contact Info" />
      </Form>
    )}
  </Formik>
);
