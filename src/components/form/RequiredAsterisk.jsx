import React from 'react';

import styles from './RequiredAsterisk.module.scss';

export const RequiredAsterisk = () => {
  return (
    <span data-testid="requiredAsterisk" className={styles.requiredAsterisk}>
      *
    </span>
  );
};

export const requiredAsteriskMessage = (
  <div data-testid="reqAsteriskMsg">
    Fields marked with <span className={styles.requiredAsterisk}>*</span> are required.
  </div>
);

export default RequiredAsterisk;
