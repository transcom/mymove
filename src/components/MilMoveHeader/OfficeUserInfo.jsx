import React from 'react';
import { string, func } from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { Link } from 'react-router-dom';

import styles from './index.module.scss';

import { officeRoutes } from 'constants/routes';

const OfficeUserInfo = ({ handleLogout, firstName, lastName }) => {
  return (
    <div className={styles.userInfo}>
      <ul className="usa-nav__primary">
        <li className="usa-nav__primary-item userLink">
          <Link to={officeRoutes.PROFILE_PATH} title="profile-link">
            {lastName}
            {lastName && firstName && ', '}
            {firstName}
          </Link>
        </li>
        <li className="usa-nav__primary-item">
          <Button unstyled onClick={handleLogout} type="button">
            Sign out
          </Button>
        </li>
      </ul>
    </div>
  );
};

OfficeUserInfo.defaultProps = {
  firstName: null,
  lastName: null,
};

OfficeUserInfo.propTypes = {
  firstName: string,
  lastName: string,
  handleLogout: func.isRequired,
};

export default OfficeUserInfo;
