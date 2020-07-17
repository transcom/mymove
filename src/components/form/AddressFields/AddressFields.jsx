import React from 'react';
import PropTypes from 'prop-types';
import { Fieldset, Label, TextInput } from '@trussworks/react-uswds';

export const AddressFields = ({ legend, className }) => {
  return (
    <Fieldset legend={legend} className={className}>
      <Label htmlFor="mailing-address-1">Street address 1</Label>
      <TextInput id="mailing-address-1" data-cy="mailingAddress1" name="mailing-address-1" type="text" />
      <Label hint=" (optional)" htmlFor="mailing-address-2">
        Street address 2
      </Label>
      <TextInput id="mailing-address-2" data-cy="mailingAddress2" name="mailing-address-2" type="text" />
      <Label htmlFor="city">City</Label>
      <TextInput id="city" data-cy="city" name="city" type="text" />
      <Label htmlFor="state">State</Label>
      <TextInput id="state" data-cy="state" name="state" type="text" />
      <Label htmlFor="zip">ZIP</Label>
      <TextInput id="zip" data-cy="zip" inputSize="medium" name="zip" pattern="[\d]{5}(-[\d]{4})?" type="text" />
    </Fieldset>
  );
};

AddressFields.propTypes = {
  legend: PropTypes.string,
  className: PropTypes.string,
};

AddressFields.defaultProps = {
  legend: '',
  className: '',
};

export default AddressFields;
