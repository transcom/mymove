import React from 'react';
import PropTypes from 'prop-types';
import { Button, GridContainer, Grid } from '@trussworks/react-uswds';

import contactInfoStyles from './ContactInfo.module.scss';

import { ResidentialAddressShape } from 'types/address';
import { BackupContactShape } from 'types/customerShapes';

const ContactInfo = ({
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
  if (phoneIsPreferred) {
    preferredContactMethod = 'Phone';
  } else if (emailIsPreferred) {
    preferredContactMethod = 'Email';
  }

  return (
    <GridContainer className={contactInfoStyles['contact-info-container']}>
      <Grid row>
        <Grid col className={contactInfoStyles['contact-info-header']}>
          <h2>Contact info</h2>
          <Button
            unstyled
            className={contactInfoStyles['edit-btn']}
            data-testid="edit-contact-info"
            onClick={onEditClick}
          >
            Edit
          </Button>
        </Grid>
      </Grid>

      <Grid row>
        <Grid col className={contactInfoStyles['contact-info-section']}>
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
        <Grid col className={contactInfoStyles['contact-info-section']}>
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

ContactInfo.propTypes = {
  telephone: PropTypes.string.isRequired,
  secondaryTelephone: PropTypes.string,
  personalEmail: PropTypes.string.isRequired,
  phoneIsPreferred: PropTypes.bool.isRequired,
  emailIsPreferred: PropTypes.bool.isRequired,
  residentialAddress: ResidentialAddressShape.isRequired,
  backupMailingAddress: ResidentialAddressShape.isRequired,
  backupContact: BackupContactShape.isRequired,
  onEditClick: PropTypes.func.isRequired,
};

ContactInfo.defaultProps = {
  secondaryTelephone: '',
};

export default ContactInfo;
