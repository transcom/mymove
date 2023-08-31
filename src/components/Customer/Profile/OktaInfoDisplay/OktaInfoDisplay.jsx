import React from 'react';
import { string } from 'prop-types';
import { Link } from 'react-router-dom';

import oktaLogo from '../../../../shared/images/okta_logo.png';

import oktaInfoDisplayStyles from './OktaInfoDisplay.module.scss';

import descriptionListStyles from 'styles/descriptionList.module.scss';

const url = 'https://test-milmove.okta.mil/enduser/settings';

const OktaInfoDisplay = ({ editURL, username, email, firstName, lastName, edipi }) => {
  return (
    <div className={oktaInfoDisplayStyles.serviceInfoContainer}>
      <div className={oktaInfoDisplayStyles.header}>
        <a href={url}>
          <img className={oktaInfoDisplayStyles.oktaLogo} src={oktaLogo} alt="Okta logo" />
        </a>
        <Link to={editURL}>Edit</Link>
      </div>
      <div className={oktaInfoDisplayStyles.header}>
        <p>
          <b>Okta</b> is the identity provider you used when signing into MilMove. If you would like to edit any of this
          information, you can do so by clicking the <b>Edit</b> link above.
        </p>
      </div>
      <div className={oktaInfoDisplayStyles.serviceInfoSection}>
        <dl className={descriptionListStyles.descriptionList}>
          <div className={descriptionListStyles.row}>
            <dt>Username</dt>
            <dd>{username}</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>Email</dt>
            <dd>{email}</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>First Name</dt>
            <dd>{firstName}</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>Last Name</dt>
            <dd>{lastName}</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>DoD ID Number | EDIPI</dt>
            <dd>{edipi}</dd>
          </div>
        </dl>
      </div>
    </div>
  );
};

OktaInfoDisplay.propTypes = {
  editURL: string,
};

OktaInfoDisplay.defaultProps = {
  editURL: '',
};

export default OktaInfoDisplay;
