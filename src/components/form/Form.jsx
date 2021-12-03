import React from 'react';
import PropTypes from 'prop-types';
import { useFormikContext } from 'formik';
import { Form as UswdsForm } from '@trussworks/react-uswds';

export const Form = ({ errorCallback, ...props }) => {
  const { errors, touched, handleReset, handleSubmit } = useFormikContext();
  React.useEffect(() => {
    if (errorCallback) {
      errorCallback(errors, touched);
    }
  });
  // eslint-disable-next-line react/jsx-props-no-spreading
  return <UswdsForm onSubmit={handleSubmit} onReset={handleReset} role="form" {...props} />;
};

Form.defaultProps = {
  errorCallback: undefined,
};
Form.propTypes = {
  errorCallback: PropTypes.func,
};

export default Form;
