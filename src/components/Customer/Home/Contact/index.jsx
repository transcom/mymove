import React from 'react';
import { Link as USWDSLink } from '@trussworks/react-uswds';
import { string } from 'prop-types';

import styles from './Contact.module.scss';

const Contact = ({ header, dutyLocationName, officeType, telephone }) => (
  <div className={styles.contactContainer}>
    <h6 className={styles.contactHeader}>{header}</h6>
    <p>
      {dutyLocationName && (
        <>
          <strong>{dutyLocationName}</strong>
          <br />
        </>
      )}
      {officeType && (
        <>
          <span>{officeType}</span>
          <br />
        </>
      )}
      {telephone && <span>{telephone}</span>}
    </p>
    <p>
      For government support or information, consult Military OneSource&apos;s{' '}
      <USWDSLink
        target="_blank"
        rel="noopener noreferrer"
        href="https://www.militaryonesource.mil/moving-housing/moving/planning-your-move/customer-service-contacts-for-military-pcs/"
      >
        directory of PCS-related contacts
      </USWDSLink>
      .
    </p>
    <p>If you&apos;re using government movers, contact them directly for questions about your shipments.</p>
  </div>
);

Contact.propTypes = {
  dutyLocationName: string.isRequired,
  header: string.isRequired,
  officeType: string.isRequired,
  telephone: string.isRequired,
};

export default Contact;
