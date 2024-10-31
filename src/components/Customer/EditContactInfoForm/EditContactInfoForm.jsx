import classnames from 'classnames';
import React from 'react';
import PropTypes from 'prop-types';
import { Formik } from 'formik';
import * as Yup from 'yup';

import editContactInfoFormStyle from './EditContactInfoForm.module.scss';

import { AddressFields } from 'components/form/AddressFields/AddressFields';
import SectionWrapper from 'components/Customer/SectionWrapper';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import {
  backupContactInfoSchema,
  contactInfoSchema,
  requiredAddressSchema,
  preferredContactMethodValidation,
} from 'utils/validation';
import { ResidentialAddressShape } from 'types/address';
import { CustomerContactInfoFields } from 'components/form/CustomerContactInfoFields';
import { BackupContactInfoFields } from 'components/form/BackupContactInfoFields';
import { BackupContactShape } from 'types/customerShapes';

export const residentialAddressName = 'residential_address';
export const backupAddressName = 'backup_mailing_address';
export const backupContactName = 'backup_contact';

const EditContactInfoForm = ({ initialValues, onSubmit, onCancel }) => {
  const validationSchema = Yup.object()
    .shape({
      ...contactInfoSchema.fields,
      [residentialAddressName]: requiredAddressSchema.required(),
      [backupAddressName]: requiredAddressSchema.required(),
      [backupContactName]: backupContactInfoSchema.required(),
    })
    .test('contactMethodRequired', 'Please select a preferred method of contact.', preferredContactMethodValidation);

  const sectionStyles = classnames(formStyles.formSection, editContactInfoFormStyle.formSection);

  return (
    <Formik initialValues={initialValues} onSubmit={onSubmit} validateOnMount validationSchema={validationSchema}>
      {({ isValid, isSubmitting, handleSubmit, values, setValues }) => {
        const handleCurrentZipCityChange = (value) => {
          setValues(
            {
              ...values,
              residential_address: {
                ...values.residential_address,
                city: value.city ? value.city : '',
                state: value.state ? value.state : '',
                county: value.county ? value.county : '',
                postalCode: value.postalCode ? value.postalCode : '',
              },
            },
            { shouldValidate: true },
          );
        };
        const handleBackupZipCityChange = (value) => {
          setValues(
            {
              ...values,
              backup_mailing_address: {
                ...values.backup_mailing_address,
                city: value.city ? value.city : '',
                state: value.state ? value.state : '',
                county: value.county ? value.county : '',
                postalCode: value.postalCode ? value.postalCode : '',
              },
            },
            { shouldValidate: true },
          );
        };
        return (
          <Form className={classnames(formStyles.form, editContactInfoFormStyle.form)}>
            <h1>Edit contact info</h1>

            <SectionWrapper className={sectionStyles}>
              <h2>Your contact info</h2>

              <CustomerContactInfoFields labelHint="Required" />
            </SectionWrapper>

            <SectionWrapper className={sectionStyles}>
              <h2>Current address</h2>

              <AddressFields
                name={residentialAddressName}
                labelHint="Required"
                zipCityEnabled
                handleLocationChange={handleCurrentZipCityChange}
              />
            </SectionWrapper>

            <SectionWrapper className={sectionStyles}>
              <h2>Backup address</h2>
              <p>
                Provide a physical address where either you can be reached or someone can contact you while you are in
                transit during your move.
              </p>

              <AddressFields
                name={backupAddressName}
                labelHint="Required"
                zipCityEnabled
                handleLocationChange={handleBackupZipCityChange}
              />
            </SectionWrapper>

            <SectionWrapper className={sectionStyles}>
              <h2>Backup contact</h2>
              <p>
                If we can&apos;t reach you, who can we contact? Any person you assign as a backup contact must be 18
                years of age or older.
              </p>

              <BackupContactInfoFields name={backupContactName} labelHint="Required" />
            </SectionWrapper>

            <div className={formStyles.formActions}>
              <WizardNavigation
                editMode
                disableNext={!isValid || isSubmitting}
                onNextClick={handleSubmit}
                onCancelClick={onCancel}
              />
            </div>
          </Form>
        );
      }}
    </Formik>
  );
};

EditContactInfoForm.propTypes = {
  initialValues: PropTypes.shape({
    telephone: PropTypes.string.isRequired,
    secondary_telephone: PropTypes.string,
    personal_email: PropTypes.string.isRequired,
    phone_is_preferred: PropTypes.bool,
    email_is_preferred: PropTypes.bool,
    [residentialAddressName]: ResidentialAddressShape.isRequired,
    [backupAddressName]: ResidentialAddressShape.isRequired,
    [backupContactName]: BackupContactShape.isRequired,
  }),
  onCancel: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
};

EditContactInfoForm.defaultProps = {
  initialValues: {
    secondaryTelephone: '',
    phoneIsPreferred: false,
    emailIsPreferred: false,
  },
};

export default EditContactInfoForm;
