import React from 'react';
import sadComputer from 'shared/images/sad-computer.png';
const FailWhale = () => (
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
      <p>
        If you continue to receive this error, call (302) 4MY-MOVE or{' '}
        <a href="mailto:transcom.scott.tcj5j4.mbx.ppcf@mail.mil">email us</a>.
      </p>
    </div>
  </div>
);

export default FailWhale;
