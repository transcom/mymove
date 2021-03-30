import React, { useEffect } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { Alert } from '@trussworks/react-uswds';

import { clearFlashMessage as clearFlashMessageAction } from 'store/flash/actions';
import { FlashMessageShape } from 'types/flash';

export const FlashMessage = ({ flash, clearFlashMessage }) => {
  useEffect(() => () => {
    // Clear this flash message on unmount (this will happen on navigation or if flash state changes)
    clearFlashMessage(flash?.key);
  });
  const { message, title, type, slim } = flash;

  return (
    // We use {title || undefined} here because an empty string as the title will render a blank header in Firefox,
    // so we must pass in undefined if we want to see no header at all.
    <Alert slim={slim} type={type} heading={title || undefined}>
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
      return <Component key={showFlash} {...props} />;
    }

    return null;
  };

  ConnectedFlashMessage.displayName = 'ConnectedFlashMessage';

  ConnectedFlashMessage.propTypes = {
    flash: PropTypes.shape({
      type: PropTypes.string,
      title: PropTypes.string,
      message: PropTypes.string,
      key: PropTypes.string,
      slim: PropTypes.bool,
    }),
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
