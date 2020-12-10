import React from 'react';
import PropTypes from 'prop-types';
import { useField } from 'formik';
import { IMaskInput } from 'react-imask';
import { FormGroup, Label, TextInput as UswdsTextInput } from '@trussworks/react-uswds';

import { ErrorMessage } from '../index';

export const TextInput = ({ label, labelClassName, id, name, labelHint, ...props }) => {
  /* eslint-disable react/jsx-props-no-spreading */
  const [, meta] = useField({ id, name, ...props });
  const hasError = meta.touched && !!meta.error;
  return (
    <FormGroup error={hasError}>
      <Label className={labelClassName} hint={labelHint} error={hasError} htmlFor={id || name}>
        {label}
      </Label>
      <TextInputMinimal id={id} name={name} {...props} />
    </FormGroup>
  );
  /* eslint-enable react/jsx-props-no-spreading */
};

TextInput.propTypes = {
  labelClassName: PropTypes.string,
  labelHint: PropTypes.string,
  id: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  label: PropTypes.string.isRequired,
};

TextInput.defaultProps = {
  labelHint: '',
  labelClassName: '',
};

export default TextInput;

export const TextInputMinimal = ({ id, name, ...props }) => {
  /* eslint-disable react/jsx-props-no-spreading */
  const [field, meta] = useField({ id, name, ...props });
  const hasError = meta.touched && !!meta.error;
  return (
    <>
      <ErrorMessage display={hasError}>{meta.error}</ErrorMessage>
      <UswdsTextInput id={id} name={name} {...field} {...props} />
    </>
  );
  /* eslint-enable react/jsx-props-no-spreading */
};

TextInputMinimal.propTypes = {
  id: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
};

export const TextMaskedInput = ({
  label,
  labelClassName,
  id,
  name,
  labelHint,
  defaultValue,
  mask,
  blocks,
  lazy,
  ...props
}) => {
  /* eslint-disable react/jsx-props-no-spreading */
  const [field, meta, helpers] = useField({ id, name, ...props });
  const hasError = meta.touched && !!meta.error;
  const { value } = field;
  return (
    <FormGroup error={hasError}>
      <Label className={labelClassName} hint={labelHint} error={hasError} htmlFor={id || name}>
        {label}
      </Label>
      <ErrorMessage display={hasError}>{meta.error}</ErrorMessage>
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
    </FormGroup>
  );
};

TextMaskedInput.propTypes = {
  labelClassName: PropTypes.string,
  labelHint: PropTypes.string,
  id: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  label: PropTypes.string.isRequired,
  defaultValue: PropTypes.string,
  mask: PropTypes.string,
  blocks: PropTypes.oneOfType([PropTypes.object]),
  lazy: PropTypes.bool,
};

TextMaskedInput.defaultProps = {
  labelHint: '',
  labelClassName: '',
  defaultValue: '',
  mask: '',
  blocks: {},
  lazy: true, // make placeholder not visible
};
