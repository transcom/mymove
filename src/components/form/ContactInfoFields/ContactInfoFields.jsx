import React from 'react';
import PropTypes from 'prop-types';
import { Field } from 'formik';
import { Fieldset } from '@trussworks/react-uswds';
import { v4 as uuidv4 } from 'uuid';

import { TextInput } from 'components/form/fields';

export const ContactInfoFields = ({ legend, className, subtitle, values, handleChange, name }) => {
  const contactInfoFieldsUUID = uuidv4();

  return (
    <Fieldset legend={legend} className={className}>
      {subtitle && <span>{subtitle}</span>}
      <Field
        as={TextInput}
        label="First name"
        labelHint="(optional)"
        id={`firstName_${contactInfoFieldsUUID}`}
        data-testid="firstName"
        name={`${name}.firstName`}
        type="text"
        onChange={handleChange}
        value={values.firstName}
      />
      <Field
        as={TextInput}
        label="Last name"
        labelHint="(optional)"
        id={`lastName_${contactInfoFieldsUUID}`}
        data-testid="lastName"
        name={`${name}.lastName`}
        type="text"
        onChange={handleChange}
        value={values.lastName}
      />

      <Field
        label="Phone"
        labelHint="(optional)"
        id={`phone_${contactInfoFieldsUUID}`}
        data-testid="phone"
        name={`${name}.phone`}
        type="text"
        onChange={handleChange}
        value={values.phone}
        validationStatus="error"
      />
      <Field
        as={TextInput}
        label="Email"
        labelHint="(optional)"
        id={`email_${contactInfoFieldsUUID}`}
        data-testid="email"
        name={`${name}.email`}
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
  name: PropTypes.string.isRequired,
  handleChange: PropTypes.func.isRequired,
};

ContactInfoFields.defaultProps = {
  legend: '',
  className: '',
  subtitle: '',
  values: {},
};

export default ContactInfoFields;
