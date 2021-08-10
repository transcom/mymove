import React from 'react';
import classnames from 'classnames';
import PropTypes from 'prop-types';
import { useField } from 'formik';
import { IMaskInput } from 'react-imask';
import { FormGroup, Label } from '@trussworks/react-uswds';

import { ErrorMessage } from '../index';

const MaskedTextField = ({
  label,
  labelClassName,
  formGroupClassName,
  id,
  name,
  labelHint,
  defaultValue,
  mask,
  blocks,
  lazy,
  warning,
  onAccept,
  ...props
}) => {
  const [field, meta, helpers] = useField({ id, name, ...props });
  const hasError = meta.touched && !!meta.error;
  const { value } = field;

  let handleAccept = (val, masked) => {
    helpers.setValue(masked.unmaskedValue);
    helpers.setTouched(true);
  };

  if (onAccept) {
    handleAccept = (val, masked) => {
      onAccept(helpers, val, masked);
    };
  }

  return (
    <FormGroup className={classnames(!!warning && !hasError && `warning`, formGroupClassName)} error={hasError}>
      <Label className={labelClassName} hint={labelHint} error={hasError} htmlFor={id || name}>
        {label}
      </Label>
      <ErrorMessage display={hasError}>{meta.error}</ErrorMessage>
      {!!warning && !hasError && (
        <p className="usa-hint" data-testid="textInputWarning">
          {warning}
        </p>
      )}

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
        onAccept={handleAccept}
        {...props}
      />
      {/* eslint-enable react/jsx-props-no-spreading */}
    </FormGroup>
  );
};

MaskedTextField.propTypes = {
  formGroupClassName: PropTypes.string,
  labelClassName: PropTypes.string,
  labelHint: PropTypes.string,
  id: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  label: PropTypes.string.isRequired,
  defaultValue: PropTypes.string,
  mask: PropTypes.oneOfType([PropTypes.string, PropTypes.func]),
  blocks: PropTypes.oneOfType([PropTypes.object]),
  lazy: PropTypes.bool,
  warning: PropTypes.string,
  onAccept: PropTypes.func,
};

MaskedTextField.defaultProps = {
  labelHint: '',
  labelClassName: '',
  formGroupClassName: '',
  defaultValue: '',
  mask: '',
  blocks: {},
  lazy: true, // make placeholder not visible
  warning: '',
  onAccept: undefined,
};

export default MaskedTextField;
