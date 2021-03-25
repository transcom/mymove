import React from 'react';
import PropTypes from 'prop-types';
import { useFormikContext } from 'formik';
import { Form as UswdsForm } from '@trussworks/react-uswds';

export const Form = ({ errorCallback, ...props }) => {
  const { errors, touched, handleReset, handleSubmit } = useFormikContext();
  React.useEffect(() => {
    if (errorCallback) {
      errorCallback(errors, touched);
      // RA Summary: eslint: no-console - System Information Leak: External
      // RA: The linter flags any use of console.
      // RA: This console serves to indicate the status of the swaggerRequest to a user for debugging purposes.
      // RA: Given that this is a simple string with no interpolation
      // RA: nor variable names, SQL strings, system path information, or source or program code,
      // RA: this is not a finding.
      // RA Developer Status: Mitigated
      // RA Validator Status: {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}
      // RA Validator: jneuner@mitre.org
      // RA Modified Severity: CAT III
      // eslint-disable-next-line
      console.log(`test`);
    }
  });
  // eslint-disable-next-line react/jsx-props-no-spreading
  return <UswdsForm onSubmit={handleSubmit} onReset={handleReset} {...props} />;
};

Form.defaultProps = {
  errorCallback: undefined,
};
Form.propTypes = {
  errorCallback: PropTypes.func,
};

export default Form;
