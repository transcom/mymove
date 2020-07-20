import React from 'react';
import PropTypes from 'prop-types';
import { Fieldset, Label, TextInput } from '@trussworks/react-uswds';

export const ContactInfoFields = ({ legend, className, subtitle, values, handleChange }) => {
  return (
    <Fieldset legend={legend} className={className}>
      {subtitle && <span>{subtitle}</span>}
      <Label hint="(optional)" htmlFor="first-name">
        First name
      </Label>
      <TextInput
        id="first-name"
        data-testid="firstName"
        name="firstName"
        type="text"
        onChange={handleChange}
        value={values.firstName}
      />
      <Label hint="(optional)" htmlFor="last-name">
        Last name
      </Label>
      <TextInput
        id="last-name"
        data-testid="lastName"
        name="lastName"
        type="text"
        onChange={handleChange}
        value={values.lastName}
      />
      <Label hint="(optional)" htmlFor="phone">
        Phone
      </Label>
      <TextInput id="phone" data-testid="phone" name="phone" type="text" onChange={handleChange} value={values.phone} />
      <Label hint="(optional)" htmlFor="state">
        Email
      </Label>
      <TextInput id="email" data-testid="email" name="email" type="text" onChange={handleChange} value={values.email} />
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
