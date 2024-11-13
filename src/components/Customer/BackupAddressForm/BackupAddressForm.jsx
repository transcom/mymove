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
import { ResidentialAddressShape } from 'types/address';

const BackupAddressForm = ({ formFieldsName, initialValues, onSubmit, onBack }) => {
  const validationSchema = Yup.object().shape({
    [formFieldsName]: requiredAddressSchema.required(),
  });

  return (
    <Formik
      initialValues={initialValues}
      onSubmit={onSubmit}
      validateOnChange
      validateOnBlur
      validateOnMount
      validationSchema={validationSchema}
    >
      {({ isValid, isSubmitting, handleSubmit, validateForm, ...formikProps }) => {
        return (
          <Form className={formStyles.form}>
            <h1>Backup address</h1>

            <p>
              Provide a physical address where either you can be reached or someone can contact you while you are in
              transit during your move.
            </p>

            <SectionWrapper className={formStyles.formSection}>
              <AddressFields
                labelHint="Required"
                name={formFieldsName}
                locationLookup
                validateForm={validateForm}
                formikProps={formikProps}
              />
            </SectionWrapper>

            <div className={formStyles.formActions}>
              <WizardNavigation
                onBackClick={onBack}
                disableNext={!isValid || isSubmitting}
                onNextClick={handleSubmit}
              />
            </div>
          </Form>
        );
      }}
    </Formik>
  );
};

BackupAddressForm.propTypes = {
  formFieldsName: PropTypes.string.isRequired,
  initialValues: ResidentialAddressShape.isRequired,
  onBack: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
};

export default BackupAddressForm;
