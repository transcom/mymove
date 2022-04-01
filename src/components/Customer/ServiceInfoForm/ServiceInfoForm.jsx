import React from 'react';
import PropTypes from 'prop-types';
import { Formik } from 'formik';
import * as Yup from 'yup';
import { Grid } from '@trussworks/react-uswds';

import { ORDERS_RANK_OPTIONS } from 'constants/orders';
import { SERVICE_MEMBER_AGENCY_LABELS } from 'content/serviceMemberAgencies';
import { Form } from 'components/form/Form';
import TextField from 'components/form/fields/TextField/TextField';
import { DropdownInput } from 'components/form/fields/DropdownInput';
import { DutyLocationInput } from 'components/form/fields/DutyLocationInput';
import SectionWrapper from 'components/Customer/SectionWrapper';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { dropdownInputOptions } from 'shared/formatters';
import formStyles from 'styles/form.module.scss';
import { DutyLocationShape } from 'types/dutyLocation';

const ServiceInfoForm = ({ initialValues, onSubmit, onCancel, newDutyLocation }) => {
  const branchOptions = dropdownInputOptions(SERVICE_MEMBER_AGENCY_LABELS);
  const rankOptions = dropdownInputOptions(ORDERS_RANK_OPTIONS);

  const validationSchema = Yup.object().shape({
    first_name: Yup.string().required('Required'),
    middle_name: Yup.string(),
    last_name: Yup.string().required('Required'),
    suffix: Yup.string(),
    affiliation: Yup.mixed().oneOf(Object.keys(SERVICE_MEMBER_AGENCY_LABELS)).required('Required'),
    edipi: Yup.string()
      .matches(/[0-9]{10}/, 'Enter a 10-digit DOD ID number')
      .required('Required'),
    rank: Yup.mixed().oneOf(Object.keys(ORDERS_RANK_OPTIONS)).required('Required'),
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
      {({ isValid, isSubmitting, handleSubmit }) => {
        return (
          <Form className={formStyles.form}>
            <h1>Edit service info</h1>
            <SectionWrapper className={formStyles.formSection}>
              <Grid row gap>
                <Grid mobileLg={{ col: 6 }}>
                  <TextField label="First name" name="first_name" id="firstName" required />
                </Grid>
                <Grid mobileLg={{ col: 6 }}>
                  <TextField label="Middle name" name="middle_name" id="middleName" labelHint="Optional" />
                </Grid>
              </Grid>

              <Grid row gap>
                <Grid mobileLg={{ col: 6 }}>
                  <TextField label="Last name" name="last_name" id="lastName" required />
                </Grid>
                <Grid mobileLg={{ col: 6 }}>
                  <TextField label="Suffix" name="suffix" id="suffix" labelHint="Optional" />
                </Grid>
              </Grid>

              <Grid row gap>
                <Grid mobileLg={{ col: 6 }}>
                  <DropdownInput
                    label="Branch of service"
                    name="affiliation"
                    id="affiliation"
                    required
                    options={branchOptions}
                  />
                </Grid>
              </Grid>

              <Grid row gap>
                <Grid mobileLg={{ col: 6 }}>
                  <DropdownInput label="Rank" name="rank" id="rank" required options={rankOptions} />
                </Grid>
              </Grid>

              <Grid row gap>
                <Grid mobileLg={{ col: 6 }}>
                  <TextField
                    label="DoD ID number"
                    name="edipi"
                    id="edipi"
                    required
                    maxLength="10"
                    inputMode="numeric"
                    pattern="[0-9]{10}"
                  />
                </Grid>
              </Grid>

              <DutyLocationInput label="Current duty location" name="current_location" id="current_location" required />
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
  newDutyLocation: DutyLocationShape,
};

ServiceInfoForm.defaultProps = {
  newDutyLocation: {},
};

export default ServiceInfoForm;
