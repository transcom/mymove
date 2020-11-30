import React from 'react';
import { string, bool } from 'prop-types';

import styles from './Contact.module.scss';

const Contact = ({ header, dutyStationName, moveSubmitted, officeType, telephone }) => (
  <div className={styles.contactContainer}>
    <h6 className={styles.contactHeader}>{header}</h6>
    <p>
      <strong>{dutyStationName}</strong>
      <br />
      <span>{officeType}</span>
      <br />
      <span>{telephone}</span>
    </p>
    {moveSubmitted && (
      <p data-testid="move-submitted-instructions">
        Talk to your move counselor or directly with your movers if you have questions during your move.
      </p>
    )}
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
