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

export const CurrentResidentialAddressWithCustomValidators = () => (
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
        <AddressFields
          legend="Current mailing address"
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

export const WithAdditionalText = () => (
  <Formik
    initialValues={{
      backup_mailing_address: {
        street_address_1: '',
        street_address_2: '',
        city: '',
        state: '',
        postal_code: '',
      },
    }}
    validationSchema={Yup.object().shape({ backup_mailing_address: requiredAddressSchema.required() })}
  >
    {() => (
      <Form className={formStyles.form}>
        <AddressFields
          legend="Backup mailing address"
          name="backup_mailing_address"
          render={(fields) => (
            <>
              <p>
                Where should we send mail if we can’t reach you at your primary address? You might use a parent’s or
                friend’s address, or a post office box.
              </p>
              {fields}
            </>
          )}
        />
      </Form>
    )}
  </Formik>
);

export const InsideSectionWrapper = () => (
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
        <SectionWrapper className={formStyles.formSection}>
          <AddressFields legend="Current mailing address" name="residential_address" />
        </SectionWrapper>
      </Form>
    )}
  </Formik>
);
