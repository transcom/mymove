import { useField } from 'formik';
import { FormGroup, Label } from '@trussworks/react-uswds';
import React from 'react';
import PropTypes from 'prop-types';

import { ErrorMessage } from 'components/form/ErrorMessage';
import SingleDatePicker from 'shared/JsonSchemaForm/SingleDatePicker';
import { formatDate } from 'shared/dates';

export const DatePickerInput = (props) => {
  const dateFormat = 'DD MMM YYYY';
  //  react/prop-types
  const { label, name, labelClassName } = props;
  const [field, meta, helpers] = useField(props);
  const hasError = meta.touched && !!meta.error;
  return (
    <FormGroup error={hasError}>
      <Label className={labelClassName} error={hasError} htmlFor={field.name}>
        {label}
      </Label>
      <ErrorMessage display={hasError}>{meta.error}</ErrorMessage>
      <SingleDatePicker
        title={label}
        name={name}
        placeholder={dateFormat}
        format={dateFormat}
        onChange={(value) => helpers.setValue(formatDate(value, dateFormat))}
        onBlur={() => helpers.setTouched(true)}
        value={field.value}
      />
    </FormGroup>
  );
};

DatePickerInput.propTypes = {
  // label displayed for input
  label: PropTypes.string.isRequired,
  // name is for the input
  name: PropTypes.string.isRequired,
};

export default DatePickerInput;
