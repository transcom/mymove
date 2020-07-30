import React from 'react';
import PropTypes from 'prop-types';
import { Fieldset, Label } from '@trussworks/react-uswds';
import { v4 as uuidv4 } from 'uuid';

import { TextInput } from 'components/form/fields';

export const AddressFields = ({ legend, className, values, name, renderExistingAddressCheckbox }) => {
  const addressFieldsUUID = uuidv4();

  return (
    <Fieldset legend={legend} className={className}>
      {renderExistingAddressCheckbox()}
      <TextInput
        label="Street address 1"
        id={`mailingAddress1_${addressFieldsUUID}`}
        data-testid="mailingAddress1"
        name={`${name}.mailingAddress1`}
        type="text"
        value={values.mailingAddress1}
      />
      <TextInput
        label="Street address 2"
        labelHint=" (optional)"
        id={`mailingAddress2_${addressFieldsUUID}`}
        data-testid="mailingAddress2"
        name={`${name}.mailingAddress2`}
        type="text"
        value={values.mailingAddress2}
      />
      <TextInput
        label="City"
        id={`city_${addressFieldsUUID}`}
        data-testid="city"
        name={`${name}.city`}
        type="text"
        value={values.city}
      />
      <TextInput
        label="State"
        id={`state_${addressFieldsUUID}`}
        data-testid="state"
        name={`${name}.state`}
        type="text"
        value={values.state}
      />
      <TextInput
        label="ZIP"
        id={`zip_${addressFieldsUUID}`}
        data-testid="zip"
        inputSize="medium"
        name={`${name}.zip`}
        pattern="[\d]{5}(-[\d]{4})?"
        type="text"
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
  name: PropTypes.string.isRequired,
  renderExistingAddressCheckbox: PropTypes.func,
};

AddressFields.defaultProps = {
  legend: '',
  className: '',
  values: {},
  renderExistingAddressCheckbox: () => {},
};

export default AddressFields;
