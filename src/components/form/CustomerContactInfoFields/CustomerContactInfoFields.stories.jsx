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
      telephone: '',
      secondary_telephone: '',
      personal_email: '',
    }}
  >
    {() => (
      <Form className={formStyles.form}>
        <CustomerContactInfoFields legend="Your contact info" />
      </Form>
    )}
  </Formik>
);
