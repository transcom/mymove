import React from 'react';
import { Formik } from 'formik';
import * as Yup from 'yup';

import { AddressFields } from './AddressFields';

import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import { requiredAddressSchema } from 'utils/validation';

export default {
  title: 'Components/Fieldsets/AddressFields',
};

export const Basic = () => (
  <Formik
    initialValues={{
      address: {
        street_address_1: '',
        street_address_2: '',
        city: '',
        state: '',
        postal_code: '',
      },
    }}
  >
    {() => (
      <Form className={formStyles.form}>
        <AddressFields legend="Basic address" name="address" />
      </Form>
    )}
  </Formik>
);

export const CurrentResidentialAddress = () => (
  <Formik
    initialValues={{
      residential_address: {
        street_address_1: '',
        street_address_2: '',
        city: '',
        state: '',
        postal_code: '',
      },
    }}
    validationSchema={Yup.object().shape({ residential_address: requiredAddressSchema.required() })}
  >
    {() => (
      <Form className={formStyles.form}>
        <AddressFields legend="Current mailing address" name="residential_address" />
      </Form>
    )}
  </Formik>
);

export const CurrentResidentialAddressWithInitialValues = () => (
  <Formik
    initialValues={{
      residential_address: {
        street_address_1: '123 Main St',
        street_address_2: '#1A',
        city: 'New York',
        state: 'NY',
        postal_code: '10002',
      },
    }}
    validationSchema={Yup.object().shape({ residential_address: requiredAddressSchema.required() })}
  >
    {() => (
      <Form className={formStyles.form}>
        <AddressFields legend="Current mailing address" name="residential_address" />
      </Form>
    )}
  </Formik>
);
