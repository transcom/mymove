import React from 'react';
import qs from 'query-string';
import { bool, shape, string } from 'prop-types';
import { Button, ButtonGroup } from '@trussworks/react-uswds';
import { useHistory } from 'react-router-dom-old';

import '../../styles/office.scss';
import styles from './InvalidPermissions.module.scss';

import '@trussworks/react-uswds/lib/index.css';

import { LogoutUser } from 'utils/api';
import { logOut } from 'store/auth/actions';
import { withContext } from 'shared/AppContext';
import Alert from 'shared/Alert';
import { LocationShape } from 'types/index';

const InvalidPermissions = ({ context, location }) => {
  const history = useHistory();
  const { siteName } = context;
  const { traceId } = qs.parse(location.search);
  const signoutClass = siteName === 'my.move.mil' ? styles.signInButton : 'usa-button';

  const handleLogOut = () => {
    logOut();
    LogoutUser().then((r) => {
      const redirectURL = r.body;
      if (redirectURL) {
        window.location.href = redirectURL;
      } else {
        history.push({
          pathname: '/sign-in',
          state: { hasLoggedOut: true },
        });
      }
    });
  };

  return (
    <div className="usa-prose grid-container padding-top-3">
      <div className="grid-row">
        <div>
          <h1>You do not have permission to access this site.</h1>
          <p>
            You are not signed in with a role that gives you access. If you believe you should have access, contact your
            administrator.
          </p>
          {traceId && traceId !== '' && (
            <Alert type="warning" slim>
              If you believe this is an error, try logging out and back in.
              <br />
              <br />
              If that doesn&apos;t work, please contact the{' '}
              <a className={styles.link} href="https://move.mil/customer-service#technical-help-desk">
                Technical Help Desk
              </a>{' '}
              and give them this code: <strong>{traceId}</strong>
            </Alert>
          )}
          <ButtonGroup type="default">
            <Button
              aria-label="Sign Out"
              className={signoutClass}
              data-testid="signout"
              onClick={handleLogOut}
              type="button"
            >
              Sign Out
            </Button>
          </ButtonGroup>
        </div>
      </div>
    </div>
  );
};

InvalidPermissions.propTypes = {
  context: shape({
    siteName: string,
    showLoginWarning: bool,
  }).isRequired,
  location: LocationShape.isRequired,
};

InvalidPermissions.defaultProps = {};

export default withContext(InvalidPermissions);
