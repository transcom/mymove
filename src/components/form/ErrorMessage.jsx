import React from 'react';
import PropTypes from 'prop-types';
import { ErrorMessage as UswdsErrorMessage } from '@trussworks/react-uswds';

export const ErrorMessage = ({ display, children, ...props }) => {
  //  react/jsx-props-no-spreading
  return display && children ? <UswdsErrorMessage {...props}>{children}</UswdsErrorMessage> : null;
};

ErrorMessage.defaultProps = {
  children: null,
};

ErrorMessage.propTypes = {
  display: PropTypes.bool.isRequired,
  children: PropTypes.string,
};

export default ErrorMessage;
