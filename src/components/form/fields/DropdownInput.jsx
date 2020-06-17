import { useField } from 'formik';
import { Dropdown, FormGroup, Label } from '@trussworks/react-uswds';
import React from 'react';
import PropTypes from 'prop-types';

import { ErrorMessage } from 'components/form/ErrorMessage';

export const DropdownInput = (props) => {
  // eslint-disable-next-line react/prop-types
  const { label, options } = props;
  const [field, meta] = useField(props);
  const hasError = meta.touched && !!meta.error;

  return (
    <FormGroup error={hasError}>
      <Label error={hasError} htmlFor={field.name}>
        {label}
      </Label>
      <ErrorMessage display={hasError}>{meta.error}</ErrorMessage>
      {/* eslint-disable-next-line react/jsx-props-no-spreading */}
      <Dropdown {...field}>
        <option value="">- Select -</option>
        {options &&
          // eslint-disable-next-line react/prop-types
          options.map(([optionValue, optionLabel]) => (
            <option key={optionValue} value={optionValue}>
              {optionLabel}
            </option>
          ))}
      </Dropdown>
    </FormGroup>
  );
};

DropdownInput.propTypes = {
  // label displayed for input
  label: PropTypes.string.isRequired,
  // name is for the input
  name: PropTypes.string.isRequired,
  // options for dropdown selection for this input
  // ex: [ [ "key", "value" ] ]
  options: PropTypes.arrayOf(PropTypes.arrayOf(PropTypes.string)).isRequired,
};

export default DropdownInput;
