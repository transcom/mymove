import React from 'react';
import { uniqueId } from 'lodash';

const RadioButton = ({ name, label, onChange, value, checked, inputClassName, labelClassName }) => {
  const radioId = uniqueId(label);
  return (
    <>
      <input
        className={inputClassName}
        id={radioId}
        type="radio"
        name={name}
        value={value}
        checked={checked}
        onChange={onChange}
      />
      <label className={labelClassName} htmlFor={radioId}>
        {label}
      </label>
    </>
  );
};

export default RadioButton;
