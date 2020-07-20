import React from 'react';
import PropTypes from 'prop-types';
import { Fieldset, Label, TextInput } from '@trussworks/react-uswds';
import { v4 as uuidv4 } from 'uuid';

export const ContactInfoFields = ({ legend, className, subtitle, values, handleChange }) => {
  const contactInfoFieldsUUID = uuidv4();

  return (
    <Fieldset legend={legend} className={className}>
      {subtitle && <span>{subtitle}</span>}
      <Label hint="(optional)" htmlFor={`firstName_${contactInfoFieldsUUID}`}>
        First name
      </Label>
      <TextInput
        id={`firstName_${contactInfoFieldsUUID}`}
        data-testid="firstName"
        name="firstName"
        type="text"
        onChange={handleChange}
        value={values.firstName}
      />
      <Label hint="(optional)" htmlFor={`lastName_${contactInfoFieldsUUID}`}>
        Last name
      </Label>
      <TextInput
        id={`lastName_${contactInfoFieldsUUID}`}
        data-testid="lastName"
        name="lastName"
        type="text"
        onChange={handleChange}
        value={values.lastName}
      />
      <Label hint="(optional)" htmlFor={`phone_${contactInfoFieldsUUID}`}>
        Phone
      </Label>
      <TextInput
        id={`phone_${contactInfoFieldsUUID}`}
        data-testid="phone"
        name="phone"
        type="text"
        onChange={handleChange}
        value={values.phone}
      />
      <Label hint="(optional)" htmlFor={`email_${contactInfoFieldsUUID}`}>
        Email
      </Label>
      <TextInput
        id={`email_${contactInfoFieldsUUID}`}
        data-testid="email"
        name="email"
        type="text"
        onChange={handleChange}
        value={values.email}
      />
    </Fieldset>
  );
};

ContactInfoFields.propTypes = {
  legend: PropTypes.string,
  className: PropTypes.string,
  subtitle: PropTypes.string,
  values: PropTypes.shape({
    firstName: PropTypes.string,
    lastName: PropTypes.string,
    phone: PropTypes.string,
    email: PropTypes.string,
  }),
  handleChange: PropTypes.func.isRequired,
};

ContactInfoFields.defaultProps = {
  legend: '',
  className: '',
  subtitle: '',
  values: {},
};

export default ContactInfoFields;
