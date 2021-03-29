import React from 'react';
import { Formik } from 'formik';
import * as Yup from 'yup';
import PropTypes from 'prop-types';

import { BackupContactInfoFields } from 'components/form/BackupContactInfoFields';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { Form } from 'components/form/Form';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import formStyles from 'styles/form.module.scss';

const BackupContactForm = ({ initialValues, onSubmit, onBack }) => {
  const validationSchema = Yup.object().shape({
    name: Yup.string().required('Required'),
    email: Yup.string()
      .matches(/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+.[a-zA-Z]{2,}$/, 'Must be a valid email address')
      .required('Required'),
    telephone: Yup.string().min(12, 'Number must have 10 digits and a valid area code').required('Required'), // min 12 includes hyphens
  });

  return (
    <Formik initialValues={initialValues} onSubmit={onSubmit} validationSchema={validationSchema} validateOnMount>
      {({ isValid, isSubmitting, handleSubmit }) => {
        return (
          <Form className={formStyles.form}>
            <SectionWrapper className={formStyles.formSection}>
              <div className="tablet:margin-top-neg-3">
                <BackupContactInfoFields
                  legend="Backup contact"
                  render={(fields) => (
                    <>
                      <p>
                        If we can&apos;t reach you, who can we contact? Any person you assign as a backup contact must
                        be 18 years of age or older.
                      </p>
                      {fields}
                    </>
                  )}
                />
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
