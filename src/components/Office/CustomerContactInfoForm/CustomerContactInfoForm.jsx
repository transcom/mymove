import React from 'react';
import { Field, Formik } from 'formik';
import * as Yup from 'yup';
import PropTypes from 'prop-types';
import { Checkbox, Radio, FormGroup, Grid } from '@trussworks/react-uswds';

import styles from './CustomerContactInfoForm.module.scss';

import { BackupContactInfoFields } from 'components/form/BackupContactInfoFields';
import { CustomerAltContactInfoFields } from 'components/form/CustomerAltContactInfoFields';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { phoneSchema, requiredAddressSchema } from 'utils/validation';
import { ResidentialAddressShape } from 'types/address';
import Hint from 'components/Hint';

const CustomerContactInfoForm = ({ initialValues, onSubmit, onBack }) => {
  const validationSchema = Yup.object().shape({
    firstName: Yup.string().required('Required'),
    lastName: Yup.string().required('Required'),
    middleName: Yup.string(),
    suffix: Yup.string(),
    customerEmail: Yup.string()
      .matches(/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+.[a-zA-Z]{2,}$/, 'Must be a valid email address')
      .required('Required'),
    customerTelephone: phoneSchema.required('Required'),
    secondaryPhone: phoneSchema,
    customerAddress: requiredAddressSchema.required(),
    backupAddress: requiredAddressSchema.required(),
    name: Yup.string().required('Required'),
    email: Yup.string()
      .matches(/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+.[a-zA-Z]{2,}$/, 'Must be a valid email address')
      .required('Required'),
    telephone: Yup.string()
      .min(12, 'Please enter a valid phone number. Phone numbers must be entered as ###-###-####.')
      .required('Required'), // min 12 includes hyphens
    phoneIsPreferred: Yup.boolean(),
    emailIsPreferred: Yup.boolean(),
    cacUser: Yup.boolean().required('Required'),
  });

  return (
    <Grid row>
      <Grid col>
        <div className={styles.customerContactForm}>
          <Formik initialValues={initialValues} onSubmit={onSubmit} validationSchema={validationSchema} validateOnMount>
            {({ isValid, handleSubmit, values, ...formikProps }) => {
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
                    <AddressFields name="customerAddress" values={values} locationLookup formikProps={formikProps} />
                    <h3 className={styles.sectionHeader}>Backup Address</h3>
                    <AddressFields name="backupAddress" values={values} locationLookup formikProps={formikProps} />
                  </SectionWrapper>
                  <SectionWrapper className={`${formStyles.formSection} ${styles.formSectionHeader}`}>
                    <h2 className={styles.sectionHeader}>Backup contact</h2>

                    <BackupContactInfoFields />
                  </SectionWrapper>
                  <SectionWrapper className={`${formStyles.formSection} ${styles.formSectionHeader}`}>
                    <h3>CAC Validation</h3>
                    <FormGroup>
                      <legend className="usa-label">
                        Is the customer a non-CAC user or do they need to bypass CAC validation?
                      </legend>
                      <Hint>
                        If this is checked yes, then the customer has already validated their CAC or their identity has
                        been validated by a trusted office user.
                      </Hint>
                      <div className="grid-row grid-gap">
                        <Field
                          as={Radio}
                          id="yesCacUser"
                          label="Yes"
                          name="cacUser"
                          value="true"
                          data-testid="cac-user-yes"
                          type="radio"
                        />
                        <Field
                          as={Radio}
                          id="NonCacUser"
                          label="No"
                          name="cacUser"
                          value="false"
                          data-testid="cac-user-no"
                          type="radio"
                        />
                      </div>
                    </FormGroup>
                  </SectionWrapper>
                  <div className={formStyles.formActions}>
                    <WizardNavigation
                      editMode
                      disableNext={!isValid}
                      onCancelClick={onBack}
                      onNextClick={handleSubmit}
                    />
                  </div>
                </Form>
              );
            }}
          </Formik>
        </div>
      </Grid>
    </Grid>
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
    cacUser: PropTypes.bool,
  }).isRequired,
  onSubmit: PropTypes.func,
  onBack: PropTypes.func,
};

CustomerContactInfoForm.defaultProps = {
  onSubmit: () => {},
  onBack: () => {},
};

export default CustomerContactInfoForm;
