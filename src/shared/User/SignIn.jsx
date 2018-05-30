import React from 'react';
import lockImage from 'shared/images/lock.png';
const SignIn = () => (
  <div className="usa-grid">
    <div className="usa-width-one-sixth">&nbsp; </div>
    <div className="usa-width-two-thirds">
      <h2 className="align-center">Welcome to my.move.mil!</h2>
      <br />
      <p>
        This is a new system from USTRANSCOM to support the relocation of
        families during PCS.
      </p>
      <p>
        Right now, use of this system is by invitation only. If you haven't
        received an invitation, please go to DPS to schedule your move.
      </p>
      <p>
        Over the coming months, we'll be rolling this new tool out to more and
        more people. Stay tuned.
      </p>
      <div className="align-center">
        <a href="/auth/login-gov" className="usa-button  usa-button-big ">
          Sign in
        </a>
      </div>
    </div>
  </div>
);

export default SignIn;
