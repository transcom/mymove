import React from 'react';
import PropTypes from 'prop-types';
import { v4 as uuidv4 } from 'uuid';
import { Fieldset } from '@trussworks/react-uswds';

import TextField from 'components/form/fields/TextField';

export const AddressFields = ({ legend, className, name, render }) => {
  const addressFieldsUUID = uuidv4();

  return (
    <Fieldset legend={legend} className={className}>
      {render(
        <>
          <TextField label="Address 1" id={`mailingAddress1_${addressFieldsUUID}`} name={`${name}.street_address_1`} />
          <TextField
            label="Address 2"
            labelHint="Optional"
            id={`mailingAddress2_${addressFieldsUUID}`}
            name={`${name}.street_address_2`}
          />
          <TextField label="City" id={`city_${addressFieldsUUID}`} name={`${name}.city`} />

          <div className="grid-row grid-gap">
            <div className="mobile-lg:grid-col-6">
              <TextField label="State" id={`state_${addressFieldsUUID}`} name={`${name}.state`} maxLength={2} />
            </div>
            <div className="mobile-lg:grid-col-6">
              <TextField label="ZIP" id={`zip_${addressFieldsUUID}`} name={`${name}.postal_code`} maxLength={10} />
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
};

AddressFields.defaultProps = {
  legend: '',
  className: '',
  render: (fields) => fields,
};

export default AddressFields;
