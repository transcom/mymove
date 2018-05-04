import React from 'react';
import DayPickerInput from 'react-day-picker/DayPickerInput';
import 'react-day-picker/lib/style.css';

export default function SingleDatePicker(props) {
  const {
    input: { value = null, onChange },
  } = props;
  return (
    <DayPickerInput onDayChange={onChange} placeholder="Date" value={value} />
  );
}
