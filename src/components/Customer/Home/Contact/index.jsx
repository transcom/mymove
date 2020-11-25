import React from 'react';
import { string, bool } from 'prop-types';

import styles from './Contact.module.scss';

const Contact = ({ header, dutyStationName, moveSubmitted, officeType, telephone }) => (
  <div className={styles['contact-container']}>
    <h6 className={styles['contact-header']}>{header}</h6>
    <p>
      <strong>{dutyStationName}</strong>
      <br />
      <span>{officeType}</span>
      <br />
      <span>{telephone}</span>
    </p>
    {moveSubmitted && <p>After you hear from your move counselor, they should be your first resource for questions.</p>}
  </div>
);

Contact.propTypes = {
  dutyStationName: string.isRequired,
  header: string.isRequired,
  moveSubmitted: bool.isRequired,
  officeType: string.isRequired,
  telephone: string.isRequired,
};

export default Contact;
