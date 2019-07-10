import React from 'react';
import { string, func, bool } from 'prop-types';
import { uniqueId } from 'lodash';
import classNames from 'classnames';

import './Checkbox.css';

const Checkbox = ({ name, label, onChange, value, checked, inputClassName, labelClassName, normalizeLabel }) => {
  const checkboxId = uniqueId(label);
  const labelClasses = classNames({ 'normalize-label': normalizeLabel }, labelClassName);
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
      <label className={labelClasses} htmlFor={checkboxId}>
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
  normalizeLabel: bool,
};

export default Checkbox;
