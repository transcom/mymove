import React from 'react';
import styles from './LoginButton.module.scss';
import { LogoutUser } from 'shared/User/api.js';

export const LoginButton = props => {
  if (!props.isLoggedIn) {
    return (
      <div className={styles['login-section']}>
        {props.showDevlocalButton && (
          <a data-hook="devlocal-signin" style={{ marginRight: '2em' }} href="/devlocal-auth/login">
            Local Sign In
          </a>
        )}
        <a data-hook="signin" href="/auth/login-gov">
          Sign In
        </a>
      </div>
    );
  } else {
    return (
      <a href="#" onClick={LogoutUser}>
        Sign Out
      </a>
    );
  }
};
