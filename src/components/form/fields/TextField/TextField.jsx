/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';
import { useField } from 'formik';
import { FormGroup, Label, TextInput, Textarea, ErrorMessage } from '@trussworks/react-uswds';

import { OptionalTag } from 'components/form/OptionalTag';
import Hint from 'components/Hint';

/**
 * This component renders a ReactUSWDS TextInput component inside of a FormGroup,
 * with a Label and ErrorMessage.
 *
 * It relies on the Formik useField hook to work, so it must ALWAYS be rendered
 * inside of a Formik form context.
 *
 * If you want to use these components outside a Formik form, you can use the
 * ReactUSWDS components directly.
 */

const TextField = ({
  name,
  id,
  label,
  labelClassName,
  labelHint,
  validate,
  type,
  optional,
  warning,
  error,
  errorMessage,
  errorClassName,
  isDisabled,
  display,
  button,
  ...inputProps
}) => {
  const [fieldProps, metaProps] = useField({ name, validate, type });
  const showError = (metaProps.touched && !!metaProps.error) || error;
  const showWarning = !showError && warning;

  const formGroupClasses = classnames({
    warning: showWarning,
  });

  const getDisplay = (displayType) => {
    switch (displayType) {
      case 'textarea':
        return <Textarea id={id} name={name} disabled={isDisabled} {...fieldProps} {...inputProps} />;
      case 'readonly':
        return <p data-testid="textFieldReadOnlyValue">{fieldProps.value}</p>;
      default:
        return <TextInput id={id} name={name} disabled={isDisabled} {...fieldProps} {...inputProps} />;
    }
  };

  return (
    <FormGroup className={formGroupClasses} error={showError}>
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
      {getDisplay(display)}

      {button || null}
    </FormGroup>
  );
};

TextField.propTypes = {
  id: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  label: PropTypes.node.isRequired,
  labelClassName: PropTypes.string,
  labelHint: PropTypes.string,
  warning: PropTypes.string,
  optional: PropTypes.bool,
  validate: PropTypes.oneOfType([PropTypes.string, PropTypes.func]),
  display: PropTypes.oneOf(['input', 'textarea', 'readonly']),
  type: PropTypes.string,
  error: PropTypes.bool,
  errorMessage: PropTypes.string,
  errorClassName: PropTypes.string,
  isDisabled: PropTypes.bool,
  button: PropTypes.node,
};

TextField.defaultProps = {
  labelHint: '',
  labelClassName: '',
  warning: '',
  optional: false,
  validate: undefined,
  type: 'text',
  error: false,
  errorMessage: '',
  errorClassName: '',
  isDisabled: false,
  display: 'input',
  button: undefined,
};

export default TextField;
