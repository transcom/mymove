import React, { useEffect, useState } from 'react';
import qs from 'query-string';
import { Button } from '@trussworks/react-uswds';

import { withContext } from 'shared/AppContext';
import Alert from 'shared/Alert';
import ConnectedEulaModal from '../../../components/EulaModal';
import styles from './SignIn.module.scss';

const SignIn = ({ context, location }) => {
  const error = qs.parse(location.search).error;
  const [showEula, setShowEula] = useState(false);
  function recentlyLoggedOut() {
    window.localStorage.removeItem('hasLoggedOut');
  }
  useEffect(() => {
    window.addEventListener('beforeunload', recentlyLoggedOut);
    return () => window.removeEventListener('beforeunload', recentlyLoggedOut);
  }, []);
  const hasLoggedOut = JSON.parse(window.localStorage.getItem('hasLoggedOut'));
  return (
    <div>
      <ConnectedEulaModal
        isOpen={showEula}
        acceptTerms={() => {
          window.location.href = '/auth/login-gov';
        }}
        closeModal={() => setShowEula(false)}
      />
      <div>&nbsp;</div>
      <div>
        {error && (
          <div>
            <Alert type="error" heading="An error occurred">
              There was an error during your last sign in attempt. Please try again.
            </Alert>
            <br />
          </div>
        )}
        {hasLoggedOut && (
          <div>
            <Alert type="success" heading="You have signed out of MilMove">
              Sign in again when you&apos;re ready to start a new session.
            </Alert>
          </div>
        )}
        <h1 className="align-center">Welcome to {context.siteName}!</h1>
        <p className="align-center">
          This is a new system from USTRANSCOM to support the relocation of families during PCS.
        </p>
        <div className="align-center">
          <Button
            aria-label="Sign In"
            className={styles['usa-button']}
            data-testid="signin"
            onClick={() => setShowEula(!showEula)}
            type="button"
          >
            Sign in
          </Button>
        </div>
      </div>
    </div>
  );
};

export default withContext(SignIn);
