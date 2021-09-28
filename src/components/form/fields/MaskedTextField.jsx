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
  inputClassName,
  errorClassName,
  id,
  name,
  labelHint,
  defaultValue,
  mask,
  blocks,
  lazy,
  warning,
  validate,
  children,
  type,
  ...props
}) => {
  const [field, meta, helpers] = useField({ id, name, validate, ...props });
  const hasError = meta.touched && !!meta.error;
  const { value } = field;
  return (
    <FormGroup className={classnames(!!warning && !hasError && `warning`, formGroupClassName)} error={hasError}>
      <Label className={labelClassName} hint={labelHint} error={hasError} htmlFor={id || name}>
        {label}
      </Label>
      <ErrorMessage display={hasError} className={errorClassName}>
        {meta.error}
      </ErrorMessage>
      {!!warning && !hasError && (
        <p className="usa-hint" data-testid="textInputWarning">
          {warning}
        </p>
      )}
      {/* eslint-disable react/jsx-props-no-spreading */}
      <IMaskInput
        className={classnames('usa-input', inputClassName)}
        type={type}
        id={id}
        name={name}
        value={value ?? defaultValue}
        mask={mask}
        blocks={blocks}
        lazy={lazy}
        onAccept={(val, masked) => {
          helpers.setValue(masked.unmaskedValue);
          // setValue is already triggering validation for this field so we should be able to skip it in setTouched
          helpers.setTouched(true, false);
        }}
        {...props}
      />
      {children}
      {/* eslint-enable react/jsx-props-no-spreading */}
    </FormGroup>
  );
};

MaskedTextField.propTypes = {
  blocks: PropTypes.oneOfType([PropTypes.object]),
  children: PropTypes.node,
  defaultValue: PropTypes.string,
  errorClassName: PropTypes.string,
  formGroupClassName: PropTypes.string,
  id: PropTypes.string.isRequired,
  inputClassName: PropTypes.string,
  label: PropTypes.string,
  labelClassName: PropTypes.string,
  labelHint: PropTypes.string,
  lazy: PropTypes.bool,
  mask: PropTypes.oneOfType([PropTypes.string, PropTypes.func]),
  name: PropTypes.string.isRequired,
  type: PropTypes.string,
  validate: PropTypes.func,
  warning: PropTypes.string,
};

MaskedTextField.defaultProps = {
  blocks: {},
  children: null,
  defaultValue: '',
  errorClassName: '',
  formGroupClassName: '',
  inputClassName: '',
  label: '',
  labelClassName: '',
  labelHint: '',
  lazy: true, // make placeholder not visible
  mask: '',
  type: 'text',
  validate: undefined,
  warning: '',
};

export default MaskedTextField;
