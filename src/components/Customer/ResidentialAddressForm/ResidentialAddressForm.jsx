import React from 'react';
import PropTypes from 'prop-types';
import { Formik } from 'formik';
import classnames from 'classnames';
import * as Yup from 'yup';

import styles from './ResidentialAddressForm.module.scss';

import { AddressFields } from 'components/form/AddressFields/AddressFields';
import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import { requiredAddressSchema } from 'utils/validation';
import { ResidentialAddressShape } from 'types/address';

const ResidentialAddressForm = ({ formFieldsName, initialValues, onSubmit, onBack }) => {
  const validationSchema = Yup.object().shape({
    [formFieldsName]: requiredAddressSchema.required(),
  });

  return (
    <Formik
      initialValues={initialValues}
      onSubmit={onSubmit}
      validateOnMount
      validateOnBlur
      validateOnChange
      validationSchema={validationSchema}
    >
      {({ isValid, isSubmitting, handleSubmit, values, ...formikProps }) => {
        return (
          <Form className={formStyles.form}>
            <h1>Current address</h1>
            <p className={styles.noBottomMargin}>Must be a physical address.</p>
            <SectionWrapper className={classnames(styles.noTopMargin, formStyles.formSection)}>
              <AddressFields labelHint="Required" name={formFieldsName} formikProps={formikProps} />
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

ResidentialAddressForm.propTypes = {
  formFieldsName: PropTypes.string.isRequired,
  initialValues: ResidentialAddressShape.isRequired,
  onBack: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
};

ResidentialAddressForm.defaultProps = {};

export default ResidentialAddressForm;
