import React, { useEffect, useState } from 'react';
import * as PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { GridContainer } from '@trussworks/react-uswds';

import ScrollToTop from 'components/ScrollToTop';
import ConnectedFlashMessage from 'containers/FlashMessage/FlashMessage';
import { clearFlashMessage } from 'store/flash/actions';

const FlashGridContainer = ({ children, messageKey, clearMessage, ...gridContainerProps }) => {
  // This reference keeps track of new flash messages, scrolls us up to the top of the page if a new message is added,
  // and adds an event listener to clear the message out when it is no longer needed.
  //   NOTE: We should probably use a hook like `useRef` here instead (this isn't really state),
  //   but we do want to trigger a re-render when it changes.
  //   Otherwise, it introduces a delay with the ScrollToTop component.
  const [currentMessageKey, setCurrentMessageKey] = useState('');
  useEffect(() => {
    const clearCurrentMessage = () => {
      clearMessage(messageKey);
    };

    if (messageKey) {
      setCurrentMessageKey(messageKey);

      // clears the current message after another element on the same page gets focused (not on window blur):
      window.addEventListener('focusout', clearCurrentMessage);
    }

    // remove the listener on component unmount:
    return () => {
      window.removeEventListener('focusout', clearCurrentMessage);
    };
  }, [messageKey, setCurrentMessageKey, clearMessage]);

  return (
    // eslint-disable-next-line react/jsx-props-no-spreading
    <GridContainer {...gridContainerProps}>
      <ScrollToTop otherDep={currentMessageKey} />
      <ConnectedFlashMessage />
      {children}
    </GridContainer>
  );
};

FlashGridContainer.propTypes = {
  children: PropTypes.node.isRequired,
  messageKey: PropTypes.string.isRequired,
  clearMessage: PropTypes.func.isRequired,
};

const mapStateToProps = (state) => {
  return {
    messageKey: state.flash.flashMessage?.key || '',
  };
};

const mapDispatchToProps = {
  clearMessage: clearFlashMessage,
};

export default connect(mapStateToProps, mapDispatchToProps)(FlashGridContainer);
