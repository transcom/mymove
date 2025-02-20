import React from 'react';
import PropTypes from 'prop-types';
import { Oval } from 'react-loader-spinner';

import styles from './LoadingSpinner.module.scss';

const LoadingSpinner = ({ message }) => (
  <div className={styles.container} data-testid="loading-spinner">
    <div className={styles.spinnerWrapper}>
      <Oval visible height="150" width="150" color="#ffbe2e" secondaryColor="#565c65" ariaLabel="oval-loading" />
      <p className={styles.message}>{message || 'Loading, please wait...'}</p>
    </div>
  </div>
);

LoadingSpinner.propTypes = {
  message: PropTypes.string,
};

LoadingSpinner.defaultProps = {
  message: '',
};

export default LoadingSpinner;
