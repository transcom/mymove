import React, { useRef } from 'react';
import PropTypes from 'prop-types';
import { v4 as uuidv4 } from 'uuid';
import { useField } from 'formik';
import { Dropdown, FormGroup, Label } from '@trussworks/react-uswds';

import { ErrorMessage } from 'components/form/ErrorMessage';
// import { OptionalTag } from 'components/form/OptionalTag';
import { DropdownArrayOf } from 'types/form';
import './DropdownInput.module.scss';

export const DropdownInput = (props) => {
  const { id, name, label, options, showDropdownPlaceholderText, isDisabled, ...inputProps } = props;
  const [field, meta] = useField(props);
  const hasError = meta.touched && !!meta.error;

  // Input elements need an ID prop to be associated with the label
  const inputId = useRef(id || `${name}_${uuidv4()}`);

  return (
    <FormGroup error={hasError}>
      <div className="labelWrapper">
        <Label error={hasError} htmlFor={inputId.current}>
          {label}
        </Label>
        {/* {optional && <OptionalTag />} */}
      </div>
      <ErrorMessage display={hasError}>{meta.error}</ErrorMessage>
      {/* eslint-disable-next-line react/jsx-props-no-spreading */}
      <Dropdown id={inputId.current} {...field} disabled={isDisabled} {...inputProps}>
        {showDropdownPlaceholderText && <option value="">- Select -</option>}
        {options &&
          options.map(({ key, value }) => (
            <option key={key} value={key}>
              {value}
            </option>
          ))}
      </Dropdown>
    </FormGroup>
  );
};

DropdownInput.propTypes = {
  // label displayed for input
  label: PropTypes.string.isRequired,
  id: PropTypes.string,
  // name is for the input
  name: PropTypes.string.isRequired,
  // options for dropdown selection for this input
  // ex: [ { key: 'key', value: 'value' } ]
  options: DropdownArrayOf.isRequired,
  showDropdownPlaceholderText: PropTypes.bool,
  isDisabled: PropTypes.bool,
};

DropdownInput.defaultProps = {
  id: undefined,
  showDropdownPlaceholderText: true,
  isDisabled: false,
};

export default DropdownInput;
