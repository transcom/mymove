import React from 'react';

import styles from './RequiredAsterisk.module.scss';

import Hint from 'components/Hint';

export const RequiredAsterisk = () => {
  return (
    <span data-testid="requiredAsterisk" className={styles.requiredAsterisk}>
      *
    </span>
  );
};

export const requiredAsteriskMessage = (
  <Hint data-testid="reqAsteriskMsg">
    Fields marked with <RequiredAsterisk /> are required.
  </Hint>
);

export default RequiredAsterisk;
