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

const CurrentDutyLocationForm = ({ initialValues, onBack, onSubmit, newDutyLocation }) => {
  const validationSchema = Yup.object().shape({
    current_location: Yup.object()
      .required('Required')
      .test(
        'existing and new duty location should not match',
        'You entered the same duty location for your origin and destination. Please change one of them.',
        (value) => value?.id !== newDutyLocation?.id,
      ),
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
  newDutyLocation: DutyLocationShape,
};

CurrentDutyLocationForm.defaultProps = {
  newDutyLocation: {},
};

export default CurrentDutyLocationForm;
