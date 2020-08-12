import React from 'react';
import PropTypes from 'prop-types';
import { Fieldset, Label, TextInput } from '@trussworks/react-uswds';
import { v4 as uuidv4 } from 'uuid';

export const AddressFields = ({ legend, className, values, handleChange, name, renderExistingAddressCheckbox }) => {
  const addressFieldsUUID = uuidv4();

  return (
    <Fieldset legend={legend} className={className}>
      {renderExistingAddressCheckbox()}
      <Label htmlFor={`street_address_1_${addressFieldsUUID}`}>Street address 1</Label>
      <TextInput
        id={`street_address_1_${addressFieldsUUID}`}
        data-testid="mailingAddress1"
        name={`${name}.street_address_1`}
        type="text"
        onChange={handleChange}
        value={values.street_address_1}
      />
      <Label hint=" (optional)" htmlFor={`street_address_2_${addressFieldsUUID}`}>
        Street address 2
      </Label>
      <TextInput
        id={`street_address_2_${addressFieldsUUID}`}
        data-testid="mailingAddress2"
        name={`${name}.street_address_2`}
        type="text"
        onChange={handleChange}
        value={values.street_address_2}
      />
      <Label htmlFor={`city_${addressFieldsUUID}`}>City</Label>
      <TextInput
        id={`city_${addressFieldsUUID}`}
        data-testid="city"
        name={`${name}.city`}
        type="text"
        onChange={handleChange}
        value={values.city}
      />
      <Label htmlFor={`state_${addressFieldsUUID}`}>State</Label>
      <TextInput
        id={`state_${addressFieldsUUID}`}
        data-testid="state"
        name={`${name}.state`}
        type="text"
        onChange={handleChange}
        value={values.state}
      />
      <Label htmlFor={`postal_code_${addressFieldsUUID}`}>ZIP</Label>
      <TextInput
        id={`postal_code_${addressFieldsUUID}`}
        data-testid="zip"
        inputSize="medium"
        name={`${name}.postal_code`}
        pattern="[\d]{5}(-[\d]{4})?"
        type="text"
        onChange={handleChange}
        value={values.postal_code}
      />
    </Fieldset>
  );
};

AddressFields.propTypes = {
  legend: PropTypes.string,
  className: PropTypes.string,
  values: PropTypes.shape({
    street_address_1: PropTypes.string,
    street_address_2: PropTypes.string,
    city: PropTypes.string,
    state: PropTypes.string,
    postal_code: PropTypes.string,
  }),
  name: PropTypes.string.isRequired,
  handleChange: PropTypes.func.isRequired,
  renderExistingAddressCheckbox: PropTypes.func,
};

AddressFields.defaultProps = {
  legend: '',
  className: '',
  values: {},
  renderExistingAddressCheckbox: () => {},
};

export default AddressFields;
