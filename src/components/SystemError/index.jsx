import React from 'react';
import { node } from 'prop-types';

import styles from './SystemError.module.scss';

const SystemError = ({ children }) => (
  <div className={styles.systemError} data-testid="system-error">
    {children}
  </div>
);

SystemError.propTypes = {
  children: node.isRequired,
};
export default SystemError;
