import React from 'react';
import { func } from 'prop-types';
import * as Yup from 'yup';
import { Formik } from 'formik';
import { Textarea, Label } from '@trussworks/react-uswds';

import SectionWrapper from 'components/Customer/SectionWrapper';
import formStyles from 'styles/form.module.scss';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import MaskedTextField from 'components/form/fields/MaskedTextField';

const validationSchema = Yup.object({
  reweighWeight: Yup.number().min(1, 'Authorized weight must be greater than or equal to 1').required('Required'),
  reweighRemarks: Yup.string().required('Required'),
});

const ReweighForm = ({ handleSubmit, handleClose }) => {
  const initialValues = {
    reweighWeight: '',
    reweighRemarks: '',
  };

  return (
    <Formik enableReinitialize initialValues={initialValues} validationSchema={validationSchema}>
      {({ handleChange, values, setTouched }) => (
        <>
          <SectionWrapper className={formStyles.formSection}>
            <MaskedTextField
              defaultValue="0"
              inputTestId="textInput"
              id="reweighWeight"
              lazy={false} // immediate masking evaluation
              label="Reweigh Weight"
              mask={Number}
              name="reweighWeight"
              scale={0} // digits after point, 0 for integers
              signed={false} // disallow negative
              thousandsSeparator=","
            >
              {' '}
              lbs
            </MaskedTextField>
            <Label htmlFor="remarks">Remarks</Label>
            <Textarea
              data-testid="remarks"
              id="reweighRemarks"
              maxLength={500}
              onChange={handleChange}
              placeholder=""
              onBlur={() => setTouched({ reweighRemarks: true }, false)}
              value={values.reweighRemarks}
            />
          </SectionWrapper>

          <WizardNavigation
            editMode
            className={formStyles.formActions}
            aria-label="Update Reweigh"
            type="submit"
            onCancelClick={handleClose}
            onNextClick={handleSubmit}
          />
        </>
      )}
    </Formik>
  );
};

// disableNext={isSubmitting || !isValid}
ReweighForm.propTypes = {
  handleSubmit: func.isRequired,
  handleClose: func.isRequired,
};

export default ReweighForm;
