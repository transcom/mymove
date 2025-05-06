import React from 'react';
import * as PropTypes from 'prop-types';
import 'styles/office.scss';
import { Link } from 'react-router-dom';

import styles from './ContactInfoDisplay.module.scss';

import descriptionListStyles from 'styles/descriptionList.module.scss';

const ContactInfoDisplay = ({ officeUserInfo, editURL }) => {
  return (
    <div className={styles.contactInfoContainer}>
      <div className={styles.contactInfoHeader}>
        <h2>Contact info</h2>
        <Link to={editURL}>Edit</Link>
      </div>

      <div className={styles.contactInfoSection}>
        <dl className={descriptionListStyles.descriptionList}>
          <div className={descriptionListStyles.row}>
            <dt>Name</dt>
            <dd data-testid="name">{officeUserInfo.name}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>Email</dt>
            <dd data-testid="email">{officeUserInfo.email}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>Phone</dt>
            <dd data-testid="phone">{officeUserInfo.telephone}</dd>
          </div>
        </dl>
      </div>
    </div>
  );
};

ContactInfoDisplay.propTypes = {
  officeUserInfo: PropTypes.shape({
    name: PropTypes.string,
    telephone: PropTypes.string,
    email: PropTypes.string,
  }).isRequired,
};

export default ContactInfoDisplay;
