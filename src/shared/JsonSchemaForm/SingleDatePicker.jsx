import React from 'react';
import DayPickerInput from 'react-day-picker/DayPickerInput';
import moment from 'moment';

import 'react-day-picker/lib/style.css';

// First date format is take to be the default
const allowedDateFormats = ['M/D/YYYY', 'YYYY/M/D', 'YYYY-M-D', 'M-D-YYYY'];
export const defaultDateFormat = allowedDateFormats[0];

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

export default function SingleDatePicker(props) {
  const { value = null, onChange, disabled } = props;
  return (
    <DayPickerInput
      onDayChange={onChange}
      placeholder=""
      parseDate={parseDate}
      format={defaultDateFormat}
      value={value}
      inputProps={{ disabled }}
    />
  );
}
