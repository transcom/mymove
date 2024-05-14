import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { Formik } from 'formik';
import * as Yup from 'yup';

import { SERVICE_MEMBER_AGENCY_LABELS } from 'content/serviceMemberAgencies';
import { Form } from 'components/form/Form';
import TextField from 'components/form/fields/TextField/TextField';
import { DropdownInput } from 'components/form/fields/DropdownInput';
import SectionWrapper from 'components/Customer/SectionWrapper';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { dropdownInputOptions } from 'utils/formatters';
import formStyles from 'styles/form.module.scss';

const DodInfoForm = ({ initialValues, onSubmit, onBack, isEmplidEnabled }) => {
  const branchOptions = dropdownInputOptions(SERVICE_MEMBER_AGENCY_LABELS);
  const [showEmplid, setShowEmplid] = useState(initialValues.affiliation === 'COAST_GUARD');

  const validationSchema = Yup.object().shape({
    affiliation: Yup.mixed().oneOf(Object.keys(SERVICE_MEMBER_AGENCY_LABELS)).required('Required'),
    emplid: Yup.string().when('showEmplid', () => {
      if (showEmplid && isEmplidEnabled)
        return Yup.string()
          .matches(/[0-9]{7}/, 'Enter a 7-digit EMPLID number')
          .required('Required');
      return Yup.string().nullable();
    }),
  });

  return (
    <Formik
      initialValues={initialValues}
      validateOnMount
      validationSchema={validationSchema}
      onSubmit={onSubmit}
      showEmplid={showEmplid}
      setShowEmplid={setShowEmplid}
    >
      {({ isValid, isSubmitting, handleSubmit, handleChange }) => {
        const handleBranchChange = (e) => {
          if (e.target.value === 'COAST_GUARD') {
            setShowEmplid(true);
          } else {
            setShowEmplid(false);
          }
        };
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
                onChange={(e) => {
                  handleChange(e);
                  handleBranchChange(e);
                }}
              />
              <TextField
                label="DOD ID number"
                name="edipi"
                id="edipi"
                required
                maxLength="10"
                inputMode="numeric"
                pattern="[0-9]{10}"
                isDisabled
              />
              {showEmplid && isEmplidEnabled && (
                <TextField
                  label="EMPLID"
                  name="emplid"
                  id="emplid"
                  required
                  maxLength="7"
                  inputMode="numeric"
                  pattern="[0-9]{7}"
                />
              )}
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
  }).isRequired,
  onSubmit: PropTypes.func.isRequired,
  onBack: PropTypes.func.isRequired,
};

export default DodInfoForm;
