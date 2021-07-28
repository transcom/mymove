import React from 'react';
import { node } from 'prop-types';

import styles from './SystemError.module.scss';

const SystemError = ({ children }) => <div className={styles.systemError}>{children}</div>;

SystemError.propTypes = {
  children: node.isRequired,
};
export default SystemError;
