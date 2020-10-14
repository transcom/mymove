/* eslint-disable react/jsx-props-no-spreading */
import React, { useEffect } from 'react';
import ReactDOM from 'react-dom';
import PropTypes from 'prop-types';
import classnames from 'classnames';
import { Modal as USWDSModal, connectModal as connectUSWDSModal, Button } from '@trussworks/react-uswds';

import styles from './Modal.module.scss';

import { ReactComponent as XLightIcon } from 'shared/icon/x-light.svg';

const Modal = ({ className, ...props }) => {
  const classes = classnames(styles.Modal, className);
  const APP_ROOT_ID = 'app-root';

  useEffect(() => {
    const appContainer = document.getElementById(APP_ROOT_ID);
    if (appContainer) {
      const scrollYPos = document.documentElement.scrollTop;
      appContainer.classList.add(styles.AppLocked);
      appContainer.style.transform = `translateY(-${scrollYPos}px)`;

      return () => {
        if (appContainer) {
          appContainer.classList.remove(styles.AppLocked);
          appContainer.style.transform = '';
          document.documentElement.scrollTo(0, scrollYPos);
        }
      };
    }

    return () => {};
  });

  return <USWDSModal className={classes} {...props} />;
};

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

export const ModalClose = ({ handleClick, className, ...buttonProps }) => (
  <Button
    type="button"
    onClick={handleClick}
    unstyled
    className={classnames(styles.ModalClose, className)}
    {...buttonProps}
  >
    <XLightIcon />
  </Button>
);

ModalClose.propTypes = {
  handleClick: PropTypes.func.isRequired,
  className: PropTypes.string,
};

ModalClose.defaultProps = {
  className: '',
};

export const ModalActions = ({ children }) => <div className={styles.ModalActions}>{children}</div>;

ModalActions.propTypes = {
  children: PropTypes.node.isRequired,
};

export const connectModal = (Component) => {
  return (props) => {
    // connectUSWDSModal handles isOpen prop & renders with container & overlay
    const ConnectedModal = connectUSWDSModal(Component);

    // Render into portal element if it exists
    const MODAL_ROOT_ID = 'modal-root';
    const modalContainer = document.getElementById(MODAL_ROOT_ID);
    if (modalContainer) {
      return ReactDOM.createPortal(<ConnectedModal {...props} />, modalContainer);
    }

    return <ConnectedModal {...props} />;
  };
};
