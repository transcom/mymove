import React from 'react';
import { Formik } from 'formik';
import PropTypes from 'prop-types';

import { BackupContactInfoFields } from 'components/form/BackupContactInfoFields';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { Form } from 'components/form/Form';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { backupContactInfoSchema } from 'utils/validation';
import formStyles from 'styles/form.module.scss';

const BackupContactForm = ({ initialValues, onSubmit, onBack }) => {
  return (
    <Formik
      initialValues={initialValues}
      onSubmit={onSubmit}
      validationSchema={backupContactInfoSchema}
      validateOnMount
    >
      {({ isValid, isSubmitting, handleSubmit }) => {
        return (
          <Form className={formStyles.form}>
            <h1>Backup contact</h1>
            <p>
              If we cannot reach you, who can we contact (such as spouse or parent)? Any person you assign as a backup
              contact must be 18 years of age or older.
            </p>
            <SectionWrapper className={formStyles.formSection}>
              <div className="tablet:margin-top-neg-3">
                <BackupContactInfoFields labelHint="Required" />
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

BackupContactForm.propTypes = {
  initialValues: PropTypes.shape({
    name: PropTypes.string,
    telephone: PropTypes.string,
    email: PropTypes.string,
  }).isRequired,
  onSubmit: PropTypes.func.isRequired,
  onBack: PropTypes.func.isRequired,
};

export default BackupContactForm;
