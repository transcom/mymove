import React from 'react';
import PropTypes from 'prop-types';

import { useField } from 'formik';
import { FormGroup, Label, TextInput as UswdsTextInput } from '@trussworks/react-uswds';
import { ErrorMessage } from '..';

const idOrNameIsRequired = ({ id, name }, _, componentName) => {
  if (!id && !name) {
    return new Error(`id or name required on '${componentName}'`);
  }
  return null;
};

export const TextInput = ({ label, id, name, ...props }) => {
  /* eslint-disable react/jsx-props-no-spreading */
  const [, meta] = useField({ id, name, ...props });
  const hasError = meta.touched && !!meta.error;
  return (
    <FormGroup error={hasError}>
      <Label error={hasError} htmlFor={id || name}>
        {label}
      </Label>
      <TextInputMinimal id={id} name={name} {...props} />
    </FormGroup>
  );
  /* eslint-enable react/jsx-props-no-spreading */
};

TextInput.propTypes = {
  id: idOrNameIsRequired,
  name: idOrNameIsRequired,
  label: PropTypes.string.isRequired,
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
  id: idOrNameIsRequired,
  name: idOrNameIsRequired,
};
