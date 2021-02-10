import React from 'react';
import DayPickerInput from 'react-day-picker/DayPickerInput';
import 'react-day-picker/lib/style.css';

import { parseDate, formatDateFromISO, defaultDateFormat } from 'shared/dates';

const getDayPickerProps = (disabledDays) => {
  return {
    modifiers: {
      disabled: disabledDays,
    },
  };
};

export default function SingleDatePicker(props) {
  const {
    value = null,
    format = defaultDateFormat,
    onChange,
    onBlur,
    disabled,
    name,
    disabledDays,
    placeholder,
    inputClassName,
  } = props;
  const parsedValue = parseDate(value);

  return (
    <DayPickerInput
      onDayChange={onChange}
      onDayPickerHide={onBlur}
      placeholder={placeholder}
      parseDate={parseDate}
      formatDate={formatDateFromISO}
      format={format}
      value={parsedValue}
      dayPickerProps={getDayPickerProps(disabledDays)}
      inputProps={{
        disabled,
        name,
        autoComplete: 'off',
        className: inputClassName || 'usa-input',
      }}
    />
  );
}
