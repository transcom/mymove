/* eslint-disable react/jsx-props-no-spreading */
import React, { useState } from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';

import styles from './MigratedModal.module.scss';

/** This is a straightforward port of the Modal component from React-USWDS 1.17
 *  into the MilMove project, as the component is being deprecated in USWDS 2.x. */

/** Modal UI component */
export const Modal = ({ title, children, actions, className }) => {
  const classes = classnames(styles.modal, className);

  return (
    <div data-testid="modal" className={classes}>
      <div className={styles.title}>{title}</div>
      <div className={styles.content}>{children}</div>
      <div className={styles.actions}>{actions}</div>
    </div>
  );
};

Modal.propTypes = {
  title: PropTypes.node,
  children: PropTypes.node.isRequired,
  actions: PropTypes.node,
  className: PropTypes.string,
};

Modal.defaultProps = {
  title: '',
  actions: '',
  className: '',
};

/** Overlay UI component (grey background) */
export const Overlay = () => <div className={styles.overlay} />;

/** Modal positioning component */
export const ModalContainer = ({ children }) => {
  return <div className={styles.modalContainer}>{children}</div>;
};

ModalContainer.propTypes = {
  children: PropTypes.node.isRequired,
};

export const connectModal = (Component) => {
  const ConnectedModal = ({ isOpen, ...props }) => {
    if (!isOpen) return null;
    return (
      <>
        <Overlay />
        <ModalContainer>
          <Component {...props} />
        </ModalContainer>
      </>
    );
  };

  ConnectedModal.propTypes = {
    isOpen: PropTypes.bool,
  };

  ConnectedModal.defaultProps = {
    isOpen: false,
  };

  return ConnectedModal;
};

export const useModal = () => {
  const [isOpen, setIsOpen] = useState(false);

  const openModal = () => {
    setIsOpen(true);
  };
  const closeModal = () => {
    setIsOpen(false);
  };

  return { isOpen, openModal, closeModal };
};
