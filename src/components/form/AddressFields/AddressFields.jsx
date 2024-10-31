import React, { useRef } from 'react';
import PropTypes from 'prop-types';
import { v4 as uuidv4 } from 'uuid';
import { Fieldset } from '@trussworks/react-uswds';

import { statesList } from '../../../constants/states';

import Hint from 'components/Hint/index';
import styles from 'components/form/AddressFields/AddressFields.module.scss';
import { technicalHelpDeskURL } from 'shared/constants';
import TextField from 'components/form/fields/TextField/TextField';
import { DropdownInput } from 'components/form/fields/DropdownInput';
import LocationInput from 'components/form/fields/LocationInput';

/**
 * @param legend
 * @param className
 * @param name
 * @param render
 * @param validators
 * @param zipCity
 * @param formikFunctionsToValidatePostalCodeOnChange If you are intending to validate the postal code on change, you
 * will need to pass the handleChange and setFieldTouched Formik functions through in an object here.
 * See ResidentialAddressForm for an example.
 * @return {JSX.Element}
 * @constructor
 */
export const AddressFields = ({
  legend,
  className,
  name,
  render,
  validators,
  zipCityEnabled,
  zipCityError,
  handleLocationChange,
  formikFunctionsToValidatePostalCodeOnChange,
  labelHint: labelHintProp,
}) => {
  const addressFieldsUUID = useRef(uuidv4());
  const infoStr = 'If you encounter any inaccurate lookup information please contact the ';
  const errorStr = 'Not all data was able to populate successfully. Contact the ';
  const assistanceStr = ' for further assistance.';

  const postalCodeField = formikFunctionsToValidatePostalCodeOnChange ? (
    <TextField
      label="ZIP"
      id={`zip_${addressFieldsUUID.current}`}
      name={`${name}.postalCode`}
      data-testid={`${name}.postalCode`}
      maxLength={10}
      labelHint={labelHintProp}
      validate={validators?.postalCode}
      isDisabled={zipCityEnabled}
      onChange={async (e) => {
        // If we are validating on change we need to also set the field to touched when it is changed.
        // Formik, by default, only sets the field to touched on blur.
        // The validation errors will not show unless the field has been touched. We await the handleChange event,
        // then we set the field to touched.
        // We send true for the shouldValidate arg to validate the field at the same time.
        await formikFunctionsToValidatePostalCodeOnChange.handleChange(e);
        formikFunctionsToValidatePostalCodeOnChange.setFieldTouched(`${name}.postalCode`, true, true);
      }}
    />
  ) : (
    <TextField
      label="ZIP"
      id={`zip_${addressFieldsUUID.current}`}
      name={`${name}.postalCode`}
      data-testid={`${name}.postalCode`}
      maxLength={10}
      labelHint={labelHintProp}
      validate={validators?.postalCode}
      isDisabled={zipCityEnabled}
    />
  );

  const stateField = zipCityEnabled ? (
    <TextField
      name={`${name}.state`}
      data-testid={`${name}.state`}
      id={`state_${addressFieldsUUID.current}`}
      label="State"
      labelHint={labelHintProp}
      validate={validators?.state}
      isDisabled={zipCityEnabled}
    />
  ) : (
    <DropdownInput
      name={`${name}.state`}
      data-testid={`${name}.state`}
      id={`state_${addressFieldsUUID.current}`}
      label="State"
      labelHint={labelHintProp}
      options={statesList}
      validate={validators?.state}
    />
  );

  return (
    <Fieldset legend={legend} className={className}>
      {render(
        <>
          <TextField
            label="Address 1"
            id={`mailingAddress1_${addressFieldsUUID.current}`}
            name={`${name}.streetAddress1`}
            labelHint={labelHintProp}
            data-testid={`${name}.streetAddress1`}
            validate={validators?.streetAddress1}
          />
          <TextField
            label="Address 2"
            labelHint={labelHintProp ? null : 'Optional'}
            id={`mailingAddress2_${addressFieldsUUID.current}`}
            name={`${name}.streetAddress2`}
            data-testid={`${name}.streetAddress2`}
            validate={validators?.streetAddress2}
          />
          <TextField
            label="Address 3"
            labelHint={labelHintProp ? null : 'Optional'}
            id={`mailingAddress3_${addressFieldsUUID.current}`}
            name={`${name}.streetAddress3`}
            data-testid={`${name}.streetAddress3`}
            validate={validators?.streetAddress3}
          />
          {handleLocationChange && (
            <>
              <LocationInput
                name={`${name}-zipCity`}
                placeholder="Start typing a Zip or City, State Zip"
                label="Location Lookup"
                handleLocationChange={handleLocationChange}
              />
              {!zipCityError && (
                <Hint className={styles.hint} id="locationInfo" data-testid="locationInfo">
                  {infoStr}
                  <a href={technicalHelpDeskURL} target="_blank" rel="noreferrer">
                    Technical Help Desk
                  </a>
                  {assistanceStr}
                </Hint>
              )}
              {zipCityError && (
                <Hint className={styles.hintError} id="zipCityError" data-testid="zipCityError">
                  {errorStr}
                  <a href={technicalHelpDeskURL} target="_blank" rel="noreferrer">
                    Technical Help Desk
                  </a>
                  {assistanceStr}
                </Hint>
              )}
            </>
          )}
          <div className="grid-row grid-gap">
            <div className="mobile-lg:grid-col-6">
              <TextField
                label="City"
                id={`city_${addressFieldsUUID.current}`}
                name={`${name}.city`}
                labelHint={labelHintProp}
                data-testid={`${name}.city`}
                validate={validators?.city}
                isDisabled={zipCityEnabled}
              />
              {handleLocationChange && (
                <TextField
                  className={styles.countyInput}
                  label="County"
                  id={`county_${addressFieldsUUID.current}`}
                  name={`${name}.county`}
                  labelHint={labelHintProp}
                  data-testid={`${name}.county`}
                  validate={validators?.county}
                  isDisabled={zipCityEnabled}
                />
              )}
            </div>
            <div className="mobile-lg:grid-col-6">
              {stateField}
              {postalCodeField}
            </div>
          </div>
        </>,
      )}
    </Fieldset>
  );
};

AddressFields.propTypes = {
  legend: PropTypes.node,
  className: PropTypes.string,
  name: PropTypes.string.isRequired,
  render: PropTypes.func,
  zipCityEnabled: PropTypes.bool,
  zipCityError: PropTypes.bool,
  handleLocationChange: PropTypes.func,
  validators: PropTypes.shape({
    streetAddress1: PropTypes.func,
    streetAddress2: PropTypes.func,
    city: PropTypes.func,
    state: PropTypes.func,
    postalCode: PropTypes.func,
    county: PropTypes.func,
  }),
  formikFunctionsToValidatePostalCodeOnChange: PropTypes.shape({
    handleChange: PropTypes.func,
    setFieldTouched: PropTypes.func,
  }),
};

AddressFields.defaultProps = {
  legend: '',
  className: '',
  render: (fields) => fields,
  zipCityEnabled: false,
  zipCityError: false,
  handleLocationChange: null,
  validators: {},
  formikFunctionsToValidatePostalCodeOnChange: null,
};

export default AddressFields;
