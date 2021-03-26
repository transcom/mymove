import React, { useRef } from 'react';
import PropTypes from 'prop-types';
import { useField } from 'formik';
import { FormGroup, Label } from '@trussworks/react-uswds';
import { v4 as uuidv4 } from 'uuid';

import { ErrorMessage } from 'components/form/ErrorMessage';
import SingleDatePicker from 'shared/JsonSchemaForm/SingleDatePicker';
import { formatDate } from 'shared/dates';

export const DatePickerInput = (props) => {
  const dateFormat = 'DD MMM YYYY';
  const { label, name, id, renderInput } = props;
  const [field, meta, helpers] = useField(props);
  const hasError = meta.touched && !!meta.error;

  // Input elements need an ID prop to be associated with the label
  const inputId = useRef(id || `${name}_${uuidv4()}`);

  return (
    <FormGroup error={hasError}>
      {renderInput(
        <>
          <Label error={hasError} htmlFor={inputId.current}>
            {label}
          </Label>
          <ErrorMessage display={hasError}>{meta.error}</ErrorMessage>
          <SingleDatePicker
            title={label}
            name={name}
            id={inputId.current}
            placeholder={dateFormat}
            format={dateFormat}
            onChange={(value) => helpers.setValue(formatDate(value, dateFormat))}
            onBlur={() => helpers.setTouched(true)}
            value={field.value}
          />
        </>,
      )}
    </FormGroup>
  );
};

DatePickerInput.propTypes = {
  // label displayed for input
  label: PropTypes.string.isRequired,
  // name is for the input
  name: PropTypes.string.isRequired,
  id: PropTypes.string,
  renderInput: PropTypes.func,
};

DatePickerInput.defaultProps = {
  renderInput: (component) => component,
  id: undefined,
};

export default DatePickerInput;
