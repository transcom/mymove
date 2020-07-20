import React from 'react';
import PropTypes from 'prop-types';
import { Fieldset, Label, TextInput } from '@trussworks/react-uswds';

export const AddressFields = ({ legend, className, values, handleChange }) => {
  return (
    <Fieldset legend={legend} className={className}>
      <Label htmlFor="mailing-address-1">Street address 1</Label>
      <TextInput
        id="mailing-address-1"
        data-testid="mailingAddress1"
        name="mailingAddress1"
        type="text"
        onChange={handleChange}
        value={values.mailingAddress1}
      />
      <Label hint=" (optional)" htmlFor="mailing-address-2">
        Street address 2
      </Label>
      <TextInput
        id="mailing-address-2"
        data-testid="mailingAddress2"
        name="mailingAddress2"
        type="text"
        onChange={handleChange}
        value={values.mailingAddress2}
      />
      <Label htmlFor="city">City</Label>
      <TextInput id="city" data-testid="city" name="city" type="text" onChange={handleChange} value={values.city} />
      <Label htmlFor="state">State</Label>
      <TextInput id="state" data-testid="state" name="state" type="text" onChange={handleChange} value={values.state} />
      <Label htmlFor="zip">ZIP</Label>
      <TextInput
        id="zip"
        data-testid="zip"
        inputSize="medium"
        name="zip"
        pattern="[\d]{5}(-[\d]{4})?"
        type="text"
        onChange={handleChange}
        value={values.zip}
      />
    </Fieldset>
  );
};

AddressFields.propTypes = {
  legend: PropTypes.string,
  className: PropTypes.string,
  values: PropTypes.shape({
    mailingAddress1: PropTypes.string,
    mailingAddress2: PropTypes.string,
    city: PropTypes.string,
    state: PropTypes.string,
    zip: PropTypes.string,
  }),
  handleChange: PropTypes.func.isRequired,
};

AddressFields.defaultProps = {
  legend: '',
  className: '',
  values: {},
};

export default AddressFields;
