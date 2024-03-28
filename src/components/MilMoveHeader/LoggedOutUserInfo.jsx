import React from 'react';
import { func } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from './index.module.scss';

const LoggedOutUserInfo = ({ handleLogin, handleRequestAccount, app }) => {
  return (
    <>
      {app === 'office' && (
        <div className={styles.userInfo}>
          <ul className="usa-nav__primary">
            <li className="usa-nav__primary-item">
              <Button
                unstyled
                aria-label="Request Account"
                onClick={handleRequestAccount}
                data-testid="requestaccount"
                type="button"
              >
                Request Account
              </Button>
            </li>
          </ul>
        </div>
      )}
      <div className={styles.userInfo}>
        <ul className="usa-nav__primary">
          <li className="usa-nav__primary-item">
            <Button unstyled aria-label="Sign In" onClick={handleLogin} data-testid="signin" type="button">
              Sign in
            </Button>
          </li>
        </ul>
      </div>
    </>
  );
};

LoggedOutUserInfo.propTypes = {
  handleLogin: func.isRequired,
};

export default LoggedOutUserInfo;
