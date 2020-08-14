import React from 'react';
import PropTypes from 'prop-types';
import { Field } from 'formik';
import { v4 as uuidv4 } from 'uuid';

import Fieldset from 'shared/Fieldset';
import { TextInput } from 'components/form/fields';

export const ContactInfoFields = ({ legend, className, subtitle, values, name, subtitleClassName, hintText }) => {
  const contactInfoFieldsUUID = uuidv4();

  return (
    <Fieldset legend={legend} className={className} hintText={hintText}>
      {subtitle && <div className={subtitleClassName}>{subtitle}</div>}
      <Field
        as={TextInput}
        label="First name"
        id={`firstName_${contactInfoFieldsUUID}`}
        data-testid="firstName"
        name={`${name}.firstName`}
        type="text"
        value={values.firstName}
      />
      <Field
        as={TextInput}
        label="Last name"
        id={`lastName_${contactInfoFieldsUUID}`}
        data-testid="lastName"
        name={`${name}.lastName`}
        type="text"
        value={values.lastName}
      />

      <Field
        as={TextInput}
        label="Phone"
        id={`phone_${contactInfoFieldsUUID}`}
        data-testid="phone"
        name={`${name}.phone`}
        type="tel"
        maxLength="10"
        value={values.phone}
      />
      <Field
        as={TextInput}
        label="Email"
        id={`email_${contactInfoFieldsUUID}`}
        data-testid="email"
        name={`${name}.email`}
        type="text"
        value={values.email}
      />
    </Fieldset>
  );
};

ContactInfoFields.propTypes = {
  legend: PropTypes.string,
  hintText: PropTypes.string,
  className: PropTypes.string,
  subtitle: PropTypes.string,
  subtitleClassName: PropTypes.string,
  values: PropTypes.shape({
    firstName: PropTypes.string,
    lastName: PropTypes.string,
    phone: PropTypes.string,
    email: PropTypes.string,
  }),
  name: PropTypes.string.isRequired,
};

ContactInfoFields.defaultProps = {
  legend: '',
  className: '',
  hintText: '',
  subtitle: '',
  subtitleClassName: '',
  values: {},
};

export default ContactInfoFields;
