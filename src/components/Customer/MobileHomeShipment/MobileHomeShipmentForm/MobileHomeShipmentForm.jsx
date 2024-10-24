import React from 'react';
import { func } from 'prop-types';
import * as Yup from 'yup';
import { Formik, Field } from 'formik';
import { Button, Form, Label, Textarea } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './MobileHomeShipmentForm.module.scss';

import SectionWrapper from 'components/Customer/SectionWrapper';
import Hint from 'components/Hint';
import Fieldset from 'shared/Fieldset';
import formStyles from 'styles/form.module.scss';
import { ShipmentShape } from 'types/shipment';
import TextField from 'components/form/fields/TextField/TextField';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import Callout from 'components/Callout';
import { ErrorMessage } from 'components/form/index';
import { convertInchesToFeetAndInches } from 'utils/formatMtoShipment';
import RequiredTag from 'components/form/RequiredTag';

const currentYear = new Date().getFullYear();
const maxYear = currentYear + 2;

const validationShape = {
  year: Yup.number().required('Required').min(1700, 'Invalid year').max(maxYear, 'Invalid year'),
  make: Yup.string().required('Required'),
  model: Yup.string().required('Required'),
  lengthFeet: Yup.number()
    .min(0)
    .nullable()
    .when('lengthInches', {
      is: (lengthInches) => !lengthInches,
      then: (schema) => schema.required('Required'),
      otherwise: (schema) => schema.notRequired(),
    }),
  lengthInches: Yup.number(),
  widthFeet: Yup.number()
    .min(0)
    .nullable()
    .when('widthInches', {
      is: (widthInches) => !widthInches,
      then: (schema) => schema.required('Required'),
      otherwise: (schema) => schema.notRequired(),
    }),
  widthInches: Yup.number().min(0),
  heightFeet: Yup.number()
    .min(0)
    .nullable()
    .when('heightInches', {
      is: (heightInches) => !heightInches,
      then: (schema) => schema.required('Required'),
      otherwise: (schema) => schema.notRequired(),
    }),
  heightInches: Yup.number().min(0),
  customerRemarks: Yup.string(),
};

const MobileHomeShipmentForm = ({ mtoShipment, onBack, onSubmit }) => {
  const { year, make, model, lengthInInches, widthInInches, heightInInches } = mtoShipment?.mobileHomeShipment || {};

  const length = convertInchesToFeetAndInches(lengthInInches);
  const width = convertInchesToFeetAndInches(widthInInches);
  const height = convertInchesToFeetAndInches(heightInInches);

  const initialValues = {
    year: year?.toString() || null,
    make: make || '',
    model: model || '',
    lengthFeet: length.feet,
    lengthInches: length.inches,
    widthFeet: width.feet,
    widthInches: width.inches,
    heightFeet: height.feet,
    heightInches: height.inches,
    customerRemarks: mtoShipment?.customerRemarks,
  };

  return (
    <Formik
      initialValues={initialValues}
      validationSchema={Yup.object().shape(validationShape)}
      onSubmit={onSubmit}
      validateOnMount
    >
      {({ isValid, handleSubmit, values, errors, touched, setFieldTouched, setFieldError, validateForm }) => {
        const lengthHasError = !!(
          (touched.lengthFeet && errors.lengthFeet) ||
          (touched.lengthInches && errors.lengthFeet)
        );
        const widthHasError = !!((touched.widthFeet && errors.widthFeet) || (touched.widthInches && errors.widthFeet));
        const heightHasError = !!(
          (touched.heightFeet && errors.heightFeet) ||
          (touched.heightInches && errors.heightFeet)
        );
        if (touched.lengthInches && !touched.lengthFeet) {
          setFieldTouched('lengthFeet', true);
        }
        if (touched.widthInches && !touched.widthFeet) {
          setFieldTouched('widthFeet', true);
        }
        if (touched.heightInches && !touched.heightFeet) {
          setFieldTouched('heightFeet', true);
        }
        // manually turn off 'required' error when page loads if field is empty.
        if (values.year === null && !touched.year && errors.year === 'Required') {
          setFieldError('year', null);
        }
        return (
          <div className={styles.formContainer}>
            <Form className={formStyles.form}>
              <SectionWrapper className={classnames(styles.sectionWrapper, formStyles.formSection, 'origin')}>
                <h2>Mobile Home Information</h2>
                <div className="grid-row grid-gap">
                  <div className="mobile-lg:grid-col-3">
                    <MaskedTextField
                      data-testid="year"
                      name="year"
                      label="Year"
                      id="year"
                      maxLength={4}
                      mask={Number}
                      scale={0}
                      signed="false"
                      lazy={false}
                      onChange={() => {
                        setFieldError('year', null);
                      }}
                      onBlur={() => {
                        setFieldTouched('year', true);
                        setFieldError('year', null);
                        validateForm();
                      }}
                      required
                    />
                  </div>
                </div>
                <div className={classnames(styles.formFieldContainer, 'mobile-lg:grid-col-7')}>
                  <TextField data-testid="make" name="make" label="Make" id="make" required />
                  <TextField data-testid="model" name="model" label="Model" id="model" required />
                </div>
              </SectionWrapper>
              <SectionWrapper className={classnames(styles.sectionWrapper, formStyles.formSection, 'origin')}>
                <h2>Mobile Home Dimensions</h2>
                <p>Enter the total outside dimensions (in Feet and Inches) of the Mobile Home.</p>
                <div>
                  <Fieldset className={styles.formFieldContainer}>
                    <div className="labelWrapper">
                      <legend className="usa-label">Length</legend>
                      <RequiredTag />
                      <ErrorMessage display={lengthHasError}>Required</ErrorMessage>
                    </div>
                    <div className={classnames(styles.formTextFieldWrapper, 'grid-row grid-gap')}>
                      <div className="mobile-lg:grid-col-3">
                        <MaskedTextField
                          data-testid="lengthFeet"
                          name="lengthFeet"
                          id="lengthFeet"
                          mask={Number}
                          scale={0} // digits after point, 0 for integers
                          signed="false" // disallow negative
                          lazy={false} // immediate masking evaluation
                          suffix="Feet"
                          errorClassName={styles.hide}
                          title="Length in feet"
                        />
                      </div>
                      <div className="mobile-lg:grid-col-3">
                        <MaskedTextField
                          data-testid="lengthInches"
                          name="lengthInches"
                          id="lengthInches"
                          mask={Number}
                          scale={0} // digits after point, 0 for integers
                          signed="false" // disallow negative
                          lazy={false} // immediate masking evaluation
                          suffix="Inches"
                          max={11}
                          errorClassName={styles.hide}
                          title="Length in inches"
                        />
                      </div>
                    </div>
                  </Fieldset>
                  <Fieldset className={styles.formFieldContainer}>
                    <div className="labelWrapper">
                      <legend className="usa-label">Width</legend>
                      <RequiredTag />
                      <ErrorMessage display={widthHasError}>Required</ErrorMessage>
                    </div>
                    <div className={classnames(styles.formTextFieldWrapper, 'grid-row grid-gap')}>
                      <div className="mobile-lg:grid-col-3">
                        <MaskedTextField
                          data-testid="widthFeet"
                          name="widthFeet"
                          id="widthFeet"
                          mask={Number}
                          scale={0} // digits after point, 0 for integers
                          signed="false" // disallow negative
                          lazy={false} // immediate masking evaluation
                          suffix="Feet"
                          errorClassName={styles.hide}
                          title="Width in feet"
                        />
                      </div>
                      <div className="mobile-lg:grid-col-3">
                        <MaskedTextField
                          data-testid="widthInches"
                          name="widthInches"
                          id="widthInches"
                          mask={Number}
                          scale={0} // digits after point, 0 for integers
                          signed="false" // disallow negative
                          lazy={false} // immediate masking evaluation
                          suffix="Inches"
                          max={11}
                          errorClassName={styles.hide}
                          title="Width in inches"
                        />
                      </div>
                    </div>
                  </Fieldset>
                  <Fieldset className={styles.formFieldContainer}>
                    <div className="labelWrapper">
                      <legend className="usa-label">Height</legend>
                      <RequiredTag />
                      <ErrorMessage display={heightHasError}>Required</ErrorMessage>
                    </div>
                    <div className={classnames(styles.formTextFieldWrapper, 'grid-row grid-gap')}>
                      <div className="mobile-lg:grid-col-3">
                        <MaskedTextField
                          data-testid="heightFeet"
                          name="heightFeet"
                          id="heightFeet"
                          mask={Number}
                          scale={0} // digits after point, 0 for integers
                          signed="false" // disallow negative
                          lazy={false} // immediate masking evaluation
                          suffix="Feet"
                          errorClassName={styles.hide}
                          title="Height in feet"
                        />
                      </div>
                      <div className="mobile-lg:grid-col-3">
                        <MaskedTextField
                          data-testid="heightInches"
                          name="heightInches"
                          id="heightInches"
                          mask={Number}
                          scale={0} // digits after point, 0 for integers
                          signed="false" // disallow negative
                          lazy={false} // immediate masking evaluation
                          suffix="Inches"
                          max={11}
                          errorClassName={styles.hide}
                          title="Height in inches"
                        />
                      </div>
                    </div>
                  </Fieldset>
                </div>
              </SectionWrapper>
              <SectionWrapper className={formStyles.formSection}>
                <Fieldset legend={<div className={formStyles.legendContent}>Remarks</div>}>
                  <Label htmlFor="customerRemarks">
                    Are there things about this mobile home shipment that your counselor or movers should know or
                    discuss with you?
                  </Label>

                  <Callout>
                    Example
                    <ul>
                      <li>
                        Is there additional information you feel is pertinent to the processing of your mobile home
                        shipment?(e.g., &lsquo;wrecker service requested&rsquo; and &lsquo;crane service needed&rsquo;).
                      </li>
                    </ul>
                  </Callout>

                  <Field
                    as={Textarea}
                    data-testid="remarks"
                    name="customerRemarks"
                    className={`${formStyles.remarks}`}
                    placeholder="Do not itemize your personal property here. Your movers will help do that when they talk to you."
                    id="customerRemarks"
                    maxLength={250}
                  />
                  <Hint>
                    <p>250 characters</p>
                  </Hint>
                </Fieldset>
              </SectionWrapper>
              <div className={styles.buttonContainer}>
                <Button className={styles.backButton} type="button" onClick={onBack} secondary outline>
                  Back
                </Button>
                <Button className={styles.saveButton} type="button" onClick={handleSubmit} disabled={!isValid}>
                  Continue
                </Button>
              </div>
            </Form>
          </div>
        );
      }}
    </Formik>
  );
};

MobileHomeShipmentForm.propTypes = {
  mtoShipment: ShipmentShape,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
};

MobileHomeShipmentForm.defaultProps = {
  mtoShipment: undefined,
};

export default MobileHomeShipmentForm;
