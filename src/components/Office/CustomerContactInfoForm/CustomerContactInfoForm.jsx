import React from 'react';
import { Formik } from 'formik';
import * as Yup from 'yup';
import PropTypes from 'prop-types';
import { Checkbox } from '@trussworks/react-uswds';

import styles from './CustomerContactInfoForm.module.scss';

import { BackupContactInfoFields } from 'components/form/BackupContactInfoFields';
import { CustomerAltContactInfoFields } from 'components/form/CustomerAltContactInfoFields';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { requiredAddressSchema } from 'utils/validation';
import { ResidentialAddressShape } from 'types/address';

const CustomerContactInfoForm = ({ initialValues, onSubmit, onBack }) => {
  const validationSchema = Yup.object().shape({
    firstName: Yup.string().required('Required'),
    lastName: Yup.string().required('Required'),
    middleName: Yup.string(),
    suffix: Yup.string(),
    customerEmail: Yup.string()
      .matches(/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+.[a-zA-Z]{2,}$/, 'Must be a valid email address')
      .required('Required'),
    customerTelephone: Yup.string().min(12, 'Number must have 10 digits and a valid area code').required('Required'), // min 12 includes hyphens
    customerAddress: requiredAddressSchema.required(),
    name: Yup.string(),
    email: Yup.string().matches(/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+.[a-zA-Z]{2,}$/, 'Must be a valid email address'),
    telephone: Yup.string().min(12, 'Number must have 10 digits and a valid area code'), // min 12 includes hyphens
  });

  return (
    <Formik initialValues={initialValues} onSubmit={onSubmit} validationSchema={validationSchema} validateOnMount>
      {({ isValid, isSubmitting, handleSubmit }) => {
        return (
          <Form className={formStyles.form}>
            <SectionWrapper className={`${formStyles.formSection} ${styles.formSectionHeader}`}>
              <CustomerAltContactInfoFields
                render={(fields) => (
                  <>
                    <h2>Contact info</h2>
                    <Checkbox
                      data-testid="useCurrentResidence"
                      label="This is not the person named on the orders."
                      name="useCurrentResidence"
                      id="useCurrentResidenceCheckbox"
                    />
                    {fields}
                  </>
                )}
              />
              <h3 className={styles.sectionHeader}>Current Address</h3>
              <AddressFields name="customerAddress" />
            </SectionWrapper>
            <SectionWrapper className={`${formStyles.formSection} ${styles.formSectionHeader}`}>
              <h2 className={styles.sectionHeader}>
                Backup contact <span className={styles.optional}>Optional</span>
              </h2>

              <BackupContactInfoFields />
            </SectionWrapper>
            <div className={formStyles.formActions}>
              <WizardNavigation
                editMode
                disableNext={!isValid || isSubmitting}
                onCancelClick={onBack}
                onNextClick={handleSubmit}
              />
            </div>
          </Form>
        );
      }}
    </Formik>
  );
};

CustomerContactInfoForm.propTypes = {
  initialValues: PropTypes.shape({
    firstName: PropTypes.string,
    lastName: PropTypes.string,
    middleName: PropTypes.string,
    suffix: PropTypes.string,
    customerTelephone: PropTypes.string,
    customerEmail: PropTypes.string,
    name: PropTypes.string,
    telephone: PropTypes.string,
    email: PropTypes.string,
    customerAddress: ResidentialAddressShape,
  }).isRequired,
  onSubmit: PropTypes.func.isRequired,
  onBack: PropTypes.func.isRequired,
};

export default CustomerContactInfoForm;
