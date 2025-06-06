import React from 'react';

import styles from './RequiredAsterisk.module.scss';

import Hint from 'components/Hint';

export const RequiredAsterisk = () => {
  return (
    <span data-testid="requiredAsterisk" className={styles.requiredAsterisk} aria-hidden="true">
      *
    </span>
  );
};

export const requiredAsteriskMessage = (
  <Hint data-testid="reqAsteriskMsg" id="reqAsteriskMsg">
    Fields marked with <RequiredAsterisk /> are required.
  </Hint>
);

export default RequiredAsterisk;
