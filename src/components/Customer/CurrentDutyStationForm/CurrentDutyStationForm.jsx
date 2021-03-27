import React from 'react';
import { func, shape } from 'prop-types';
import { Formik } from 'formik';
import * as Yup from 'yup';

import { Form } from 'components/form/Form';
import { DutyStationInput } from 'components/form/fields/DutyStationInput';
import SectionWrapper from 'components/Customer/SectionWrapper';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import formStyles from 'styles/form.module.scss';
import { DutyStationShape } from 'types/dutyStation';

const CurrentDutyStationForm = ({ initialValues, onBack, onSubmit, newDutyStation }) => {
  const validationSchema = Yup.object().shape({
    current_station: Yup.object()
      .required('Required')
      .test(
        'existing and new duty station should not match',
        'You entered the same duty station for your origin and destination. Please change one of them.',
        (value) => value?.id !== newDutyStation?.id,
      ),
  });

  return (
    <Formik initialValues={initialValues} validateOnMount validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ isValid, handleSubmit, isSubmitting }) => {
        return (
          <Form className={formStyles.form}>
            <h1>Current duty station</h1>
            <SectionWrapper className={formStyles.formSection}>
              <DutyStationInput
                label="What is your current duty station?"
                name="current_station"
                id="current_station"
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

CurrentDutyStationForm.propTypes = {
  initialValues: shape({
    current_station: DutyStationShape,
  }).isRequired,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
  newDutyStation: DutyStationShape,
};

CurrentDutyStationForm.defaultProps = {
  newDutyStation: {},
};

export default CurrentDutyStationForm;
