import React from 'react';
import lockImage from 'shared/images/lock.png';
const SignIn = () => (
  <div className="usa-grid">
    <div className="usa-width-one-sixth">&nbsp; </div>
    <div className="usa-width-two-thirds">
      <div className="align-center">
        <p>
          <img src={lockImage} width="120" alt="" />
        </p>
        <br />
        <p>
          Your account security is important. Please sign in before proceeding.
        </p>
        <br />
        <a href="/auth/login-gov" className="usa-button  usa-button-big ">
          Sign in
        </a>
        <br /> <br />
      </div>
      <div className="storage-estimate">
        <p>
          Welcome to my.move.mil! This is a new system put together by
          USTRANSCOM to support the relocation of families during the PCS
          process.{' '}
        </p>
        <p>
          Right now, use of this system is by invitation only. If you haven't
          received an invitation, you need to go to{' '}
          <a href="https://eta.sddc.army.mil/ETASSOPortal/default.aspx">DPS</a>{' '}
          to schedule your move.
        </p>
        <p>
          Over the coming months we'll be rolling this new tool out to more and
          more people. Stay tuned.
        </p>
      </div>
    </div>
  </div>
);

export default SignIn;
