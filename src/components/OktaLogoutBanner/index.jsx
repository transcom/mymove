import React from 'react';

import styles from './OktaLogoutBanner.module.scss';

const OktaLogoutBanner = () => {
  const hostname = window && window.location && window.location.hostname;
  const oktaURL =
    hostname === 'office.move.mil' || hostname === 'admin.move.mil'
      ? 'https://milmove.okta.mil/enduser/settings'
      : 'https://test-milmove.okta.mil/enduser/settings';

  return (
    <div className={styles.oktaLogoutBanner} data-testid="okta-logout-banner">
      You have been logged out of Okta. <br />
      If you need to sign in again, you can do so by clicking <strong>Sign in</strong> below. <br />
      If you have any other issues logging in or authenticating with Okta, please refer to our{' '}
      <a
        className={styles.link}
        href="https://transcom.github.io/mymove-docs/docs/getting-started/okta/okta-troubleshooting"
      >
        <strong>troubleshooting page</strong>
      </a>
      . <br />
      If you continue to have issues authenticating, please go{' '}
      <a className={styles.link} href={oktaURL}>
        <strong>here</strong>
      </a>{' '}
      and sign completely out of Okta and try logging into MilMove again.
    </div>
  );
};

export default OktaLogoutBanner;
