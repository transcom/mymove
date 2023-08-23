import classnames from 'classnames';
import React from 'react';
import PropTypes from 'prop-types';
import { Formik } from 'formik';
import * as Yup from 'yup';
import { Link } from 'react-router-dom';

import oktaLogo from '../../../shared/images/okta_logo.png';

import editOktaInfoFormStyle from './EditOktaInfoForm.module.scss';

import SectionWrapper from 'components/Customer/SectionWrapper';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import { backupContactInfoSchema, contactInfoSchema, requiredAddressSchema } from 'utils/validation';
import { ResidentialAddressShape } from 'types/address';
import { OktaProfileFields } from 'components/form/OktaProfileFields';
import { BackupContactShape } from 'types/customerShapes';

export const residentialAddressName = 'residential_address';
export const backupAddressName = 'backup_mailing_address';
export const backupContactName = 'backup_contact';

const EditOktaInfoForm = ({ initialValues, onSubmit, onCancel }) => {
  const validationSchema = Yup.object().shape({
    ...contactInfoSchema.fields,
    [residentialAddressName]: requiredAddressSchema.required(),
    [backupAddressName]: requiredAddressSchema.required(),
    [backupContactName]: backupContactInfoSchema.required(),
  });

  const sectionStyles = classnames(formStyles.formSection, editOktaInfoFormStyle.formSection);
  const url = 'https://test-milmove.okta.mil';

  return (
    <Formik initialValues={initialValues} onSubmit={onSubmit} validateOnMount validationSchema={validationSchema}>
      {({ isValid, isSubmitting, handleSubmit }) => {
        return (
          <Form className={classnames(formStyles.form, editOktaInfoFormStyle.form)}>
            <a href={url}>
              <img className={editOktaInfoFormStyle.oktaLogo} src={oktaLogo} alt="Okta logo" />
            </a>

            <SectionWrapper className={sectionStyles}>
              <h2>Your Okta Profile</h2>
              <p>
                This is the info stored in your Okta Profile used for logging into MilMove. If you wish to change any
                information, you can do so here.
              </p>
              <p>
                If you need to update your security methods or set up additional security methods, you can update/add
                those in the <Link to="https://test-milmove.okta.mil/enduser/settings">Okta Dashboard.</Link>
              </p>

              <OktaProfileFields />
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

EditOktaInfoForm.propTypes = {
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

EditOktaInfoForm.defaultProps = {
  initialValues: {
    secondaryTelephone: '',
    phoneIsPreferred: false,
    emailIsPreferred: false,
  },
};

export default EditOktaInfoForm;
