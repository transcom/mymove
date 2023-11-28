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
      You have been logged out of Okta. If you need to sign in again, you can do so by clicking <strong>Sign in</strong>
      . If you have any other issues logging in or authenticating with Okta, please refer to our troubleshooting page
      here:{' '}
      <a
        className={styles.link}
        href="https://transcom.github.io/mymove-docs/docs/getting-started/okta/okta-troubleshooting"
      >
        <strong>Okta Troubleshooting Guide</strong>
      </a>
      . If you need to log out of the Okta Dashboard to completely clear your session, you can do so{' '}
      <a className={styles.link} href={oktaURL}>
        <strong>here</strong>
      </a>
      .
    </div>
  );
};

export default OktaLogoutBanner;
