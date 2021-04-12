import React from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router-dom';

import styles from './ContactInfoDisplay.module.scss';

import { ResidentialAddressShape } from 'types/address';
import { BackupContactShape } from 'types/customerShapes';
import descriptionListStyles from 'styles/descriptionList.module.scss';

const ContactInfoDisplay = ({
  telephone,
  secondaryTelephone,
  personalEmail,
  phoneIsPreferred,
  emailIsPreferred,
  residentialAddress,
  backupMailingAddress,
  backupContact,
  editURL,
}) => {
  let preferredContactMethod = 'Unknown';
  if (phoneIsPreferred && emailIsPreferred) {
    preferredContactMethod = 'Phone, Email';
  } else if (phoneIsPreferred) {
    preferredContactMethod = 'Phone';
  } else if (emailIsPreferred) {
    preferredContactMethod = 'Email';
  }

  return (
    <div className={styles.contactInfoContainer}>
      <div className={styles.contactInfoHeader}>
        <h2>Contact info</h2>
        <Link to={editURL}>Edit</Link>
      </div>

      <div className={styles.contactInfoSection}>
        <dl className={descriptionListStyles.descriptionList}>
          <div className={descriptionListStyles.row}>
            <dt>Best contact phone</dt>
            <dd>{telephone}</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>Alt. phone</dt>
            <dd>{secondaryTelephone || '–'}</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>Personal email</dt>
            <dd>{personalEmail}</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>Preferred contact method</dt>
            <dd>{preferredContactMethod}</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>Current mailing address</dt>
            <dd>
              {residentialAddress.street_address_1} {residentialAddress.street_address_2}
              <br />
              {residentialAddress.city}, {residentialAddress.state} {residentialAddress.postal_code}
            </dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>Backup mailing address</dt>
            <dd>
              {backupMailingAddress.street_address_1} {backupMailingAddress.street_address_2}
              <br />
              {backupMailingAddress.city}, {backupMailingAddress.state} {backupMailingAddress.postal_code}
            </dd>
          </div>
        </dl>
      </div>

      <div className={styles.contactInfoSection}>
        <h3>Backup contact</h3>
        <dl className={descriptionListStyles.descriptionList}>
          <div className={descriptionListStyles.row}>
            <dt>Name</dt>
            <dd>{backupContact.name}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>Email</dt>
            <dd>{backupContact.email}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>Phone</dt>
            <dd>{backupContact.telephone}</dd>
          </div>
        </dl>
      </div>
    </div>
  );
};

ContactInfoDisplay.propTypes = {
  telephone: PropTypes.string.isRequired,
  secondaryTelephone: PropTypes.string,
  personalEmail: PropTypes.string.isRequired,
  phoneIsPreferred: PropTypes.bool,
  emailIsPreferred: PropTypes.bool,
  residentialAddress: ResidentialAddressShape.isRequired,
  backupMailingAddress: ResidentialAddressShape.isRequired,
  backupContact: BackupContactShape.isRequired,
  editURL: PropTypes.string.isRequired,
};

ContactInfoDisplay.defaultProps = {
  secondaryTelephone: '',
  phoneIsPreferred: false,
  emailIsPreferred: false,
};

export default ContactInfoDisplay;
