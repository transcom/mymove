import React from 'react';
import { string } from 'prop-types';

import styles from './Contact.module.scss';

const Contact = ({ header, dutyStationName, officeType, telephone }) => (
  <div className={`${styles['contact-container']} padding-top-2 padding-left-2 padding-right-2 padding-bottom-3`}>
    <h6 className="margin-bottom-1">{header.toUpperCase()}</h6>
    <p>
      <strong>{dutyStationName}</strong>
      <br />
      <span>{officeType}</span>
      <br />
      <span>{telephone}</span>
    </p>
  </div>
);

Contact.propTypes = {
  header: string.isRequired,
  dutyStationName: string.isRequired,
  officeType: string.isRequired,
  telephone: string.isRequired,
};

export default Contact;
