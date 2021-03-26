import React from 'react';
import { func, shape, string } from 'prop-types';
import { Formik } from 'formik';
import * as Yup from 'yup';

import { Form } from 'components/form/Form';
import TextField from 'components/form/fields/TextField';
import SectionWrapper from 'components/Customer/SectionWrapper';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import formStyles from 'styles/form.module.scss';

const NameForm = ({ initialValues, onSubmit, onBack }) => {
  const validationSchema = Yup.object().shape({
    first_name: Yup.string().required('Required'),
    middle_name: Yup.string(),
    last_name: Yup.string().required('Required'),
    suffix: Yup.string(),
  });

  return (
    <Formik initialValues={initialValues} validateOnMount validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ isValid, isSubmitting, handleSubmit }) => {
        return (
          <Form className={formStyles.form}>
            <h1>Name</h1>
            <SectionWrapper className={formStyles.formSection}>
              <TextField label="First name" name="first_name" id="firstName" required />
              <TextField label="Middle name" name="middle_name" id="middleName" labelHint="Optional" />
              <TextField label="Last name" name="last_name" id="lastName" required />
              <TextField label="Suffix" name="suffix" id="suffix" labelHint="Optional" />
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

NameForm.propTypes = {
  initialValues: shape({
    first_name: string,
    middle_name: string,
    last_name: string,
    suffix: string,
  }).isRequired,
  onSubmit: func.isRequired,
  onBack: func.isRequired,
};

export default NameForm;
