import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';
import { Field, useField } from 'formik';
import { FormGroup, Checkbox, ErrorMessage } from '@trussworks/react-uswds';

import './CheckboxField.module.scss';

/**
 * This component renders a checkbox
 *
 * It relies on the Formik useField hook to work, so it must ALWAYS be rendered
 * inside of a Formik form context.
 *
 * If you want to use these components outside a Formik form, you can use the
 * ReactUSWDS components directly.
 */

export const CheckboxField = ({
  name,
  id,
  validate,
  type,
  warning,
  error,
  errorMessage,
  errorClassName,
  isDisabled,
  ...inputProps
}) => {
  const [fieldProps, metaProps, helperProps] = useField({ name, validate, type });
  const showError = (metaProps.touched && !!metaProps.error) || error;
  const showWarning = !showError && warning;

  const formGroupClasses = classnames({
    warning: showWarning,
  });

  // This immediately triggers state change for the yup validation errors
  // If this is not present and blur is not triggered, then only after a user clicks again
  // outside of the checkbox (blur) then it will trigger errors. We want to enforce
  // errors appearing immediately on click and prior to form submission.
  const handleClick = () => {
    helperProps.setValue(!metaProps.value);
    helperProps.setTouched(true);
  };

  return (
    <FormGroup className={formGroupClasses} error={showError}>
      {showError && (
        <ErrorMessage display={showError} className={errorClassName}>
          {metaProps.error ? metaProps.error : errorMessage}
        </ErrorMessage>
      )}
      <Field
        id={id}
        as={Checkbox}
        name={name}
        disabled={isDisabled}
        onClick={handleClick}
        onBlur={() => helperProps.setTouched(true)}
        /* eslint-disable-next-line react/jsx-props-no-spreading */
        {...fieldProps}
        /* eslint-disable-next-line react/jsx-props-no-spreading */
        {...inputProps}
      />
    </FormGroup>
  );
};

CheckboxField.propTypes = {
  id: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  warning: PropTypes.string,
  validate: PropTypes.oneOfType([PropTypes.string, PropTypes.func]),
  type: PropTypes.string,
  error: PropTypes.bool,
  errorMessage: PropTypes.string,
  errorClassName: PropTypes.string,
  isDisabled: PropTypes.bool,
};

CheckboxField.defaultProps = {
  warning: '',
  validate: undefined,
  type: 'checkbox',
  error: false,
  errorMessage: '',
  errorClassName: '',
  isDisabled: false,
};

export default CheckboxField;
