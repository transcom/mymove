import React from 'react';
import PropTypes from 'prop-types';
import { Button, GridContainer, Grid } from '@trussworks/react-uswds';

import styles from './ContactInfoDisplay.module.scss';

import { ResidentialAddressShape } from 'types/address';
import { BackupContactShape } from 'types/customerShapes';

const ContactInfoDisplay = ({
  telephone,
  secondaryTelephone,
  personalEmail,
  phoneIsPreferred,
  emailIsPreferred,
  residentialAddress,
  backupMailingAddress,
  backupContact,
  onEditClick,
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
    <GridContainer className={styles['contact-info-container']}>
      <Grid row>
        <Grid col className={styles['contact-info-header']}>
          <h2>Contact info</h2>
          <Button unstyled className={styles['edit-btn']} data-testid="edit-contact-info" onClick={onEditClick}>
            Edit
          </Button>
        </Grid>
      </Grid>

      <Grid row>
        <Grid col className={styles['contact-info-section']}>
          <dl>
            <dt>Best contact phone</dt>
            <dd>{telephone}</dd>

            <dt>Alt. phone</dt>
            <dd>{secondaryTelephone || '–'}</dd>

            <dt>Personal email</dt>
            <dd>{personalEmail}</dd>

            <dt>Preferred contact method</dt>
            <dd>{preferredContactMethod}</dd>

            <dt>Current mailing address</dt>
            <dd>
              {residentialAddress.street_address_1} {residentialAddress.street_address_2}
              <br />
              {residentialAddress.city}, {residentialAddress.state} {residentialAddress.postal_code}
            </dd>

            <dt>Backup mailing address</dt>
            <dd>
              {backupMailingAddress.street_address_1} {backupMailingAddress.street_address_2}
              <br />
              {backupMailingAddress.city}, {backupMailingAddress.state} {backupMailingAddress.postal_code}
            </dd>
          </dl>
        </Grid>
      </Grid>

      <Grid row>
        <Grid col className={styles['contact-info-section']}>
          <h3>Backup contact</h3>
          <dl>
            <dt>Name</dt>
            <dd>{backupContact.name}</dd>

            <dt>Email</dt>
            <dd>{backupContact.email}</dd>

            <dt>Phone</dt>
            <dd>{backupContact.telephone}</dd>
          </dl>
        </Grid>
      </Grid>
    </GridContainer>
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
  onEditClick: PropTypes.func.isRequired,
};

ContactInfoDisplay.defaultProps = {
  secondaryTelephone: '',
  phoneIsPreferred: false,
  emailIsPreferred: false,
};

export default ContactInfoDisplay;
