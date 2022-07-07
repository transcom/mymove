import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';
import { useField } from 'formik';
import { FormGroup, Label, TextInput, ErrorMessage } from '@trussworks/react-uswds';

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
  ...inputProps
}) => {
  const [fieldProps, metaProps] = useField({ name, validate, type });
  const showError = (metaProps.touched && !!metaProps.error) || error;
  const showWarning = !showError && warning;

  const formGroupClasses = classnames({
    warning: showWarning,
  });

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
      {/* eslint-disable-next-line react/jsx-props-no-spreading */}
      <TextInput id={id} name={name} disabled={isDisabled} {...fieldProps} {...inputProps} />
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
  type: PropTypes.string,
  error: PropTypes.bool,
  errorMessage: PropTypes.string,
  errorClassName: PropTypes.string,
  isDisabled: PropTypes.bool,
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
};

export default TextField;
