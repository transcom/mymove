import React from 'react';
import { func } from 'prop-types';
import * as Yup from 'yup';
import { Formik, Field } from 'formik';
import { Button, Form, Radio, Label, Textarea } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './BoatShipmentForm.module.scss';

import RequiredTag from 'components/form/RequiredTag';
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
  hasTrailer: Yup.boolean().required('Required'),
  isRoadworthy: Yup.boolean().when('hasTrailer', {
    is: true,
    then: (schema) => schema.required('Required'),
  }),
  customerRemarks: Yup.string(),
};

const BoatShipmentForm = ({ mtoShipment, onBack, onSubmit }) => {
  const { year, make, model, lengthInInches, widthInInches, heightInInches, hasTrailer, isRoadworthy } =
    mtoShipment?.boatShipment || {};

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
    hasTrailer: hasTrailer === null ? '' : hasTrailer?.toString(),
    isRoadworthy: isRoadworthy === null ? '' : isRoadworthy?.toString(),
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
                <h2>Boat Information</h2>
                <div className="grid-row grid-gap">
                  <div className="mobile-lg:grid-col-3">
                    <MaskedTextField
                      data-testid="year"
                      name="year"
                      label="Year"
                      id="year"
                      labelHint="Required"
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
                  <TextField data-testid="make" name="make" label="Make" id="make" required labelHint="Required" />
                  <TextField data-testid="model" name="model" label="Model" id="model" required labelHint="Required" />
                </div>
              </SectionWrapper>
              <SectionWrapper className={classnames(styles.sectionWrapper, formStyles.formSection, 'origin')}>
                <h2>Boat Dimensions</h2>
                <p>
                  Enter the total outside dimensions of the boat and the trailer (if a trailer is included) with the
                  boat sitting on the trailer. If there is no trailer, then enter the outside dimensions of the boat
                  itself.
                </p>
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
              <SectionWrapper className={classnames(styles.sectionWrapper, formStyles.formSection)}>
                <h2>Trailer Status</h2>
                <Fieldset>
                  <legend className="usa-label">Does the boat have a trailer?</legend>
                  <RequiredTag />
                  <Field
                    as={Radio}
                    id="hasTrailerYes"
                    data-testid="hasTrailerYes"
                    label="Yes"
                    name="hasTrailer"
                    value="true"
                    checked={values.hasTrailer === 'true'}
                  />
                  <Field
                    as={Radio}
                    id="hasTrailerNo"
                    data-testid="hasTrailerNo"
                    label="No"
                    name="hasTrailer"
                    value="false"
                    checked={values.hasTrailer === 'false'}
                  />
                  {values.hasTrailer === 'true' && (
                    <Fieldset className={styles.formFieldContainer}>
                      <legend className="usa-label">Is the trailer roadworthy?</legend>
                      <RequiredTag />
                      <Field
                        as={Radio}
                        id="isRoadworthyYes"
                        data-testid="isRoadworthyYes"
                        label="Yes"
                        name="isRoadworthy"
                        value="true"
                        checked={values.isRoadworthy === 'true'}
                      />
                      <Field
                        as={Radio}
                        id="isRoadworthyNo"
                        data-testid="isRoadworthyNo"
                        label="No"
                        name="isRoadworthy"
                        value="false"
                        checked={values.isRoadworthy === 'false'}
                      />
                    </Fieldset>
                  )}
                </Fieldset>
              </SectionWrapper>
              <SectionWrapper className={formStyles.formSection}>
                <Fieldset legend={<div className={formStyles.legendContent}>Remarks</div>}>
                  <Label htmlFor="customerRemarks">
                    Are there things about this boat shipment that your counselor or movers should know or discuss with
                    you?
                  </Label>

                  <Callout>
                    Examples
                    <ul>
                      <li>
                        Dimensions of the boat on the trailer are significantly different than one would expect given
                        their individual dimensions
                      </li>

                      <li>Access info for your origin or destination address/marina</li>
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

BoatShipmentForm.propTypes = {
  mtoShipment: ShipmentShape,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
};

BoatShipmentForm.defaultProps = {
  mtoShipment: undefined,
};

export default BoatShipmentForm;
