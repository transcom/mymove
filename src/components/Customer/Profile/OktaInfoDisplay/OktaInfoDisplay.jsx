import React from 'react';
import { string, PropTypes } from 'prop-types';
import { Link } from 'react-router-dom';

import oktaLogo from '../../../../shared/images/okta_logo.png';

import oktaInfoDisplayStyles from './OktaInfoDisplay.module.scss';

import descriptionListStyles from 'styles/descriptionList.module.scss';

const url = 'https://test-milmove.okta.mil/enduser/settings';

const OktaInfoDisplay = ({ editURL, oktaUsername, oktaEmail, oktaFirstName, oktaLastName, oktaEdipi }) => {
  return (
    <div className={oktaInfoDisplayStyles.serviceInfoContainer}>
      <div className={oktaInfoDisplayStyles.header}>
        <a href={url}>
          <img className={oktaInfoDisplayStyles.oktaLogo} src={oktaLogo} alt="Okta logo" />
        </a>
        <Link className={oktaInfoDisplayStyles.oktaEditLink} to={editURL}>
          Edit
        </Link>
      </div>
      <div className={oktaInfoDisplayStyles.header}>
        <p>
          <b>Okta</b> is the identity provider you used when signing into MilMove. If you would like to edit any of this
          information, you can do so by clicking the <b>Edit</b> link.
        </p>
      </div>
      <div className={oktaInfoDisplayStyles.oktaInfoSection}>
        <dl className={descriptionListStyles.descriptionList}>
          <div className={descriptionListStyles.row}>
            <dt>Username</dt>
            <dd>{oktaUsername}</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>Email</dt>
            <dd>{oktaEmail}</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>First Name</dt>
            <dd>{oktaFirstName}</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>Last Name</dt>
            <dd>{oktaLastName}</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>DoD ID Number | EDIPI</dt>
            <dd>{oktaEdipi}</dd>
          </div>
        </dl>
      </div>
    </div>
  );
};

OktaInfoDisplay.propTypes = {
  oktaUsername: PropTypes.string.isRequired,
  oktaEmail: PropTypes.string.isRequired,
  oktaFirstName: PropTypes.string.isRequired,
  oktaLastName: PropTypes.string.isRequired,
  oktaEdipi: PropTypes.string,
  editURL: string,
};

OktaInfoDisplay.defaultProps = {
  editURL: '',
  oktaEdipi: 'Not Provided',
};

export default OktaInfoDisplay;
