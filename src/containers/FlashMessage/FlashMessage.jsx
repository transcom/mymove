import React, { useEffect } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import { clearFlashMessage as clearFlashMessageAction } from 'store/flash/actions';
import Alert from 'shared/Alert';
import { FlashMessageShape } from 'types/flash';

export const FlashMessage = ({ flash, clearFlashMessage }) => {
  useEffect(() => () => {
    // Clear flash message on unmount (this will happen on navigation or if flash state changes)
    clearFlashMessage();
  });

  const { message, title, type } = flash;

  return (
    <Alert type={type} heading={title}>
      {message}
    </Alert>
  );
};

FlashMessage.displayName = 'FlashMessage';

FlashMessage.propTypes = {
  flash: FlashMessageShape.isRequired,
  clearFlashMessage: PropTypes.func.isRequired,
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
    flash: FlashMessageShape,
    clearFlashMessage: PropTypes.func.isRequired,
  };

  ConnectedFlashMessage.defaultProps = {
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
