import React, { useState } from 'react';
import styles from './LoginButton.module.scss';
import { LogoutUser } from 'utils/api';
import { Button } from '@trussworks/react-uswds';

import ConnectedEulaModal from '../../../components/EulaModal';

export const LoginButton = (props) => {
  const [showEula, setShowEula] = useState(false);

  if (!props.isLoggedIn) {
    return (
      <>
        <ConnectedEulaModal
          isOpen={showEula}
          acceptTerms={() => {
            window.location.href = '/auth/login-gov';
          }}
          closeModal={() => setShowEula(false)}
        />
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
          <Button
            aria-label="Sign In"
            className={styles.signIn}
            data-testid="signin"
            onClick={() => setShowEula(!showEula)}
            type="button"
          >
            Sign In
          </Button>
        </div>
      </>
    );
  } else {
    return (
      <Button aria-label="Sign Out" className="usa=link" data-testid="signout" onClick={LogoutUser} type="button">
        Sign Out
      </Button>
    );
  }
};
