import React from 'react';
import { func, shape } from 'prop-types';
import { Formik } from 'formik';
import * as Yup from 'yup';

import { Form } from 'components/form/Form';
import { DutyLocationInput } from 'components/form/fields/DutyLocationInput';
import SectionWrapper from 'components/Customer/SectionWrapper';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import formStyles from 'styles/form.module.scss';
import { DutyLocationShape } from 'types/dutyLocation';

const CurrentDutyLocationForm = ({ initialValues, onBack, onSubmit }) => {
  const validationSchema = Yup.object().shape({
    current_location: Yup.object().required('Required'),
  });

  return (
    <Formik initialValues={initialValues} validateOnMount validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ isValid, handleSubmit, isSubmitting }) => {
        return (
          <Form className={formStyles.form}>
            <h1>Current duty location</h1>
            <SectionWrapper className={formStyles.formSection}>
              <DutyLocationInput
                label="What is your current duty location?"
                name="current_location"
                id="current_location"
                required
              />
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

CurrentDutyLocationForm.propTypes = {
  initialValues: shape({
    current_location: DutyLocationShape,
  }).isRequired,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
};

export default CurrentDutyLocationForm;
