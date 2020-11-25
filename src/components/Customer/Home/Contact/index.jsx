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
    <p>
      {moveSubmitted && (
        <>
          You&apos;ll hear from a move counselor and the movers themselves within the next few days.
          <br />
          <br />
          Talk to either of them with questions about your move.
        </>
      )}
    </p>
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
