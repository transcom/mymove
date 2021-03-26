import React from 'react';
import PropTypes from 'prop-types';
import { Formik } from 'formik';
import * as Yup from 'yup';

import { AddressFields } from 'components/form/AddressFields/AddressFields';
import SectionWrapper from 'components/Customer/SectionWrapper';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import { requiredAddressSchema } from 'utils/validation';

const ResidentialAddressForm = ({ formFieldsName, initialValues, onSubmit, validators }) => {
  const validationSchema = Yup.object().shape({
    [formFieldsName]: requiredAddressSchema.required(),
  });

  return (
    <Formik
      initialValues={initialValues}
      onSubmit={onSubmit}
      validateOnChange={false}
      validateOnMount
      validationSchema={validationSchema}
    >
      {({ isValid, isSubmitting, handleSubmit, setFieldValue, errors }) => {
        const handleBack = (e) => {
          setFieldValue('nextPage', 'back');
          handleSubmit(e);
        };

        const handleNext = (e) => {
          setFieldValue('nextPage', 'next');
          handleSubmit(e);
        };

        return (
          <Form className={formStyles.form}>
            <h1>Current residence</h1>
            <p>
              Valid: {isValid.toString()} | Submitting: {isSubmitting.toString()} | Errors: {JSON.stringify(errors)}
            </p>

            <SectionWrapper className={formStyles.formSection}>
              <AddressFields name={formFieldsName} validators={validators} />
            </SectionWrapper>

            <div className={formStyles.formActions}>
              <WizardNavigation
                onBackClick={handleBack}
                disableNext={!isValid || isSubmitting}
                onNextClick={handleNext}
              />
            </div>
          </Form>
        );
      }}
    </Formik>
  );
};

ResidentialAddressForm.propTypes = {
  formFieldsName: PropTypes.string.isRequired,
  initialValues: PropTypes.shape({
    street_address_1: PropTypes.string,
    street_address_2: PropTypes.string,
    city: PropTypes.string,
    state: PropTypes.string,
    postal_code: PropTypes.string,
  }).isRequired,
  onSubmit: PropTypes.func.isRequired,
  validators: PropTypes.shape({
    streetAddress1: PropTypes.func,
    streetAddress2: PropTypes.func,
    city: PropTypes.func,
    state: PropTypes.func,
    postalCode: PropTypes.func,
  }),
};

ResidentialAddressForm.defaultProps = {
  validators: {},
};

export default ResidentialAddressForm;
