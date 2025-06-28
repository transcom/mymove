/* eslint-disable react/jsx-props-no-spreading */
import React, { useEffect, useRef, useState } from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';
import { useField } from 'formik';
import { FormGroup, Label, TextInput, Textarea, ErrorMessage } from '@trussworks/react-uswds';

import styles from './TextField.module.scss';

import Hint from 'components/Hint';
import RequiredAsterisk from 'components/form/RequiredAsterisk';

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
  disablePaste,
  showRequiredAsterisk,
  prefix,
  ...inputProps
}) => {
  const [fieldProps, metaProps] = useField({ name, validate, type });
  const showError = (metaProps.touched && !!metaProps.error) || error;
  const showWarning = !showError && warning;

  const formGroupClasses = classnames({
    warning: showWarning,
  });

  const pasteHandler = disablePaste ? (e) => e.preventDefault() : undefined;

  const prefixRef = useRef(null);
  const [prefixWidth, setPrefixWidth] = useState(0);

  useEffect(() => {
    if (prefixRef.current) {
      setPrefixWidth(prefixRef.current.offsetWidth + 16);
    }
  }, [prefix]);

  const getDisplay = (displayType) => {
    switch (displayType) {
      case 'textarea':
        return (
          <Textarea
            id={id}
            name={name}
            disabled={isDisabled}
            onPaste={pasteHandler}
            {...fieldProps}
            {...inputProps}
            aria-describedby={showError ? `${id}-error` : undefined}
          />
        );
      case 'readonly':
        return (
          <label htmlFor={id || name} id={id} data-testid={label} aria-label={name}>
            {fieldProps.value}
          </label>
        );
      default:
        if (!prefix) {
          return (
            <TextInput
              id={id}
              name={name}
              disabled={isDisabled}
              onPaste={pasteHandler}
              {...fieldProps}
              {...inputProps}
              aria-describedby={showError ? `${id}-error` : undefined}
            />
          );
        }

        return (
          <div className={styles.inputWithPrefix}>
            <span ref={prefixRef} className={styles.prefix}>
              {prefix}
            </span>
            <TextInput
              id={id}
              name={name}
              disabled={isDisabled}
              onPaste={pasteHandler}
              {...fieldProps}
              {...inputProps}
              aria-describedby={showError ? `${id}-error` : undefined}
              className={classnames(inputProps.className, styles.prefixedInput)}
              style={{
                paddingLeft: `${prefixWidth}px`,
                ...inputProps.style,
              }}
            />
          </div>
        );
    }
  };

  return (
    <FormGroup className={formGroupClasses} error={showError}>
      <div className="labelWrapper">
        <Label className={labelClassName} hint={labelHint} error={showError} htmlFor={id || name}>
          <span>
            {label} {showRequiredAsterisk && <RequiredAsterisk />}
          </span>
        </Label>
        {optional}
      </div>

      {showError && (
        <ErrorMessage id={`${id}-error`} role="alert" aria-live="assertive" className={errorClassName}>
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
  disablePaste: PropTypes.bool,
  prefix: PropTypes.string,
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
  disablePaste: false,
  prefix: '',
};

export default TextField;
