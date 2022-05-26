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
  const { label, name, id, className, renderInput, disabled, required } = props;
  const [field, meta, helpers] = useField(props);
  const hasError = meta.touched && !!meta.error;

  // Input elements need an ID prop to be associated with the label
  const inputId = useRef(id || `${name}_${uuidv4()}`);

  return (
    <FormGroup error={hasError}>
      {renderInput(
        <>
          <div className="labelWrapper">
            <Label error={hasError} htmlFor={inputId.current}>
              {label}
            </Label>
          </div>
          <ErrorMessage display={hasError}>{meta.error}</ErrorMessage>
          <SingleDatePicker
            title={label}
            name={name}
            id={inputId.current}
            inputClassName={className}
            placeholder={dateFormat}
            format={dateFormat}
            onChange={(value) => helpers.setValue(formatDate(value, dateFormat))}
            onBlur={() => helpers.setTouched(true)}
            value={field.value}
            required={required}
            disabled={disabled}
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
  className: PropTypes.string,
  renderInput: PropTypes.func,
  disabled: PropTypes.bool,
  required: PropTypes.bool,
};

DatePickerInput.defaultProps = {
  renderInput: (component) => component,
  id: undefined,
  className: undefined,
  disabled: false,
  required: false,
};

export default DatePickerInput;
