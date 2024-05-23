import React from 'react';

import sadComputer from 'shared/images/sad-computer.png';

const SomethingWentWrong = () => (
  <div className="usa-grid">
    <div className="usa-width-one-whole align-center">
      <p>
        <img src={sadComputer} width="200" height="200" alt="" />
      </p>
      <h2>
        Oops!
        <br />
        Something went wrong.
      </h2>
      <p>Please try again in a few moments.</p>
      <br />
      <br />
      <br />
      <p data-testid="contactMsg">
        If you continue to receive this error, call (800) 462-2176, Option 2 or{' '}
        <a href="mailto:usarmy.scott.sddc.mbx.G6-SRC-MilMove-HD@army.mil" className="usa-link">
          email us
        </a>
        .
      </p>
    </div>
  </div>
);

export default SomethingWentWrong;
