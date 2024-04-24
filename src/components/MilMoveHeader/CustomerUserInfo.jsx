import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { NavLink, useLocation, generatePath } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './index.module.scss';

import { customerRoutes } from 'constants/routes';

const CustomerUserInfo = ({ handleLogout, showProfileLink, moveId }) => {
  let { state } = useLocation();
  state = { ...state, moveId };

  return (
    <div className={styles.userInfo}>
      <ul className="usa-nav__primary">
        {showProfileLink && (
          <li className="usa-nav__primary-item">
            <NavLink
              to={generatePath(customerRoutes.PROFILE_PATH)}
              className={styles.profileLink}
              title="profile-link"
              aria-label="profile-link"
              state={state}
            >
              <FontAwesomeIcon className="fa-2x" icon={['far', 'user']} />
            </NavLink>
          </li>
        )}

        <li className="usa-nav__primary-item">
          <Button unstyled onClick={handleLogout} type="button">
            Sign out
          </Button>
        </li>
      </ul>
    </div>
  );
};

CustomerUserInfo.propTypes = {
  handleLogout: PropTypes.func.isRequired,
  showProfileLink: PropTypes.bool,
};

CustomerUserInfo.defaultProps = {
  showProfileLink: false,
};

export default CustomerUserInfo;
