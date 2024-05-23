import React from 'react';

import smartCard from 'shared/images/smart-card.png';

const SmartCardRedirect = () => (
  <div className="usa-grid">
    <div className="usa-width-one-whole align-center">
      <p>
        <img src={smartCard} width="200" height="200" alt="" />
      </p>
      <h2>You must sign in with your smart card first.</h2>
      <p data-testid="helperText">
        Please sign out and authenticate with your smart card.
        <br />
        Once you sign in with your smart card, it is an optional authentication method going forward.
        <br />
        You can then use any other authenticator you have set up.
      </p>
      <br />
      <br />
      <br />
      <p data-testid="contactMsg">
        If you continue to receive this error even after authenticating with a smart card,
        <br /> call (800) 462-2176, Option 2 or{' '}
        <a href="mailto:usarmy.scott.sddc.mbx.G6-SRC-MilMove-HD@army.mil" className="usa-link">
          email us
        </a>
        .
      </p>
    </div>
  </div>
);

export default SmartCardRedirect;
