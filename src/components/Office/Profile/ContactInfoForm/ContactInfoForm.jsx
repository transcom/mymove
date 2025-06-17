import React from 'react';
import PropTypes from 'prop-types';
import { Formik } from 'formik';
import * as Yup from 'yup';
import { Fieldset, Button } from '@trussworks/react-uswds';

import TextField from 'components/form/fields/TextField/TextField';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import { phoneSchema } from 'utils/validation';
import { requiredAsteriskMessage } from 'components/form/RequiredAsterisk';

const ContactInfoForm = ({ initialValues, onSubmit, onCancel }) => {
  const validationSchema = Yup.object().shape({
    telephone: phoneSchema.required('Required'),
  });

  return (
    <Formik initialValues={initialValues} onSubmit={onSubmit} validationSchema={validationSchema} validateOnMount>
      {({ isValid, handleSubmit, isSubmitting }) => {
        return (
          <Form className={formStyles.form}>
            <SectionWrapper className={formStyles.formSection}>
              <h2>Your contact info</h2>
              {requiredAsteriskMessage}
              <Fieldset>
                <div className="grid-row grid-gap">
                  <div className="grid-col-6">
                    <TextField label="First name" name="firstName" id="firstName" disabled />
                  </div>
                  <div className="grid-col-6">
                    <TextField label="Middle name" name="middleName" id="middleName" disabled />
                  </div>
                  <div className="grid-col-6">
                    <TextField label="Last name" name="lastName" id="lastName" disabled />
                  </div>
                  <div className="mobile-lg:grid-col-6">
                    <TextField label="Email" id="email" name="email" disabled />
                  </div>
                </div>

                <div className="grid-row grid-gap">
                  <div className="mobile-lg:grid-col-6">
                    <MaskedTextField
                      label="Phone"
                      id="telephone"
                      name="telephone"
                      type="tel"
                      minimum="12"
                      mask="000{-}000{-}0000"
                      required
                      showRequiredAsterisk
                    />
                  </div>
                </div>
              </Fieldset>
            </SectionWrapper>
            <div className={formStyles.formActions}>
              <Button type="button" secondary onClick={onCancel}>
                Cancel
              </Button>
              <Button disabled={isSubmitting || !isValid} type="submit" onClick={handleSubmit}>
                Save
              </Button>
            </div>
          </Form>
        );
      }}
    </Formik>
  );
};

ContactInfoForm.propTypes = {
  initialValues: PropTypes.shape({
    telephone: PropTypes.string.isRequired,
  }),
  onCancel: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
};

ContactInfoForm.defaultProps = {
  initialValues: {},
};

export default ContactInfoForm;
