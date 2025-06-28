import classnames from 'classnames';
import React from 'react';
import PropTypes from 'prop-types';
import { Formik } from 'formik';
import { Link } from 'react-router-dom';

import oktaLogo from '../../../shared/images/okta_logo.png';

import editOktaInfoFormStyle from './EditOktaInfoForm.module.scss';

import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import { OktaInfoFields } from 'components/form/OktaInfoFields';
import { oktaInfoSchema } from 'utils/validation';

const EditOktaInfoForm = ({ initialValues, onSubmit, onCancel, isMoveLocked }) => {
  const sectionStyles = classnames(formStyles.formSection, editOktaInfoFormStyle.formSection);
  const hostname = window && window.location && window.location.hostname;
  const url =
    hostname === 'my.move.mil'
      ? 'https://milmove.okta.mil/enduser/settings'
      : 'https://test-milmove.okta.mil/enduser/settings';

  return (
    <Formik initialValues={initialValues} onSubmit={onSubmit} validateOnMount validationSchema={oktaInfoSchema}>
      {({ isValid, isSubmitting, handleSubmit }) => {
        return (
          <Form className={classnames(formStyles.form, editOktaInfoFormStyle.form)}>
            <a href={url}>
              <img className={editOktaInfoFormStyle.oktaLogo} src={oktaLogo} alt="Okta logo" />
            </a>

            <SectionWrapper className={sectionStyles}>
              <h2>Your Okta Profile</h2>
              <p>
                This is the information stored in your Okta Profile used for logging into MilMove. If you wish to change
                any information, you can do so here by changing the appropriate fields and clicking <b>Save</b>.
              </p>
              <p>
                If you need to update your security methods or set up additional security methods, you can update/add
                those in the <Link to={url}>Okta Dashboard.</Link>
              </p>

              <OktaInfoFields />
            </SectionWrapper>

            <div className={formStyles.formActions}>
              <WizardNavigation
                editMode
                disableNext={!isValid || isSubmitting}
                onNextClick={handleSubmit}
                onCancelClick={onCancel}
                readOnly={isMoveLocked}
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
    oktaUsername: PropTypes.string.isRequired,
    oktaEmail: PropTypes.string.isRequired,
    oktaFirstName: PropTypes.string.isRequired,
    oktaLastName: PropTypes.string.isRequired,
  }),
  onCancel: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
};

EditOktaInfoForm.defaultProps = {
  initialValues: {
    oktaEdipi: 'Not Provided',
  },
};

export default EditOktaInfoForm;
