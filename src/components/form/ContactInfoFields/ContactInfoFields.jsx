import React from 'react';
import PropTypes from 'prop-types';
import { Fieldset, Label, TextInput } from '@trussworks/react-uswds';

export const ContactInfoFields = ({ legend, className }) => {
  return (
    <Fieldset legend={legend} className={className}>
      <Label hint="(optional)" htmlFor="first-name">
        First name
      </Label>
      <TextInput id="first-name" name="first-name" type="text" />
      <Label hint="(optional)" htmlFor="last-name">
        Last name
      </Label>
      <TextInput id="last-name" name="last-name" type="text" />
      <Label hint="(optional)" htmlFor="phone">
        Phone
      </Label>
      <TextInput id="phone" name="phone" type="text" />
      <Label hint="(optional)" htmlFor="state">
        Email
      </Label>
      <TextInput id="email" name="email" type="text" />
    </Fieldset>
  );
};

ContactInfoFields.propTypes = {
  legend: PropTypes.string,
  className: PropTypes.string,
};

ContactInfoFields.defaultProps = {
  legend: '',
  className: '',
};

export default ContactInfoFields;
