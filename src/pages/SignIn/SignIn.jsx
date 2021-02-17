import React, { useState } from 'react';
import qs from 'query-string';
import { bool, shape, string } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from './SignIn.module.scss';
import '@trussworks/react-uswds/lib/index.css';

import { withContext } from 'shared/AppContext';
import Alert from 'shared/Alert';
import ConnectedEulaModal from 'components/EulaModal';
import { LocationShape } from 'types/index';

const SignIn = ({ context, location }) => {
  const [showEula, setShowEula] = useState(false);

  const { error } = qs.parse(location.search);
  const hash = qs.parse(location.hash);
  const { siteName, showLoginWarning } = context;

  return (
    <div className="grid-container usa-prose">
      <ConnectedEulaModal
        isOpen={showEula}
        acceptTerms={() => {
          window.location.href = '/auth/login-gov';
        }}
        closeModal={() => setShowEula(false)}
      />
      <div className="grid-row">
        <div className="grid-col-8 grid-offset-2">
          {error && (
            <div>
              <Alert type="error" heading="An error occurred">
                There was an error during your last sign in attempt. Please try again.
              </Alert>
              <br />
            </div>
          )}
          {'timedout' in hash && (
            <div>
              <Alert type="error" heading="Logged out">
                You have been logged out due to inactivity.
              </Alert>
            </div>
          )}

          <h1 className="align-center">Welcome to {siteName}!</h1>
          <p>This is a new system from USTRANSCOM to support the relocation of families during PCS.</p>
          {showLoginWarning && (
            <div>
              <p>
                Right now, use of this system is by invitation only. If you haven&apos;t received an invitation, please
                go to{' '}
                <a href="https://eta.sddc.army.mil/ETASSOPortal/default.aspx" className="usa-link">
                  DPS
                </a>{' '}
                to schedule your move.
              </p>
              <p>
                Over the coming months, we&apos;ll be rolling this new tool out to more and more people. Stay tuned.
              </p>
            </div>
          )}
          <div className="align-center">
            <Button
              aria-label="Sign In"
              className={siteName === 'my.move.mil' ? styles.signInButton : 'usa-button'}
              onClick={() => setShowEula(!showEula)}
              type="button"
            >
              Sign in
            </Button>
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
  location: LocationShape.isRequired,
};

export default withContext(SignIn);
