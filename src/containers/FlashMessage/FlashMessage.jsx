import React, { useEffect } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import { clearFlashMessage as clearFlashMessageAction } from 'store/flash/actions';
import Alert from 'shared/Alert';

export const FlashMessage = ({ children, flash, clearFlashMessage }) => {
  useEffect(() => () => {
    // Clear flash message on unmount (this will happen on navigation or if flash state changes)
    clearFlashMessage();
  });

  const { message, title, type, key } = flash;

  // display children if passed in as a custom flash (using the key)
  if (key && children) {
    return children;
  }

  // otherwise fall back to alert
  // message & type are required
  return (
    <Alert type={type} heading={title}>
      {message}
    </Alert>
  );
};

FlashMessage.displayName = 'FlashMessage';

FlashMessage.propTypes = {
  children: PropTypes.node,
  flash: PropTypes.shape({
    type: PropTypes.string,
    title: PropTypes.string,
    message: PropTypes.string,
    key: PropTypes.string.isRequired,
  }).isRequired,
  clearFlashMessage: PropTypes.func.isRequired,
};

FlashMessage.defaultProps = {
  children: null,
};

const connectFlashMessage = (Component) => {
  const ConnectedFlashMessage = (props) => {
    const { flash } = props;

    // Only render flash if a key is defined
    const showFlash = flash?.key;

    if (showFlash) {
      // eslint-disable-next-line react/jsx-props-no-spreading
      return <Component {...props} />;
    }

    return null;
  };

  ConnectedFlashMessage.displayName = 'ConnectedFlashMessage';

  ConnectedFlashMessage.propTypes = {
    children: PropTypes.node,
    flash: PropTypes.shape({
      type: PropTypes.string,
      title: PropTypes.string,
      message: PropTypes.string,
      key: PropTypes.string,
    }),
    clearFlashMessage: PropTypes.func.isRequired,
  };

  ConnectedFlashMessage.defaultProps = {
    children: null,
    flash: null,
  };

  const mapDispatchToProps = {
    clearFlashMessage: clearFlashMessageAction,
  };

  const mapStateToProps = (state) => ({
    flash: state.flash.flashMessage,
  });

  return connect(mapStateToProps, mapDispatchToProps)(ConnectedFlashMessage);
};

export default connectFlashMessage(FlashMessage);
