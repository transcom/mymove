import React from 'react';
import PropTypes from 'prop-types';
import { useField } from 'formik';
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
