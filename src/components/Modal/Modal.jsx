/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import ReactDOM from 'react-dom';
import PropTypes from 'prop-types';
import classnames from 'classnames';
import { Modal as USWDSModal, connectModal as connectUSWDSModal } from '@trussworks/react-uswds';

import styles from './Modal.module.scss';

const Modal = ({ className, ...props }) => {
  const classes = classnames(styles.Modal, className);
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

export const ModalActions = ({ children }) => <div className={styles.ModalActions}>{children}</div>;

ModalActions.propTypes = {
  children: PropTypes.node.isRequired,
};

export const connectModal = (Component) => {
  const ConnectedModal = connectUSWDSModal(Component);
  const MODAL_ROOT_ID = 'modal-root';
  // Render into portal element if it exists
  const modalContainer = document.getElementById(MODAL_ROOT_ID);
  if (modalContainer) {
    return ReactDOM.createPortal(ConnectedModal, modalContainer);
  }

  return ConnectedModal;
};
