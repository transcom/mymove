import React from 'react';
import { Fieldset, Radio } from '@trussworks/react-uswds';
import { Field } from 'formik';
import { func } from 'prop-types';
import classnames from 'classnames';

import styles from './BoatShipmentForm.module.scss';

import formStyles from 'styles/form.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import TextField from 'components/form/fields/TextField/TextField';
import { ErrorMessage } from 'components/form/index';
import { boatShipmentTypes } from 'constants/shipments';

const BoatShipmentForm = ({
  lengthHasError,
  widthHasError,
  heightHasError,
  values,
  setFieldTouched,
  setFieldError,
  validateForm,
  dimensionError,
}) => {
  return (
    <div className={styles.formContainer}>
      <SectionWrapper className={formStyles.formSection}>
        <h2>Boat Information</h2>
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
        <h2>Boat Dimensions</h2>
        <p>
          Enter the total outside dimensions of the boat and the trailer (if a trailer is included) with the boat
          sitting on the trailer. If there is no trailer, then enter the outside dimensions of the boat itself.
        </p>
        <div>
          <ErrorMessage display={dimensionError}>
            <p>
              The dimensions do not meet the requirements for a boat shipment. Please cancel and select a different
              shipment type.
            </p>
          </ErrorMessage>
          <Fieldset className={styles.formFieldContainer}>
            <div className="labelWrapper">
              <legend className="usa-label">Length</legend>
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
                  onChange={() => {
                    setFieldError('heightFeet', null);
                    setFieldError('widthFeet', null);
                    setFieldError('lengthFeet', null);
                  }}
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
                  onChange={() => {
                    setFieldError('heightFeet', null);
                    setFieldError('widthFeet', null);
                    setFieldError('lengthFeet', null);
                  }}
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
                  onChange={() => {
                    setFieldError('heightFeet', null);
                    setFieldError('widthFeet', null);
                    setFieldError('lengthFeet', null);
                  }}
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
        <h2>Trailer Status</h2>
        <Fieldset>
          <legend className="usa-label">Does the boat have a trailer?</legend>
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
        <h2>Shipment Method</h2>
        <Fieldset>
          <legend className="usa-label">What is the method of shipment?</legend>
          <Field
            as={Radio}
            id="typeTowAway"
            data-testid="typeTowAway"
            label="Boat Tow-Away (BTA)"
            name="type"
            value={boatShipmentTypes.TOW_AWAY}
            checked={values.type === boatShipmentTypes.TOW_AWAY}
          />
          <Field
            as={Radio}
            id="typeHaulAway"
            data-testid="typeHaulAway"
            label="Boat Haul-Away (BHA)"
            name="type"
            value={boatShipmentTypes.HAUL_AWAY}
            checked={values.type === boatShipmentTypes.HAUL_AWAY}
          />
        </Fieldset>
      </SectionWrapper>
    </div>
  );
};

export default BoatShipmentForm;

BoatShipmentForm.propTypes = {
  lengthHasError: func.isRequired,
  widthHasError: func.isRequired,
  heightHasError: func.isRequired,
  setFieldTouched: func.isRequired,
  setFieldError: func.isRequired,
  validateForm: func.isRequired,
};
