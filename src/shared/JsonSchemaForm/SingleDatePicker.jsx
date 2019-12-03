import React from 'react';
import DayPickerInput from 'react-day-picker/DayPickerInput';
import moment from 'moment';
import { defaultDateFormat } from 'shared/utils';

import 'react-day-picker/lib/style.css';

// First date format is take to be the default
const allowedDateFormats = [
  defaultDateFormat,
  'YYYY/M/D',
  // 'YYYY-M-D',
  // 'M-D-YYYY',
  // 'D-MMM-YYYY',
  // 'MMM-D-YYYY',
  // 'DD-MMM-YY',
];

function parseDate(str, _format, locale = 'en') {
  // Ignore default format, and attempt to parse date using allowed formats
  const m = moment(str, allowedDateFormats, locale, true);
  if (m.isValid()) {
    return m.toDate();
  }

  return undefined;
}

function formatDate(date, format = defaultDateFormat, locale = 'en') {
  return moment(date, allowedDateFormats, locale, true)
    .locale(locale)
    .format(format);
}

const getDayPickerProps = disabledDays => {
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
  const formatted = parseDate(value);

  return (
    <DayPickerInput
      onDayChange={onChange}
      onDayPickerHide={onBlur}
      placeholder={placeholder}
      parseDate={parseDate}
      formatDate={formatDate}
      format={format}
      value={formatted}
      dayPickerProps={getDayPickerProps(disabledDays)}
      inputProps={{ disabled, name, onChange, autoComplete: 'off', className: inputClassName || 'usa-input' }}
    />
  );
}
