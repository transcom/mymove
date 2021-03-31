import React from 'react';
import PropTypes from 'prop-types';
import { v4 as uuidv4 } from 'uuid';
import { Fieldset } from '@trussworks/react-uswds';

import TextField from 'components/form/fields/TextField';

export const ContactInfoFields = ({ legend, className, name, render }) => {
  const contactInfoFieldsUUID = uuidv4();

  return (
    <Fieldset legend={legend} className={className}>
      {render(
        <>
          <TextField label="First name" id={`firstName_${contactInfoFieldsUUID}`} name={`${name}.firstName`} />
          <TextField label="Last name" id={`lastName_${contactInfoFieldsUUID}`} name={`${name}.lastName`} />

          <TextField
            label="Phone"
            id={`phone_${contactInfoFieldsUUID}`}
            name={`${name}.phone`}
            type="tel"
            maxLength="10"
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
