import React from 'react';
import DayPickerInput from 'react-day-picker/DayPickerInput';
import 'react-day-picker/lib/style.css';

export default function SingleDatePicker(props) {
  const { value = null, onChange, disabled } = props;
  return (
    <DayPickerInput
      onDayChange={onChange}
      placeholder="Date"
      value={value}
      inputProps={{ disabled }}
    />
  );
}
