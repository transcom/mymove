import React from 'react';
import PropTypes from 'prop-types';
import { Field } from 'formik';
import { v4 as uuidv4 } from 'uuid';
import { Fieldset } from '@trussworks/react-uswds';

import { TextInput } from 'components/form/fields';

export const AddressFields = ({ legend, className, values, name, render }) => {
  const addressFieldsUUID = uuidv4();

  return (
    <Fieldset legend={legend} className={className}>
      {render(
        <>
          <Field
            as={TextInput}
            label="Address 1"
            id={`mailingAddress1_${addressFieldsUUID}`}
            data-testid="mailingAddress1"
            name={`${name}.street_address_1`}
            type="text"
            value={values.street_address_1}
          />
          <Field
            as={TextInput}
            label="Address 2"
            labelHint="Optional"
            id={`mailingAddress2_${addressFieldsUUID}`}
            data-testid="mailingAddress2"
            name={`${name}.street_address_2`}
            type="text"
            value={values.street_address_2}
          />
          <Field
            as={TextInput}
            label="City"
            id={`city_${addressFieldsUUID}`}
            data-testid="city"
            name={`${name}.city`}
            type="text"
            value={values.city}
          />
          <Field
            as={TextInput}
            label="State"
            id={`state_${addressFieldsUUID}`}
            data-testid="state"
            name={`${name}.state`}
            type="text"
            value={values.state}
            maxLength={2}
          />
          <Field
            as={TextInput}
            label="ZIP"
            id={`zip_${addressFieldsUUID}`}
            data-testid="zip"
            inputSize="medium"
            name={`${name}.postal_code`}
            type="text"
            value={values.postal_code}
            maxLength={10}
          />
        </>,
      )}
    </Fieldset>
  );
};

AddressFields.propTypes = {
  legend: PropTypes.node,
  className: PropTypes.string,
  values: PropTypes.shape({
    street_address_1: PropTypes.string,
    street_address_2: PropTypes.string,
    city: PropTypes.string,
    state: PropTypes.string,
    postal_code: PropTypes.string,
  }),
  name: PropTypes.string.isRequired,
  render: PropTypes.func,
};

AddressFields.defaultProps = {
  legend: '',
  className: '',
  values: {},
  render: (fields) => fields,
};

export default AddressFields;
