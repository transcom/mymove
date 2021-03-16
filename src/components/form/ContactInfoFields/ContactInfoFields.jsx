import React from 'react';
import PropTypes from 'prop-types';
import { v4 as uuidv4 } from 'uuid';
import { Fieldset } from '@trussworks/react-uswds';

import TextField from 'components/form/fields/TextField';

export const ContactInfoFields = ({ legend, className, values, name, render }) => {
  const contactInfoFieldsUUID = uuidv4();

  return (
    <Fieldset legend={legend} className={className}>
      {render(
        <>
          <TextField
            label="First name"
            id={`firstName_${contactInfoFieldsUUID}`}
            data-testid="firstName"
            name={`${name}.firstName`}
            type="text"
            value={values.firstName}
          />
          <TextField
            label="Last name"
            id={`lastName_${contactInfoFieldsUUID}`}
            data-testid="lastName"
            name={`${name}.lastName`}
            type="text"
            value={values.lastName}
          />

          <TextField
            label="Phone"
            id={`phone_${contactInfoFieldsUUID}`}
            data-testid="phone"
            name={`${name}.phone`}
            type="tel"
            maxLength="10"
            value={values.phone}
          />
          <TextField
            label="Email"
            id={`email_${contactInfoFieldsUUID}`}
            data-testid="email"
            name={`${name}.email`}
            type="text"
            value={values.email}
          />
        </>,
      )}
    </Fieldset>
  );
};

ContactInfoFields.propTypes = {
  legend: PropTypes.node,
  className: PropTypes.string,
  values: PropTypes.shape({
    firstName: PropTypes.string,
    lastName: PropTypes.string,
    phone: PropTypes.string,
    email: PropTypes.string,
  }),
  name: PropTypes.string.isRequired,
  render: PropTypes.func,
};

ContactInfoFields.defaultProps = {
  legend: '',
  className: '',
  values: {},
  render: (fields) => fields,
};

export default ContactInfoFields;
