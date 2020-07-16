import React from 'react';
import PropTypes from 'prop-types';
import { Fieldset } from '@trussworks/react-uswds';

import { TextInput } from 'components/form/fields/TextInput';

export const AddressFields = ({ addressType, legend }) => {
  return (
    <Fieldset legend={legend}>
      <TextInput
        label="Street address 1"
        id={`${addressType}-mailing-address-1`}
        name={`${addressType}-mailing-address-1`}
        type="text"
      />
      <TextInput
        label="Street address 2"
        id={`${addressType}-mailing-address-2`}
        name={`${addressType}-mailing-address-2`}
        type="text"
        hint=" (optional)"
      />
      <TextInput label="City" id="city" name={`${addressType}-city`} type="text" />
      <TextInput label="State" id="state" name={`${addressType}-state`} type="text" />
      <TextInput
        label="ZIP"
        id="zip"
        inputSize="medium"
        name={`${addressType}-zip`}
        pattern="[\d]{5}(-[\d]{4})?"
        type="text"
      />
    </Fieldset>
  );
};

AddressFields.propTypes = {
  addressType: PropTypes.string,
  legend: PropTypes.string,
};

AddressFields.defaultProps = {
  addressType: '',
  legend: '',
};

export default AddressFields;
