import React from 'react';
import { string } from 'prop-types';

import styles from './ModifiedBy.module.scss';

const milMoveName = 'MilMove';
const primeName = 'Prime';

const UserModifiedBy = ({ firstName, lastName, email, phone }) => {
  return (
    <div className={styles.ModifiedBy}>
      <span className={styles.name}>
        {lastName}, {firstName}
      </span>
      {email && phone && (
        <div className={styles.contactInfo}>
          {email} | {phone}
        </div>
      )}
    </div>
  );
};

// If an event is modified by the MilMove or Prime systems, it will contain limited
// or no user information. This is used to identify moves that were modified by those
// systems.
const SystemModifiedBy = ({ systemName }) => {
  return (
    <div className={styles.ModifiedBy}>
      <span className={styles.name}>{systemName}</span>
    </div>
  );
};

const ModifiedBy = ({ firstName, lastName, email, phone }) => {
  const isUserEmpty = firstName === '' && lastName === '' && email === '' && phone === '';
  const isUserPrime = firstName === primeName && lastName === '' && email === '' && phone === '';
  if (isUserEmpty) {
    return <SystemModifiedBy systemName={milMoveName} />;
  }
  if (isUserPrime) {
    return <SystemModifiedBy systemName={primeName} />;
  }
  return (
    <div className={styles.ModifiedBy}>
      <span className={styles.name}>
        {lastName}, {firstName}
      </span>
      {email && phone && (
        <div className={styles.contactInfo}>
          {email} | {phone}
        </div>
      )}
    </div>
  );
};

UserModifiedBy.propTypes = {
  firstName: string.isRequired,
  lastName: string.isRequired,
  email: string.isRequired,
  phone: string.isRequired,
};

SystemModifiedBy.propTypes = {
  systemName: string.isRequired,
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
