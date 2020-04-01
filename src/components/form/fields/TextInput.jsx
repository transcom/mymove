import React from 'react';
import PropTypes from 'prop-types';

import { useField } from 'formik';
import { FormGroup, Label, TextInput as UswdsTextInput } from '@trussworks/react-uswds';
import { ErrorMessage } from '..';

export const TextInput = ({ label, id, name, ...props }) => {
  /* eslint-disable react/jsx-props-no-spreading */
  const [field, meta] = useField({ id, name, ...props });
  const hasError = meta.touched && !!meta.error;
  return id || name ? (
    <FormGroup error={hasError}>
      <Label error={hasError} htmlFor={id || name}>
        {label}
      </Label>
      <ErrorMessage display={hasError}>{meta.error}</ErrorMessage>
      <UswdsTextInput id={id} name={name} {...field} {...props} />
    </FormGroup>
  ) : null;
  /* eslint-enable react/jsx-props-no-spreading */
};

TextInput.propTypes = {
  id: ({ id, name }, _, componentName) => {
    if (!id && !name) {
      return new Error(`id or name required on '${componentName}'`);
    }
    return null;
  },
  name: ({ id, name }, _, componentName) => {
    if (!id && !name) {
      return new Error(`id or name required on '${componentName}'`);
    }
    return null;
  },
  label: PropTypes.string.isRequired,
};

export default TextInput;
