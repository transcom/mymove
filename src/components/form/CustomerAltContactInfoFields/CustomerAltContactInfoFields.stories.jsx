import React from 'react';
import { Formik } from 'formik';

import { CustomerAltContactInfoFields } from './index';

import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';

export default {
  title: 'Components/Fieldsets/AltContactInfoFields',
};

export const DefaultState = () => (
  <Formik initialValues={{}}>
    {() => (
      <Form className={formStyles.form}>
        <CustomerAltContactInfoFields legend="Contact info" />
      </Form>
    )}
  </Formik>
);
export const WithInitialValues = () => (
  <Formik
    initialValues={{
      first_name: 'Leo',
      middle_name: 'Star',
      last_name: 'Spaceman',
      suffix: 'Mr.',
      customer_telephone: '555-555-5555',
      customer_email: 'test@sample.com',
    }}
  >
    {() => (
      <Form className={formStyles.form}>
        <CustomerAltContactInfoFields legend="Contact info" />
      </Form>
    )}
  </Formik>
);
