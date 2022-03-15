import React from 'react';
import { string } from 'prop-types';

import styles from './ModifiedBy.module.scss';

const ModifiedBy = ({ firstName, lastName, email, phone }) => (
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

ModifiedBy.defaultProps = {
  firstName: null,
  lastName: null,
  email: null,
  phone: null,
};

ModifiedBy.propTypes = {
  firstName: string,
  lastName: string,
  email: string,
  phone: string,
};

export default ModifiedBy;
