import React, { useRef } from 'react';
import PropTypes from 'prop-types';
import { useField } from 'formik';
import { FormGroup, Label } from '@trussworks/react-uswds';
import { v4 as uuidv4 } from 'uuid';

import RequiredAsterisk from '../RequiredAsterisk';

import styles from './DatePickerInput.module.scss';

import { ErrorMessage } from 'components/form/ErrorMessage';
import SingleDatePicker from 'shared/JsonSchemaForm/SingleDatePicker';
import { datePickerFormat, formatDate } from 'shared/dates';

export const DatePickerInput = (props) => {
  const {
    label,
    showOptional,
    name,
    id,
    className,
    renderInput,
    disabled,
    disabledDays,
    required,
    hint,
    disableErrorLabel,
    onChange,
    showRequiredAsterisk,
  } = props;
  const [field, meta, helpers] = useField(props);
  const hasError = disableErrorLabel ? false : meta.touched && !!meta.error;
  const defaultOnChange = (value, _, dayPickerInput) => {
    if (value === undefined && dayPickerInput.getInput().value === '') {
      // The user cleared the date input, so we shouldn't bother attempting to format it.
      helpers.setValue(undefined);
    } else {
      helpers.setValue(formatDate(value, datePickerFormat));
    }
  };

  // Input elements need an ID prop to be associated with the label
  const inputId = useRef(id || `${name}_${uuidv4()}`);
  return (
    <FormGroup error={hasError}>
      {renderInput(
        <>
          <div className="labelWrapper">
            <Label hint={hint} className={styles.label} error={hasError} htmlFor={inputId.current}>
              <span>
                {label} {showRequiredAsterisk && <RequiredAsterisk />}
              </span>
              {showOptional && <div className={styles.optionalLabel}>Optional</div>}
            </Label>
          </div>
          <ErrorMessage display={hasError}>{meta.error}</ErrorMessage>
          <SingleDatePicker
            title={label}
            name={name}
            id={inputId.current}
            inputClassName={className}
            placeholder={datePickerFormat}
            format={datePickerFormat}
            onChange={onChange || defaultOnChange}
            onBlur={() => helpers.setTouched(true)}
            value={field.value}
            required={required}
            disabled={disabled}
            disabledDays={disabledDays}
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
  hint: PropTypes.string,
  required: PropTypes.bool,
  disableErrorLabel: PropTypes.bool,
  onChange: PropTypes.func,
};

DatePickerInput.defaultProps = {
  renderInput: (component) => component,
  id: undefined,
  className: undefined,
  disabled: false,
  required: false,
  hint: undefined,
  disableErrorLabel: false,
  onChange: undefined,
};

export default DatePickerInput;
