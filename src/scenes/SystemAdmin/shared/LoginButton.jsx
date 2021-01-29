import React from 'react';
import styles from './LoginButton.module.scss';
import { LogoutUser } from 'utils/api';

export const LoginButton = (props) => {
  if (!props.isLoggedIn) {
    return (
      <div className={styles['login-section']}>
        {props.showDevlocalButton && (
          <a
            data-hook="devlocal-signin"
            style={{ marginRight: '2em' }}
            href="/devlocal-auth/login"
            className="usa-link"
          >
            Local Sign In
          </a>
        )}
        <a data-hook="signin" href="/auth/login-gov" className="usa-link">
          Sign In
        </a>
      </div>
    );
  } else {
    return (
      <a href="#" onClick={LogoutUser} className="usa-link">
        Sign Out
      </a>
    );
  }
};
