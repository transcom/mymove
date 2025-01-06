import React from 'react';

import sadComputer from 'shared/images/sad-computer.png';

const Inaccessible = () => (
  <div className="usa-grid">
    <div className="usa-width-one-whole align-center">
      <p>
        <img src={sadComputer} width="200" height="200" alt="" />
      </p>
      <h2>Page is not accessible.</h2>
      <p data-testid="contactMsg">
        If you feel this message was received in error, please call (800) 462-2176, Option 2 or{' '}
        <a href="mailto:usarmy.scott.sddc.mbx.G6-SRC-MilMove-HD@army.mil" className="usa-link">
          email us
        </a>
        .
      </p>
    </div>
  </div>
);

export const INACCESSIBLE_API_RESPONSE = 'Page is inaccessible';

export default Inaccessible;
