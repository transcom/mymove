import React, { useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import { Formik } from 'formik';
import * as Yup from 'yup';
import { Grid } from '@trussworks/react-uswds';

import { SERVICE_MEMBER_AGENCY_LABELS } from 'content/serviceMemberAgencies';
import { Form } from 'components/form/Form';
import TextField from 'components/form/fields/TextField/TextField';
import { DropdownInput } from 'components/form/fields/DropdownInput';
import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { dropdownInputOptions } from 'utils/formatters';
import formStyles from 'styles/form.module.scss';
import { DutyLocationShape } from 'types/dutyLocation';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { requiredAsteriskMessage } from 'components/form/RequiredAsterisk';

const ServiceInfoForm = ({ initialValues, onSubmit, onCancel }) => {
  const branchOptions = dropdownInputOptions(SERVICE_MEMBER_AGENCY_LABELS);
  const [showEmplid, setShowEmplid] = useState(initialValues.affiliation === 'COAST_GUARD');
  const [isDodidDisabled, setIsDodidDisabled] = useState(false);

  useEffect(() => {
    // checking feature flag to see if DODID input should be disabled
    // this data pulls from Okta and doens't let the customer update it
    const fetchData = async () => {
      setIsDodidDisabled(await isBooleanFlagEnabled('okta_dodid_input'));
    };
    fetchData();
  }, []);

  const validationSchema = Yup.object().shape({
    first_name: Yup.string().required('Required'),
    middle_name: Yup.string(),
    last_name: Yup.string().required('Required'),
    suffix: Yup.string(),
    affiliation: Yup.mixed().oneOf(Object.keys(SERVICE_MEMBER_AGENCY_LABELS)).required('Required'),
    edipi: isDodidDisabled
      ? Yup.string().notRequired()
      : Yup.string()
          .matches(/[0-9]{10}/, 'Enter a 10-digit DOD ID number')
          .required('Required'),
    emplid: Yup.string().when('showEmplid', () => {
      if (showEmplid)
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
            <h1>Edit service info</h1>
            <SectionWrapper className={formStyles.formSection}>
              {requiredAsteriskMessage}
              <Grid row gap>
                <Grid mobileLg={{ col: 6 }}>
                  <TextField label="First name" name="first_name" id="firstName" showRequiredAsterisk required />
                </Grid>
                <Grid mobileLg={{ col: 6 }}>
                  <TextField label="Middle name" name="middle_name" id="middleName" />
                </Grid>
              </Grid>

              <Grid row gap>
                <Grid mobileLg={{ col: 6 }}>
                  <TextField label="Last name" name="last_name" id="lastName" showRequiredAsterisk required />
                </Grid>
                <Grid mobileLg={{ col: 6 }}>
                  <TextField label="Suffix" name="suffix" id="suffix" />
                </Grid>
              </Grid>

              <Grid row gap>
                <Grid mobileLg={{ col: 6 }}>
                  <DropdownInput
                    label="Branch of service"
                    name="affiliation"
                    id="affiliation"
                    showRequiredAsterisk
                    required
                    options={branchOptions}
                    onChange={(e) => {
                      handleChange(e);
                      handleBranchChange(e);
                    }}
                  />
                </Grid>
                {showEmplid && (
                  <Grid mobileLg={{ col: 6 }}>
                    <TextField
                      label="EMPLID"
                      name="emplid"
                      id="emplid"
                      showRequiredAsterisk
                      required
                      maxLength="7"
                      inputMode="numeric"
                      pattern="[0-9]{7}"
                    />
                  </Grid>
                )}
              </Grid>

              <Grid row gap>
                <Grid mobileLg={{ col: 6 }}>
                  <TextField
                    label="DoD ID number"
                    name="edipi"
                    id="edipi"
                    showRequiredAsterisk
                    required
                    maxLength="10"
                    inputMode="numeric"
                    pattern="[0-9]{10}"
                    isDisabled={isDodidDisabled}
                  />
                </Grid>
              </Grid>
            </SectionWrapper>

            <div className={formStyles.formActions}>
              <WizardNavigation
                editMode
                onCancelClick={onCancel}
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

ServiceInfoForm.propTypes = {
  initialValues: PropTypes.shape({
    current_location: DutyLocationShape,
  }).isRequired,
  onCancel: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
};

export default ServiceInfoForm;
