import React from 'react';
import { string, func, bool } from 'prop-types';
import { uniqueId } from 'lodash';

const Checkbox = ({ name, label, onChange, value, checked, inputClassName, labelClassName }) => {
  const checkboxId = uniqueId(label);
  return (
    <>
      <input
        className={inputClassName}
        id={checkboxId}
        type="checkbox"
        name={name}
        value={value}
        checked={checked}
        onChange={onChange}
      />
      <label className={labelClassName} htmlFor={checkboxId}>
        {label}
      </label>
    </>
  );
};

Checkbox.propTypes = {
  name: string.isRequired,
  label: string.isRequired,
  onChange: func.isRequired,
  checked: bool.isRequired,
  inputClassName: string,
  labelClassName: string,
};

export default Checkbox;
