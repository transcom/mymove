import React from 'react';
import { string, func, bool } from 'prop-types';
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

RadioButton.propTypes = {
  name: string.isRequired,
  label: string.isRequired,
  onChange: func.isRequired,
  value: string.isRequired,
  checked: bool.isRequired,
  inputClassName: string,
  labelClassName: string,
};

export default RadioButton;
