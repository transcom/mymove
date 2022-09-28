import React from 'react';
import PropTypes from 'prop-types';
import { Formik } from 'formik';
import * as Yup from 'yup';

import { AddressFields } from 'components/form/AddressFields/AddressFields';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import { requiredAddressSchema } from 'utils/validation';
import { W2AddressShape } from 'types/address';

const W2AddressForm = ({ formFieldsName, initialValues, validators }) => {
  const validationSchema = Yup.object().shape({
    [formFieldsName]: requiredAddressSchema.required(),
  });

  return (
    <Formik initialValues={initialValues} validateOnChange={false} validateOnMount validationSchema={validationSchema}>
      <Form className={formStyles.form}>
        <SectionWrapper className={formStyles.formSection}>
          <h2>W-2 address</h2>
          <p>What is the address on your W-2?</p>
          <AddressFields name={formFieldsName} validators={validators} />
        </SectionWrapper>
      </Form>
    </Formik>
  );
};

W2AddressForm.propTypes = {
  formFieldsName: PropTypes.string.isRequired,
  initialValues: W2AddressShape.isRequired,
  validators: PropTypes.shape({
    streetAddress1: PropTypes.func,
    streetAddress2: PropTypes.func,
    city: PropTypes.func,
    state: PropTypes.func,
    postalCode: PropTypes.func,
  }),
};

W2AddressForm.defaultProps = {
  validators: {},
};

export default W2AddressForm;
