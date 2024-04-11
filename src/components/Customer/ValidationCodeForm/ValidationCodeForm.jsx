import React from 'react';
import { Formik } from 'formik';
import * as Yup from 'yup';

import { Form } from 'components/form/Form';
import TextField from 'components/form/fields/TextField/TextField';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import formStyles from 'styles/form.module.scss';

const ValidationCodeForm = ({ initialValues, onSubmit }) => {
  const validationSchema = Yup.object().shape({
    code: Yup.string()
      .matches(/[0-9]{20}/, 'Enter a 20-digit number')
      .required('Required'),
  });

  return (
    <Formik initialValues={initialValues} validateOnMount validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ isValid, handleSubmit }) => {
        return (
          <Form className={formStyles.form}>
            <h1>Please input your validation code</h1>
            <TextField label="Validation code" name="code" id="code" required maxLength="20" />

            <div className={formStyles.formActions}>
              <WizardNavigation disableNext={!isValid} onNextClick={handleSubmit} />
            </div>
          </Form>
        );
      }}
    </Formik>
  );
};

export default ValidationCodeForm;
