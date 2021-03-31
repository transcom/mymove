import React from 'react';
import PropTypes from 'prop-types';
import { useField } from 'formik';
import { IMaskInput } from 'react-imask';
import { FormGroup, Label } from '@trussworks/react-uswds';

import { ErrorMessage } from '../index';

const MaskedTextField = ({
  label,
  labelClassName,
  id,
  name,
  labelHint,
  defaultValue,
  mask,
  blocks,
  lazy,
  warning,
  ...props
}) => {
  const [field, meta, helpers] = useField({ id, name, ...props });
  const hasError = meta.touched && !!meta.error;
  const { value } = field;
  return (
    <FormGroup className={!!warning && !hasError && `warning`} error={hasError}>
      <Label className={labelClassName} hint={labelHint} error={hasError} htmlFor={id || name}>
        {label}
      </Label>
      <ErrorMessage display={hasError}>{meta.error}</ErrorMessage>
      {/* eslint-disable react/jsx-props-no-spreading */}
      <IMaskInput
        className="usa-input"
        type="text"
        id={id}
        name={name}
        value={value ?? defaultValue}
        mask={mask}
        blocks={blocks}
        lazy={lazy}
        onAccept={(val, masked) => {
          helpers.setValue(masked.unmaskedValue);
          helpers.setTouched(true);
        }}
        {...props}
      />
      {/* eslint-enable react/jsx-props-no-spreading */}

      {!!warning && !hasError && (
        <p className="usa-hint" data-testid="textInputWarning">
          {warning}
        </p>
      )}
    </FormGroup>
  );
};

MaskedTextField.propTypes = {
  labelClassName: PropTypes.string,
  labelHint: PropTypes.string,
  id: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  label: PropTypes.string.isRequired,
  defaultValue: PropTypes.string,
  mask: PropTypes.string,
  blocks: PropTypes.oneOfType([PropTypes.object]),
  lazy: PropTypes.bool,
  warning: PropTypes.string,
};

MaskedTextField.defaultProps = {
  labelHint: '',
  labelClassName: '',
  defaultValue: '',
  mask: '',
  blocks: {},
  lazy: true, // make placeholder not visible
  warning: '',
};

export default MaskedTextField;
