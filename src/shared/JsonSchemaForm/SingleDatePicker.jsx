import React from 'react';
import DayPickerInput from 'react-day-picker/DayPickerInput';
import moment from 'moment';
import { defaultDateFormat } from 'shared/utils';

import 'react-day-picker/lib/style.css';

// First date format is take to be the default
const allowedDateFormats = [defaultDateFormat, 'YYYY/M/D', 'YYYY-M-D', 'M-D-YYYY', 'D-MMM-YYYY', 'MMM-D-YYYY'];

function parseDate(str, _format, locale = 'en') {
  // Ignore default format, and attempt to parse date using allowed formats
  for (var i = 0; i < allowedDateFormats.length; i++) {
    let format = allowedDateFormats[i]; // eslint-disable-line security/detect-object-injection
    const m = moment(str, format, locale, true);
    if (m.isValid()) {
      return m.toDate();
    }
  }
  return undefined;
}

function formatDate(date, format = 'L', locale = 'en') {
  return moment(date)
    .locale(locale)
    .format(defaultDateFormat);
}

const getDayPickerProps = disabledDays => {
  return {
    modifiers: {
      disabled: disabledDays,
    },
  };
};

export default function SingleDatePicker(props) {
  const { value = null, onChange, disabled, name, disabledDays } = props;
  const formatted = parseDate(value);
  return (
    <DayPickerInput
      onDayChange={onChange}
      placeholder=""
      parseDate={parseDate}
      formatDate={formatDate}
      value={formatted}
      dayPickerProps={getDayPickerProps(disabledDays)}
      inputProps={{ disabled, name }}
    />
  );
}
