import React from 'react';
import { Formik } from 'formik';
import * as Yup from 'yup';

import { AddressFields } from './AddressFields';

import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import { requiredAddressSchema } from 'utils/validation';
import SectionWrapper from 'components/Customer/SectionWrapper';

export default {
  title: 'Components/Fieldsets/AddressFields',
};

export const Basic = () => (
  <Formik
    initialValues={{
      address: {
        streetAddress1: '',
        streetAddress2: '',
        city: '',
        state: '',
        postalCode: '',
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
        streetAddress1: '',
        streetAddress2: '',
        city: '',
        state: '',
        postalCode: '',
      },
    }}
    validationSchema={Yup.object().shape({ residential_address: requiredAddressSchema.required() })}
  >
    {() => (
      <Form className={formStyles.form}>
        <AddressFields legend="Pickup Address" name="residential_address" />
      </Form>
    )}
  </Formik>
);

export const CurrentResidentialAddressWithInitialValues = () => (
  <Formik
    initialValues={{
      residential_address: {
        streetAddress1: '123 Main St',
        streetAddress2: '#1A',
        city: 'New York',
        state: 'NY',
        postalCode: '10002',
      },
    }}
    validationSchema={Yup.object().shape({ residential_address: requiredAddressSchema.required() })}
  >
    {() => (
      <Form className={formStyles.form}>
        <AddressFields legend="Pickup Address" name="residential_address" />
      </Form>
    )}
  </Formik>
);

export const CurrentResidentialAddressWithCustomValidators = () => (
  <Formik
    initialValues={{
      residential_address: {
        streetAddress1: '',
        streetAddress2: '',
        city: '',
        state: '',
        postalCode: '',
      },
    }}
    validationSchema={Yup.object().shape({ residential_address: requiredAddressSchema.required() })}
  >
    {() => (
      <Form className={formStyles.form}>
        <AddressFields
          legend="Pickup Address"
          name="residential_address"
          validators={{
            city: (value) => (value === 'Nowhere' ? 'No one lives there' : ''),
            postalCode: (value) => (value !== '99999' ? 'ZIP code must be 99999' : ''),
          }}
        />
      </Form>
    )}
  </Formik>
);

export const InsideSectionWrapper = () => (
  <Formik
    initialValues={{
      residential_address: {
        streetAddress1: '',
        streetAddress2: '',
        city: '',
        state: '',
        postalCode: '',
      },
    }}
    validationSchema={Yup.object().shape({ residential_address: requiredAddressSchema.required() })}
  >
    {() => (
      <Form className={formStyles.form}>
        <SectionWrapper className={formStyles.formSection}>
          <AddressFields legend="Pickup Address" name="residential_address" />
        </SectionWrapper>
      </Form>
    )}
  </Formik>
);
