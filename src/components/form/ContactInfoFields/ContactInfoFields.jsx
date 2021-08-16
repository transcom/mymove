import React from 'react';
import PropTypes from 'prop-types';
import { v4 as uuidv4 } from 'uuid';
import { Fieldset } from '@trussworks/react-uswds';

import TextField from 'components/form/fields/TextField';
import MaskedTextField from 'components/form/fields/MaskedTextField';

export const ContactInfoFields = ({ legend, className, name, render }) => {
  const contactInfoFieldsUUID = uuidv4();

  return (
    <Fieldset legend={legend} className={className}>
      {render(
        <>
          <TextField label="First name" id={`firstName_${contactInfoFieldsUUID}`} name={`${name}.firstName`} />
          <TextField label="Last name" id={`lastName_${contactInfoFieldsUUID}`} name={`${name}.lastName`} />

          <MaskedTextField
            label="Phone"
            id={`phone_${contactInfoFieldsUUID}`}
            name={`${name}.phone`}
            type="tel"
            minimum="12"
            mask="000{-}000{-}0000"
            required
          />

          <TextField label="Email" id={`email_${contactInfoFieldsUUID}`} name={`${name}.email`} />
        </>,
      )}
    </Fieldset>
  );
};

ContactInfoFields.propTypes = {
  legend: PropTypes.node,
  className: PropTypes.string,
  name: PropTypes.string.isRequired,
  render: PropTypes.func,
};

ContactInfoFields.defaultProps = {
  legend: '',
  className: '',
  render: (fields) => fields,
};

export default ContactInfoFields;
