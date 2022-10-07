import React from 'react';
import sadComputer from 'shared/images/sad-computer.png';

const SomethingWentWrong = ({ error, info }) => (
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
      <div>
        error: {error.toString()}
        info: {info.componentStack}
      </div>
      <br />
      <br />
      <br />
      <br />
      <p>
        If you continue to receive this error, call (302) 4MY-MOVE or{' '}
        <a href="mailto:transcom.scott.tcj5j4.mbx.ppcf@mail.mil" className="usa-link">
          email us
        </a>
        .
      </p>
    </div>
  </div>
);

export default SomethingWentWrong;
