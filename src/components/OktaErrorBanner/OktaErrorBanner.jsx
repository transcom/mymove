import React from 'react';

import styles from './OktaErrorBanner.module.scss';

const OktaErrorBanner = () => {
  return (
    <div className={styles.oktaErrorBanner} data-testid="okta-error-banner">
      You must use a different e-mail when authenticating with Okta.
      <br />
      Access to this application is denied with the previously used authentication method.
    </div>
  );
};

export default OktaErrorBanner;
