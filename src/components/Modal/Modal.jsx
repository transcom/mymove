/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';
import { Modal as USWDSModal } from '@trussworks/react-uswds';

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
