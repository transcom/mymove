import React from 'react';
import { Formik } from 'formik';
import * as Yup from 'yup';

import { AddressFields } from './AddressFields';

import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';

export default {
  title: 'Components/Fieldsets/AddressFields',
};

const requiredAddressSchema = Yup.object().shape({
  street_address_1: Yup.string().required('Required'),
  street_address_2: Yup.string(),
  city: Yup.string().required('Required'),
  state: Yup.string().length(2, 'Must use state abbreviation').required('Required'),
  postal_code: Yup.string()
    // RA Summary: eslint - security/detect-unsafe-regex - Denial of Service: Regular Expression
    // RA: Locates potentially unsafe regular expressions, which may take a very long time to run, blocking the event loop
    // RA: Per MilMove SSP, predisposing conditions are regex patterns from untrusted sources or unbounded matching.
    // RA: The regex pattern is a constant string set at compile-time and it is bounded to 10 characters (zip code).
    // RA Developer Status: Mitigated
    // RA Validator Status:  Mitigated
    // RA Modified Severity: N/A
    // eslint-disable-next-line security/detect-unsafe-regex
    .matches(/^(\d{5}([-]\d{4})?)$/, 'Must be valid zip code')
    .required('Required'),
});

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
    validationSchema={requiredAddressSchema}
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
    validationSchema={requiredAddressSchema}
  >
    {() => (
      <Form className={formStyles.form}>
        <AddressFields legend="Current mailing address" name="residential_address" />
      </Form>
    )}
  </Formik>
);
