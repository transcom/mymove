import React from 'react';
import { Formik } from 'formik';
import PropTypes from 'prop-types';

import { CustomerContactInfoFields } from 'components/form/CustomerContactInfoFields';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { Form } from 'components/form/Form';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import formStyles from 'styles/form.module.scss';
import { contactInfoSchema } from 'utils/validation';

const ContactInfoForm = ({ initialValues, onSubmit, onBack }) => {
  return (
    <Formik initialValues={initialValues} onSubmit={onSubmit} validationSchema={contactInfoSchema} validateOnMount>
      {({ isValid, isSubmitting, handleSubmit }) => {
        return (
          <Form className={formStyles.form}>
            <h1>Your contact info</h1>
            <SectionWrapper className={formStyles.formSection}>
              <div className="tablet:margin-top-neg-3">
                <CustomerContactInfoFields labelHint="Required" />
              </div>
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

ContactInfoForm.propTypes = {
  initialValues: PropTypes.shape({
    telephone: PropTypes.string,
    secondary_telephone: PropTypes.string,
    personal_email: PropTypes.string,
    phone_is_preferred: PropTypes.bool,
    email_is_preferred: PropTypes.bool,
  }).isRequired,
  onSubmit: PropTypes.func.isRequired,
  onBack: PropTypes.func.isRequired,
};

export default ContactInfoForm;
