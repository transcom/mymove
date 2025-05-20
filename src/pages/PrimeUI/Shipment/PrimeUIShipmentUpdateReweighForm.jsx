import React from 'react';
import { func, shape, string } from 'prop-types';
import * as Yup from 'yup';
import { Formik } from 'formik';
import { Textarea, Label } from '@trussworks/react-uswds';

import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';

const validationSchema = Yup.object({
  reweighWeight: Yup.number().min(1, 'Authorized weight must be greater than or equal to 1').required('Required'),
  reweighRemarks: Yup.string().required('Required'),
});

const ReweighForm = ({ onSubmit, handleClose, initialValues }) => {
  return (
    <Formik initialValues={initialValues} onSubmit={onSubmit} validationSchema={validationSchema}>
      {({ isValid, isSubmitting, handleSubmit, handleChange, setTouched, values }) => (
        <Form className={formStyles.form}>
          <SectionWrapper className={formStyles.formSection}>
            <MaskedTextField
              defaultValue="0"
              inputTestId="reweighWeightInput"
              id="reweighWeight"
              lazy={false} // immediate masking evaluation
              label="Reweigh Weight (lbs)"
              mask={Number}
              name="reweighWeight"
              scale={0} // digits after point, 0 for integers
              signed={false} // disallow negative
              thousandsSeparator=","
            />
            <Label htmlFor="remarks">Remarks</Label>
            <Textarea
              data-testid="remarks"
              id="reweighRemarks"
              maxLength={500}
              placeholder=""
              onChange={handleChange}
              onBlur={() => setTouched({ reweighRemarks: true }, false)}
              value={values.reweighRemarks}
            />
          </SectionWrapper>

          <WizardNavigation
            editMode
            className={formStyles.formActions}
            aria-label="Update Reweigh"
            type="submit"
            disableNext={!isValid || isSubmitting}
            onCancelClick={handleClose}
            onNextClick={handleSubmit}
          />
        </Form>
      )}
    </Formik>
  );
};

// disableNext={isSubmitting || !isValid}
ReweighForm.propTypes = {
  onSubmit: func.isRequired,
  handleClose: func.isRequired,
  initialValues: shape({
    reweighWeight: string,
    reweighRemarks: string,
  }).isRequired,
};

export default ReweighForm;
