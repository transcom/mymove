import React from 'react';
import classnames from 'classnames';
import PropTypes from 'prop-types';
import { useField } from 'formik';
import { IMaskInput } from 'react-imask';
import { FormGroup, Label } from '@trussworks/react-uswds';

import styles from './MaskedTextField.module.scss';

import { OptionalTag } from 'components/form/OptionalTag';
import { ErrorMessage } from 'components/form/index';
import Hint from 'components/Hint';

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
  inputTestId,
  optional,
  errorMessage,
  error,
  suffix,
  ...props
}) => {
  const [field, metaProps, helpers] = useField({ id, name, validate, ...props });
  const showError = (metaProps.touched && !!metaProps.error) || error;
  const showWarning = !showError && warning;
  const { value } = field;
  return (
    <FormGroup className={classnames(!!warning && !showError && `warning`, formGroupClassName)} error={showError}>
      <div className="labelWrapper">
        <Label className={labelClassName} hint={labelHint} error={showError} htmlFor={id || name}>
          {label}
        </Label>
        {optional && <OptionalTag />}
      </div>
      {showError && (
        <ErrorMessage display={showError} className={errorClassName}>
          {metaProps.error ? metaProps.error : errorMessage}
        </ErrorMessage>
      )}
      {showWarning && <Hint data-testid="textInputWarning">{warning}</Hint>}
      {/* eslint-disable react/jsx-props-no-spreading */}
      {suffix ? (
        <div className={suffix && styles.hasSuffix}>
          <IMaskInput
            className={classnames('usa-input', inputClassName)}
            data-testid={inputTestId}
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
            {...field}
          />
          {suffix && <div className="suffix">{suffix}</div>}
        </div>
      ) : (
        <IMaskInput
          className={classnames('usa-input', inputClassName)}
          data-testid={inputTestId}
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
          {...field}
        />
      )}
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
  suffix: PropTypes.string,
  label: PropTypes.string,
  labelClassName: PropTypes.string,
  labelHint: PropTypes.string,
  lazy: PropTypes.bool,
  mask: PropTypes.oneOfType([PropTypes.string, PropTypes.func]),
  name: PropTypes.string.isRequired,
  type: PropTypes.string,
  validate: PropTypes.func,
  warning: PropTypes.string,
  inputTestId: PropTypes.string,
  optional: PropTypes.bool,
  error: PropTypes.bool,
  errorMessage: PropTypes.string,
};

MaskedTextField.defaultProps = {
  blocks: {},
  children: null,
  defaultValue: '',
  errorClassName: '',
  formGroupClassName: '',
  suffix: '',
  inputClassName: '',
  label: '',
  labelClassName: '',
  labelHint: '',
  lazy: true, // make placeholder not visible
  mask: '',
  type: 'text',
  validate: undefined,
  warning: '',
  inputTestId: '',
  optional: false,
  error: false,
  errorMessage: '',
};

export default MaskedTextField;
