import { useField } from 'formik';
import { FormGroup, Label } from '@trussworks/react-uswds';
import { ErrorMessage } from 'components/form/ErrorMessage';
import SingleDatePicker from 'shared/JsonSchemaForm/SingleDatePicker';
import { formatDate } from 'shared/dates';
import React from 'react';
import PropTypes from 'prop-types';

export const DatePickerInput = (props) => {
  const dateFormat = 'DD MMM YYYY';
  // eslint-disable-next-line react/prop-types
  const { label, name } = props;
  const [field, meta, helpers] = useField(props);
  const hasError = meta.touched && !!meta.error;

  return (
    <FormGroup error={hasError}>
      <Label error={hasError} htmlFor={field.name}>
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
  // label optionally displayed for input
  label: PropTypes.string,
  // name is for the input
  name: PropTypes.string.isRequired,
};

export default DatePickerInput;
