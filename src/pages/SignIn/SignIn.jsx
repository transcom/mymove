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
import { isDevelopment } from 'shared/constants';
import { useTitle } from 'hooks/custom';
import ConnectedFlashMessage from 'containers/FlashMessage/FlashMessage';

const SignIn = ({ context, showLocalDevLogin, showTestharnessList }) => {
  const location = useLocation();
  const [showEula, setShowEula] = useState(false);
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

  return (
    <div className={classNames(styles.center, 'usa-prose grid-container padding-top-3')}>
      <ConnectedEulaModal
        isOpen={showEula}
        acceptTerms={() => {
          window.location.href = '/auth/okta';
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
                onClick={() => setShowEula(!showEula)}
                type="button"
              >
                Sign in
              </Button>

              {showLocalDevLogin && (
                <a className="usa-button" data-testid="devlocal-signin" href="/devlocal-auth/login">
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
