import React, { useRef } from 'react';
import { func, node, string } from 'prop-types';
import { v4 as uuidv4 } from 'uuid';
import { Fieldset } from '@trussworks/react-uswds';

import TextField from 'components/form/fields/TextField';
import MaskedTextField from 'components/form/fields/MaskedTextField';

export const CustomerAltContactInfoFields = ({ legend, className, render }) => {
  const CustomerAltContactInfoFieldsUUID = useRef(uuidv4());

  return (
    <Fieldset legend={legend} className={className}>
      {render(
        <>
          <div className="grid-row grid-gap">
            <div className="grid-col-6">
              <TextField label="First name" name="first_name" id="firstName" required />
            </div>
            <div className="grid-col-6">
              <TextField label="Middle name" name="middle_name" id="middleName" labelHint="Optional" />
            </div>
            <div className="grid-col-6">
              <TextField label="Last name" name="last_name" id="lastName" required />
            </div>
            <div className="grid-col-6">
              <TextField label="Suffix" name="suffix" id="suffix" labelHint="Optional" />
            </div>
          </div>
          <div className="grid-row grid-gap">
            <div className="mobile-lg:grid-col-7">
              <MaskedTextField
                label="Phone"
                id={`customer_telephone_${CustomerAltContactInfoFieldsUUID.current}`}
                name="customer_telephone"
                type="tel"
                minimum="12"
                mask="000{-}000{-}0000"
                required
              />
            </div>
          </div>

          <TextField
            label="Email"
            id={`customer_email_${CustomerAltContactInfoFieldsUUID.current}`}
            name="customer_email"
            required
          />
        </>,
      )}
    </Fieldset>
  );
};

CustomerAltContactInfoFields.propTypes = {
  legend: node,
  className: string,
  render: func,
};

CustomerAltContactInfoFields.defaultProps = {
  legend: '',
  className: '',
  render: (fields) => fields,
};

export default CustomerAltContactInfoFields;
