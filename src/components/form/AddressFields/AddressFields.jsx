import React, { useRef } from 'react';
import PropTypes from 'prop-types';
import { v4 as uuidv4 } from 'uuid';
import { Fieldset } from '@trussworks/react-uswds';

import { statesList } from '../../../constants/states';

import TextField from 'components/form/fields/TextField/TextField';
import { DropdownInput } from 'components/form/fields/DropdownInput';

/**
 * @param legend
 * @param className
 * @param name
 * @param render
 * @param validators
 * @param formikFunctionsToValidatePostalCodeOnChange If you are intending to validate the postal code on change, you
 * will need to pass the handleChange and setFieldTouched Formik functions through in an object here.
 * See ResidentialAddressForm for an example.
 * @param address1LabelHint string to override display labelHint if street 1 is Optional/Required per context.
 * This is specifically designed to handle unique display between customer and office/prime sim for address 1.
 * @return {JSX.Element}
 * @constructor
 */
export const AddressFields = ({
  legend,
  className,
  name,
  render,
  validators,
  formikFunctionsToValidatePostalCodeOnChange,
  labelHint: labelHintProp,
  address1LabelHint,
}) => {
  const addressFieldsUUID = useRef(uuidv4());

  let postalCodeField;

  if (formikFunctionsToValidatePostalCodeOnChange) {
    postalCodeField = (
      <TextField
        label="ZIP"
        id={`zip_${addressFieldsUUID.current}`}
        name={`${name}.postalCode`}
        maxLength={10}
        labelHint={labelHintProp}
        validate={validators?.postalCode}
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
    );
  } else {
    postalCodeField = (
      <TextField
        label="ZIP"
        id={`zip_${addressFieldsUUID.current}`}
        name={`${name}.postalCode`}
        maxLength={10}
        labelHint={labelHintProp}
        validate={validators?.postalCode}
      />
    );
  }

  const getAddress1LabelHintText = (labelHint, address1Label) => {
    if (address1Label === null) {
      return labelHint;
    }

    // Override default and use what is passed in.
    if (address1Label && address1Label.trim().length > 0) {
      return address1Label;
    }

    return null;
  };

  return (
    <Fieldset legend={legend} className={className}>
      {render(
        <>
          <TextField
            label="Address 1"
            id={`mailingAddress1_${addressFieldsUUID.current}`}
            name={`${name}.streetAddress1`}
            labelHint={getAddress1LabelHintText(labelHintProp, address1LabelHint)}
            validate={validators?.streetAddress1}
          />
          <TextField
            label="Address 2"
            labelHint={labelHintProp ? null : 'Optional'}
            id={`mailingAddress2_${addressFieldsUUID.current}`}
            name={`${name}.streetAddress2`}
            validate={validators?.streetAddress2}
          />
          <TextField
            label="Address 3"
            labelHint={labelHintProp ? null : 'Optional'}
            id={`mailingAddress3_${addressFieldsUUID.current}`}
            name={`${name}.streetAddress3`}
            validate={validators?.streetAddress3}
          />
          <TextField
            label="City"
            id={`city_${addressFieldsUUID.current}`}
            name={`${name}.city`}
            labelHint={labelHintProp}
            validate={validators?.city}
          />

          <div className="grid-row grid-gap">
            <div className="mobile-lg:grid-col-6">
              <DropdownInput
                name={`${name}.state`}
                id={`state_${addressFieldsUUID.current}`}
                label="State"
                labelHint={labelHintProp}
                options={statesList}
                validate={validators?.state}
              />
            </div>
            <div className="mobile-lg:grid-col-6">{postalCodeField}</div>
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
  validators: PropTypes.shape({
    streetAddress1: PropTypes.func,
    streetAddress2: PropTypes.func,
    city: PropTypes.func,
    state: PropTypes.func,
    postalCode: PropTypes.func,
  }),
  formikFunctionsToValidatePostalCodeOnChange: PropTypes.shape({
    handleChange: PropTypes.func,
    setFieldTouched: PropTypes.func,
  }),
  address1LabelHint: PropTypes.string,
};

AddressFields.defaultProps = {
  legend: '',
  className: '',
  render: (fields) => fields,
  validators: {},
  formikFunctionsToValidatePostalCodeOnChange: null,
  address1LabelHint: null,
};

export default AddressFields;
