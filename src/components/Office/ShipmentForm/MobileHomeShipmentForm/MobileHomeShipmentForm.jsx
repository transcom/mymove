import React from 'react';
import { Fieldset } from '@trussworks/react-uswds';
import { func } from 'prop-types';
import classnames from 'classnames';

import styles from './MobileHomeShipmentForm.module.scss';

import formStyles from 'styles/form.module.scss';
import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import TextField from 'components/form/fields/TextField/TextField';
import { ErrorMessage } from 'components/form/index';
import RequiredAsterisk, { requiredAsteriskMessage } from 'components/form/RequiredAsterisk';

const MobileHomeShipmentForm = ({
  lengthHasError,
  widthHasError,
  heightHasError,
  setFieldTouched,
  setFieldError,
  validateForm,
  dimensionError,
}) => {
  return (
    <div className={styles.formContainer}>
      <SectionWrapper className={formStyles.formSection}>
        <h2>Mobile Home Information</h2>
        {requiredAsteriskMessage}
        <div className="grid-row grid-gap">
          <div className="mobile-lg:grid-col-3">
            <MaskedTextField
              data-testid="year"
              name="year"
              label="Year"
              id="year"
              aria-label="Year of the mobile home"
              showRequiredAsterisk
              required
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
            />
          </div>
        </div>
        <div className={classnames(styles.formFieldContainer, 'mobile-lg:grid-col-7')}>
          <TextField
            data-testid="make"
            name="make"
            label="Make"
            id="make"
            aria-label="Make of the mobile home"
            showRequiredAsterisk
            required
          />
          <TextField
            data-testid="model"
            name="model"
            label="Model"
            id="model"
            aria-label="Model of the mobile home"
            showRequiredAsterisk
            required
          />
        </div>
      </SectionWrapper>
      <SectionWrapper className={classnames(styles.sectionWrapper, formStyles.formSection, 'origin')}>
        <h2>Mobile Home Dimensions</h2>
        <p>Enter the total outside dimensions of the mobile home.</p>
        {requiredAsteriskMessage}
        <div>
          <ErrorMessage display={dimensionError}>
            <p>
              The dimensions do not meet the requirements for a mobile home shipment. Please cancel and select a
              different shipment type.
            </p>
          </ErrorMessage>
          <Fieldset className={styles.formFieldContainer}>
            <div className="labelWrapper">
              <legend
                className="usa-label"
                data-testid="mobileHomeLength"
                aria-label="Required: What is the length of the mobile home?"
              >
                <span required>
                  Length <RequiredAsterisk />
                </span>
              </legend>
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
              <legend
                className="usa-label"
                data-testid="mobileHomeWidth"
                aria-label="Required: What is the width of the mobile home?"
              >
                <span required>
                  Width <RequiredAsterisk />
                </span>
              </legend>
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
              <legend
                className="usa-label"
                data-testid="mobileHomeHeight"
                aria-label="Required: What is the height of the mobile home?"
              >
                <span required>
                  Height <RequiredAsterisk />
                </span>
              </legend>
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
    </div>
  );
};

export default MobileHomeShipmentForm;

MobileHomeShipmentForm.propTypes = {
  lengthHasError: func.isRequired,
  widthHasError: func.isRequired,
  heightHasError: func.isRequired,
  setFieldTouched: func.isRequired,
  setFieldError: func.isRequired,
  validateForm: func.isRequired,
};
