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
