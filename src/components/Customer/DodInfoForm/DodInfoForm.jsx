import React from 'react';
import PropTypes from 'prop-types';
import { Formik } from 'formik';
import * as Yup from 'yup';

import { ORDERS_RANK_OPTIONS } from 'constants/orders';
import { SERVICE_MEMBER_AGENCY_LABELS } from 'content/serviceMemberAgencies';
import { Form } from 'components/form/Form';
import TextField from 'components/form/fields/TextField';
import { DropdownInput } from 'components/form/fields/DropdownInput';
import SectionWrapper from 'components/Customer/SectionWrapper';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { dropdownInputOptions } from 'shared/formatters';
import formStyles from 'styles/form.module.scss';

const DodInfoForm = ({ initialValues, onSubmit, onBack }) => {
  const branchOptions = dropdownInputOptions(SERVICE_MEMBER_AGENCY_LABELS);
  const rankOptions = dropdownInputOptions(ORDERS_RANK_OPTIONS);

  const validationSchema = Yup.object().shape({
    affiliation: Yup.mixed().oneOf(Object.keys(SERVICE_MEMBER_AGENCY_LABELS)).required('Required'),
    edipi: Yup.string()
      .matches(/[0-9]{10}/, 'Enter a 10-digit DOD ID number')
      .required('Required'),
    rank: Yup.mixed().oneOf(Object.keys(ORDERS_RANK_OPTIONS)).required('Required'),
  });

  return (
    <Formik initialValues={initialValues} validateOnMount validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ isValid, isSubmitting, handleSubmit }) => {
        return (
          <Form className={formStyles.form}>
            <h1>Create your profile</h1>
            <p>Before we can schedule your move, we need to know a little more about you.</p>
            <SectionWrapper className={formStyles.formSection}>
              <DropdownInput
                label="Branch of service"
                name="affiliation"
                id="affiliation"
                required
                options={branchOptions}
              />
              <TextField
                label="DOD ID number"
                name="edipi"
                id="edipi"
                required
                maxLength="10"
                inputMode="numeric"
                pattern="[0-9]{10}"
              />
              <DropdownInput label="Rank" name="rank" id="rank" required options={rankOptions} />
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

DodInfoForm.propTypes = {
  initialValues: PropTypes.shape({
    affiliation: PropTypes.string,
    edipi: PropTypes.string,
    rank: PropTypes.string,
  }).isRequired,
  onSubmit: PropTypes.func.isRequired,
  onBack: PropTypes.func.isRequired,
};

export default DodInfoForm;
