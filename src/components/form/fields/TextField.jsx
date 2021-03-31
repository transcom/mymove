import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';
import { useField } from 'formik';
import { FormGroup, Label, TextInput, ErrorMessage } from '@trussworks/react-uswds';

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

const TextField = ({ name, id, label, labelClassName, labelHint, validate, type, warning, ...inputProps }) => {
  const [fieldProps, metaProps] = useField({ name, validate, type });
  const showError = metaProps.touched && !!metaProps.error;
  const showWarning = !showError && warning;

  const formGroupClasses = classnames({
    warning: showWarning,
  });

  return (
    <FormGroup className={formGroupClasses} error={showError}>
      <Label className={labelClassName} hint={labelHint} error={showError} htmlFor={id || name}>
        {label}
      </Label>

      {showError && <ErrorMessage>{metaProps.error}</ErrorMessage>}

      {/* eslint-disable-next-line react/jsx-props-no-spreading */}
      <TextInput id={id} name={name} {...fieldProps} {...inputProps} />

      {showWarning && (
        <p className="usa-hint" data-testid="textInputWarning">
          {warning}
        </p>
      )}
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
  validate: PropTypes.oneOfType([PropTypes.string, PropTypes.func]),
  type: PropTypes.string,
};

TextField.defaultProps = {
  labelHint: '',
  labelClassName: '',
  warning: '',
  validate: undefined,
  type: 'text',
};

export default TextField;
