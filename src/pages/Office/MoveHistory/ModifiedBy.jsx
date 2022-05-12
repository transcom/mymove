import React from 'react';
import { string } from 'prop-types';

import styles from './ModifiedBy.module.scss';

const ModifiedBy = ({ firstName, lastName, email, phone }) => {
  // If an event is modified by the MilMove system, it will contain no
  // information. This is used to idetify moves that were modified by the
  // MilMove system itself.
  const isUserEmpty = firstName === '' && lastName === '' && email === '' && phone === '';
  if (isUserEmpty) {
    const systemName = 'MilMove';
    return (
      <div className={styles.ModifiedBy}>
        <span className={styles.name}>{systemName}</span>
      </div>
    );
  }

  return (
    <div className={styles.ModifiedBy}>
      {lastName && firstName && (
        <span className={styles.name}>
          {lastName}, {firstName}
        </span>
      )}
      {email && phone && (
        <div className={styles.contactInfo}>
          {email} | {phone}
        </div>
      )}
    </div>
  );
};

ModifiedBy.defaultProps = {
  firstName: '',
  lastName: '',
  email: '',
  phone: '',
};

ModifiedBy.propTypes = {
  firstName: string,
  lastName: string,
  email: string,
  phone: string,
};

export default ModifiedBy;
