import React from 'react';
import { func, node, shape, string } from 'prop-types';
import { Field } from 'formik';
import { v4 as uuidv4 } from 'uuid';
import { Checkbox, Fieldset } from '@trussworks/react-uswds';

import { TextInput } from 'components/form/fields';

export const ContactInfoFieldset = ({
  legend,
  className,
  onChangePreferEmail,
  onChangePreferPhone,
  values,
  name,
  render,
}) => {
  const contactInfoFieldsetUUID = uuidv4();

  return (
    <Fieldset legend={legend} className={className}>
      {render(
        <>
          <Field
            as={TextInput}
            label="Best contact phone"
            id={`phone_${contactInfoFieldsetUUID}`}
            data-testid="phone"
            name={`${name}.phone`}
            type="tel"
            maxLength="10"
            value={values.phone}
          />
          <Field
            as={TextInput}
            label="Alt. phone"
            id={`alternate_phone_${contactInfoFieldsetUUID}`}
            data-testid="alternamte-phone"
            name={`${name}.alternate_phone`}
            type="tel"
            maxLength="10"
            value={values.alternatePhone}
          />
          <Field
            as={TextInput}
            label="Personal email"
            id={`email_${contactInfoFieldsetUUID}`}
            data-testid="email"
            name={`${name}.email`}
            type="text"
            value={values.email}
          />
          <p>Preferred contact method</p>
          <Checkbox
            id={`prefer_phone_${contactInfoFieldsetUUID}`}
            label="Phone"
            name={`${name}.prefer_phone`}
            onChange={onChangePreferPhone}
          />
          <Checkbox
            id={`prefer_email${contactInfoFieldsetUUID}`}
            label="Email"
            name={`${name}.prefer_email`}
            onChange={onChangePreferEmail}
          />
        </>,
      )}
    </Fieldset>
  );
};

ContactInfoFieldset.propTypes = {
  legend: node,
  className: string,
  values: shape({
    phone: string,
    alternatePhone: string,
    email: string,
  }),
  name: string.isRequired,
  render: func,
  onChangePreferPhone: func.isRequired,
  onChangePreferEmail: func.isRequired,
};

ContactInfoFieldset.defaultProps = {
  legend: '',
  className: '',
  values: {},
  render: (fields) => fields,
};

export default ContactInfoFieldset;
