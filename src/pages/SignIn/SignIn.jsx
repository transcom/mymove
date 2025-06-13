import React, { useEffect, useState } from 'react';
import qs from 'query-string';
import { bool, shape, string } from 'prop-types';
import { Button, ButtonGroup } from '@trussworks/react-uswds';
import { useLocation, useNavigate } from 'react-router-dom';
import '../../styles/signinImports.scss';
import classNames from 'classnames';

import styles from './SignIn.module.scss';

import '@trussworks/react-uswds/lib/index.css';

import { withContext } from 'shared/AppContext';
import Alert from 'shared/Alert';
import ConnectedEulaModal from 'components/EulaModal';
import { FEATURE_FLAG_KEYS, isDevelopment } from 'shared/constants';
import { useTitle } from 'hooks/custom';
import ConnectedFlashMessage from 'containers/FlashMessage/FlashMessage';
import { isBooleanFlagEnabledUnauthenticated } from 'utils/featureFlags';
import { generalRoutes } from 'constants/routes';

const SignIn = ({ context, showLocalDevLogin, showTestharnessList }) => {
  const location = useLocation();
  const [showEula, setShowEula] = useState(false);
  const [isSigningIn, setIsSigningIn] = useState(false);
  const [isSigningUp, setIsSigningUp] = useState(false);
  const [customerRegistrationFF, setCustomerRegistrationFF] = useState(false);

  const navigate = useNavigate();

  const { error } = qs.parse(location.search);
  const { siteName, showLoginWarning } = context;

  useTitle();

  useEffect(() => {
    function unload() {
      navigate('', { replace: true, state: null });
    }
    window.addEventListener('beforeunload', unload);
    return () => window.removeEventListener('beforeunload', unload);
  }, [navigate]);

  useEffect(() => {
    isBooleanFlagEnabledUnauthenticated(FEATURE_FLAG_KEYS.CUSTOMER_REGISTRATION)?.then((enabled) => {
      setCustomerRegistrationFF(enabled);
    });
  }, []);

  return (
    <div className={classNames(styles.center, 'usa-prose grid-container padding-top-3')}>
      <ConnectedEulaModal
        isOpen={showEula}
        acceptTerms={() => {
          if (isSigningIn) {
            window.location.href = '/auth/okta';
          } else if (isSigningUp) {
            navigate(generalRoutes.CREATE_ACCOUNT_PATH);
          }
        }}
        closeModal={() => setShowEula(false)}
      />
      <div className="grid-row">
        <div>
          {error && (
            <div>
              <Alert type="error" heading="An error occurred">
                There was an error during your last sign in attempt. Please try again.
              </Alert>
              <br />
            </div>
          )}
          {location.state && location.state.timedout && (
            <div>
              <Alert type="error" heading="Logged out">
                You have been logged out due to inactivity.
              </Alert>
            </div>
          )}
          {location.state && location.state.hasLoggedOut && (
            <div>
              <Alert type="success" heading="You have signed out of MilMove">
                Sign in again when you&apos;re ready to start a new session.
              </Alert>
            </div>
          )}
          {location.state && location.state.noValidCAC && (
            <div>
              <Alert type="warning" heading="CAC Validation is required at first sign-in">
                If you do not have a Common Access Card (CAC) do not request your account here. You must visit your
                nearest personal property office where they will assist you with creating your MilMove account.
              </Alert>
            </div>
          )}

          {siteName === 'office.move.mil' && <ConnectedFlashMessage />}

          <h1 className="align-center">Welcome to {siteName}!</h1>
          {showLoginWarning && (
            <div>
              <h2 style={{ color: 'red' }}>
                Use of this system is by invitation only, following mandatory screening for{' '}
                <a
                  href="https://dps.move.mil/cust/standard/user/home.xhtml"
                  style={{ color: 'red', textDecoration: 'underline' }}
                >
                  eligibility in MilMove.
                </a>{' '}
              </h2>
              <h2 style={{ color: 'red' }}>
                DO NOT PROCEED if you have not gone through that{' '}
                <a
                  href="https://dps.move.mil/cust/standard/user/home.xhtml"
                  style={{ color: 'red', textDecoration: 'underline' }}
                >
                  screening process which begins with you selecting &quot;New Shipment&quot; (click here to begin).
                </a>{' '}
              </h2>
              <h2 style={{ color: 'red' }}>
                Failure to do so will likely result in you having to resubmit your shipment in the{' '}
                <a
                  style={{ color: 'red', textDecoration: 'underline' }}
                  href="https://dps.move.mil/cust/standard/user/home.xhtml"
                >
                  Defense Personal Property System
                </a>{' '}
                and could cause a delay in your shipment being moved.
              </h2>
            </div>
          )}
          <div className="align-center">
            <ButtonGroup type="default">
              <Button
                aria-label="Sign In"
                className={siteName === 'my.move.mil' ? styles.signInButton : 'usa-button'}
                data-testid="signin"
                onClick={() => {
                  setIsSigningUp(false);
                  setIsSigningIn(true);
                  setShowEula(!showEula);
                }}
                type="button"
              >
                Sign in
              </Button>
              {siteName === 'my.move.mil' && customerRegistrationFF ? (
                <Button
                  aria-label="Create account"
                  className={siteName === 'my.move.mil' ? styles.signInButton : 'usa-button'}
                  data-testid="createAccount"
                  onClick={() => {
                    setIsSigningIn(false);
                    setIsSigningUp(true);
                    setShowEula(!showEula);
                  }}
                  type="button"
                >
                  Create Account
                </Button>
              ) : null}

              {showLocalDevLogin && (
                <a className="usa-button usa-button--primary" data-testid="devlocal-signin" href="/devlocal-auth/login">
                  Local Sign In
                </a>
              )}
              {showTestharnessList && (
                <a className="usa-button" data-testid="devlocal-testharnesslist" href="/testharness/list">
                  View Testharness Data Scenarios
                </a>
              )}
            </ButtonGroup>
          </div>
        </div>
      </div>
    </div>
  );
};

SignIn.propTypes = {
  context: shape({
    siteName: string,
    showLoginWarning: bool,
  }).isRequired,
  showLocalDevLogin: bool,
  showTestharnessList: bool,
};

SignIn.defaultProps = {
  showLocalDevLogin: isDevelopment,
  showTestharnessList: isDevelopment,
};

export default withContext(SignIn);
