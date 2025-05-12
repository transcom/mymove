import React from 'react';
import PropTypes from 'prop-types';
import { useLocation } from 'react-router-dom';

import styles from './ContactInfoDisplay.module.scss';

import { ResidentialAddressShape } from 'types/address';
import { BackupContactShape } from 'types/customerShapes';
import descriptionListStyles from 'styles/descriptionList.module.scss';
import { LinkButton } from 'shared/standardUI/Buttons';

const editButtonStyle = styles['edit-btn'];

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

  const { state } = useLocation();

  return (
    <div className={styles.contactInfoContainer}>
      <div className={styles.contactInfoHeader}>
        <h2>Contact info</h2>
        <LinkButton className={editButtonStyle} href={editURL} state={state}>
          Edit
        </LinkButton>
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
            <dt>Current address</dt>
            <dd>
              {residentialAddress.streetAddress1} {residentialAddress.streetAddress2}{' '}
              {residentialAddress.streetAddress3}
              <br />
              {residentialAddress.city}, {residentialAddress.state} {residentialAddress.postalCode}
            </dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>Backup address</dt>
            <dd>
              {backupMailingAddress.streetAddress1} {backupMailingAddress.streetAddress2}{' '}
              {backupMailingAddress.streetAddress3}
              <br />
              {backupMailingAddress.city}, {backupMailingAddress.state} {backupMailingAddress.postalCode}
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
