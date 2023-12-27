import React from 'react';

import styles from './OktaLogoutBanner.module.scss';

export const OktaNeedsLoggedOutBanner = () => {
  const hostname = window && window.location && window.location.hostname;
  const oktaURL =
    hostname === 'office.move.mil' || hostname === 'admin.move.mil'
      ? 'https://milmove.okta.mil/enduser/settings'
      : 'https://test-milmove.okta.mil/enduser/settings';

  return (
    <div className={styles.oktaNeedsLoggedOutBanner} data-testid="okta-logout-banner">
      You have an existing Okta session. Please log out of Okta completely.
      <br />
      <a className={styles.link} href={oktaURL} target="_blank" rel="noreferrer">
        <strong>You can access your Okta dashboard by following this link.</strong>
      </a>{' '}
      <br />
      In the top-right corner, you can click the drop down where it displays your name and select &apos;Sign Out&apos;.{' '}
      <br />
      Once you sign out of Okta, you should be able to sign into MilMove.
      <br />
      If you have issues logging in or authenticating with Okta, please refer to our{' '}
      <a
        className={styles.link}
        target="_blank"
        href="https://transcom.github.io/mymove-docs/docs/getting-started/okta/okta-troubleshooting"
        rel="noreferrer"
      >
        <strong>troubleshooting page</strong>
      </a>
      .
    </div>
  );
};

export const OktaLoggedOutBanner = () => {
  const hostname = window && window.location && window.location.hostname;
  const oktaURL =
    hostname === 'office.move.mil' || hostname === 'admin.move.mil'
      ? 'https://milmove.okta.mil/enduser/settings'
      : 'https://test-milmove.okta.mil/enduser/settings';

  return (
    <div className={styles.oktaLoggedOutBanner} data-testid="okta-logout-banner">
      You have been logged out of Okta. To sign in again, click <strong>Sign in</strong> below. <br />
      If you have any other issues please refer to our{' '}
      <a
        className={styles.link}
        target="_blank"
        href="https://transcom.github.io/mymove-docs/docs/getting-started/okta/okta-troubleshooting"
        rel="noreferrer"
      >
        <strong>troubleshooting page</strong>
      </a>
      . <br />
      If you continue to have issues, please go to your{' '}
      <a className={styles.link} href={oktaURL} target="_blank" rel="noreferrer">
        <strong>Okta dashboard</strong>
      </a>{' '}
      and sign completely out, then try logging into MilMove again.
    </div>
  );
};
