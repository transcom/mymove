import { Button } from '@trussworks/react-uswds';
import React from 'react';

import smartCard from 'shared/images/smart-card.png';
import styles from './SmartCardRedirect.module.scss';
import { logOut } from 'store/auth/actions';
import { LogoutUserWithOktaRedirect } from 'utils/api';
import { useNavigate } from 'react-router';

const SmartCardRedirect = () => {
  const navigate = useNavigate();

  const handleLogout = () => {
    logOut();
    LogoutUserWithOktaRedirect().then((r) => {
      const redirectURL = r.body;
      const urlParams = new URLSearchParams(redirectURL.split('?')[1]);
      const idTokenHint = urlParams.get('id_token_hint');
      if (redirectURL && idTokenHint !== 'devlocal') {
        window.location.href = redirectURL;
      } else {
        navigate('/sign-in');
      }
    });
  };

  return (
    <div className="usa-grid">
      <div className="usa-width-one-whole align-center">
        <p>
          <img src={smartCard} width="200" height="200" alt="" />
        </p>
        <h2>Please sign in with your smart card.</h2>
        <p data-testid="helperText">
          For your first MilMove visit, you must sign in with your smart card.
          <br />
          For future visits, you can continue to authenticate using your smart card or use any other method you have set
          up.
          <br />
          Click the button to sign out and authenticate with your smart card.
        </p>
        <div className={styles.signOutBtn}>
          <Button onClick={handleLogout}>Sign Out</Button>
        </div>
        <br />
        <br />
        <p data-testid="contactMsg">
          If you have already authenticated with your smart card once and are still seeing this message,
          <br /> call (800) 462-2176 choose Option 2 or{' '}
          <a href="mailto:usarmy.scott.sddc.mbx.G6-SRC-MilMove-HD@army.mil" className="usa-link">
            email us
          </a>
          .
        </p>
      </div>
    </div>
  );
};

export default SmartCardRedirect;
