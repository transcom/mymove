import React from 'react';
import { string, bool } from 'prop-types';
import { Link } from 'react-router-dom';

import oktaLogo from '../../../../shared/images/okta_logo.png';

import oktaInfoDisplayStyles from './OktaInfoDisplay.module.scss';

import descriptionListStyles from 'styles/descriptionList.module.scss';

const url = 'https://test-milmove.okta.mil/enduser/settings';

const OktaInfoDisplay = ({
  affiliation,
  originDutyLocationName,
  originTransportationOfficeName,
  originTransportationOfficePhone,
  edipi,
  firstName,
  isEditable,
  showMessage,
  lastName,
  editURL,
  rank,
}) => {
  return (
    <div className={oktaInfoDisplayStyles.serviceInfoContainer}>
      <div className={oktaInfoDisplayStyles.header}>
        <a href={url}>
          <img className={oktaInfoDisplayStyles.oktaLogo} src={oktaLogo} alt="Okta logo" />
        </a>
        <Link to={editURL}>Edit</Link>
      </div>
      <div className={oktaInfoDisplayStyles.serviceInfoSection}>
        <dl className={descriptionListStyles.descriptionList}>
          <div className={descriptionListStyles.row}>
            <dt>Username</dt>
            <dd>oktausername@email.com</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>Email</dt>
            <dd>oktaEmail@email.com</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>First Name</dt>
            <dd>First Name</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>Last Name</dt>
            <dd>Last Name</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>DoD ID Number</dt>
            <dd>DoDID or &apos;Not Provided&apos;</dd>
          </div>
        </dl>
      </div>
    </div>
  );
};

OktaInfoDisplay.propTypes = {
  affiliation: string.isRequired,
  originDutyLocationName: string.isRequired,
  originTransportationOfficeName: string.isRequired,
  originTransportationOfficePhone: string,
  edipi: string.isRequired,
  firstName: string.isRequired,
  isEditable: bool,
  showMessage: bool,
  lastName: string.isRequired,
  editURL: string,
  rank: string.isRequired,
};

OktaInfoDisplay.defaultProps = {
  originTransportationOfficePhone: '',
  editURL: '',
  isEditable: true,
  showMessage: false,
};

export default OktaInfoDisplay;
