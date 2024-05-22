import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './LockedMoveBanner.module.scss';

const LockedMoveBanner = ({ children }) => (
  <div className={styles.lockedMoveBanner} data-testid="locked-move-banner">
    <FontAwesomeIcon icon="lock" /> {children}
  </div>
);

export default LockedMoveBanner;
