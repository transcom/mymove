/* eslint-disable react/jsx-props-no-spreading */
import React, { useEffect } from 'react';
import ReactDOM from 'react-dom';
import PropTypes from 'prop-types';
import classnames from 'classnames';
import { Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './Modal.module.scss';

import { Modal as MigratedModal, connectModal as connectMigratedModal } from 'components/MigratedModal/MigratedModal';

const Modal = ({ className, ...props }) => {
  const classes = classnames(styles.Modal, className);
  const APP_ROOT_ID = 'app-root';

  useEffect(() => {
    const appContainer = document.getElementById(APP_ROOT_ID);
    if (appContainer) {
      const scrollYPos = document.documentElement.scrollTop;
      appContainer.classList.add(styles.AppLocked);
      window.scrollTo(0, 0);

      return () => {
        if (appContainer) {
          appContainer.classList.remove(styles.AppLocked);
          window.scrollTo(0, scrollYPos);
        }
      };
    }

    return () => {};
  });

  return <MigratedModal className={classes} {...props} />;
};

Modal.displayName = 'MilMoveModal';

Modal.propTypes = {
  className: PropTypes.string,
};

Modal.defaultProps = {
  className: '',
};

export default Modal;

export const ModalTitle = ({ children }) => <div className={styles.ModalTitle}>{children}</div>;

ModalTitle.propTypes = {
  children: PropTypes.node.isRequired,
};

const ModalClose = ({ handleClick, className, ...buttonProps }) => {
  // What the button will display, defaults to X icon
  const { buttonContent, dataTestId } = buttonProps;

  return (
    <Button
      type="button"
      onClick={handleClick}
      unstyled
      className={classnames(styles.ModalClose, className)}
      data-testid="modalCloseButton"
      aria-label="Close"
      {...buttonProps}
    >
      {buttonContent || <FontAwesomeIcon icon="times" />}
    </Button>
  );
};

ModalClose.propTypes = {
  handleClick: PropTypes.func.isRequired,
  className: PropTypes.string,
};

ModalClose.defaultProps = {
  className: '',
};

export { ModalClose };

// TODO: default styling
const ModalSubmit = ({ handleClick, className, ...buttonProps }) => {
  // What the button will display, defaults to Submit
  // What testId we should use, defaults to modalSubmitButton
  const { buttonContent, dataTestId } = buttonProps;

  return (
    <Button
      className={classnames(styles.ModalSubmit, className)}
      type="button"
      onClick={handleClick}
      data-testid={dataTestId || 'modalSubmitButton'}
      aria-label="Submit"
      {...buttonProps}
    >
      {buttonContent || 'Submit'}
    </Button>
  );
};

ModalSubmit.propTypes = {
  handleClick: PropTypes.func.isRequired,
  className: PropTypes.string,
};

ModalSubmit.defaultProps = {
  className: '',
};

export { ModalSubmit };

export const ModalActions = ({ children }) => <div className={styles.ModalActions}>{children}</div>;

ModalActions.propTypes = {
  children: PropTypes.node.isRequired,
};

export const connectModal = (Component) => {
  function getDisplayName(WrappedComponent) {
    return WrappedComponent.displayName || WrappedComponent.name || 'Component';
  }

  const ConnectedModal = (props) => {
    // connectMigratedModal handles isOpen prop & renders with container & overlay
    const ConnectedMigratedModal = connectMigratedModal(Component);

    // Render into portal element if it exists
    const MODAL_ROOT_ID = 'modal-root';
    const modalContainer = document.getElementById(MODAL_ROOT_ID);
    if (modalContainer) {
      return ReactDOM.createPortal(<ConnectedMigratedModal {...props} />, modalContainer);
    }

    return <ConnectedMigratedModal {...props} />;
  };

  ConnectedModal.displayName = `Connected${getDisplayName(Component)}`;

  return ConnectedModal;
};
