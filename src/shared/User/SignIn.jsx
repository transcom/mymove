import React from 'react';
import qs from 'query-string';

import { withContext } from 'shared/AppContext';
import Alert from 'shared/Alert';

const SignIn = ({ context, location }) => {
  const error = qs.parse(location.search).error;
  return (
    <div className="usa-grid">
      <div className="usa-width-one-sixth">&nbsp; </div>
      <div className="usa-width-two-thirds">
        {error && (
          <div>
            <Alert type="error" heading="An error occurred">
              There was an error during your last sign in attempt. Please try again.<br />
              Error code: {error}
            </Alert>
            <br />
          </div>
        )}
        <h2 className="align-center">Welcome to {context.siteName}!</h2>
        <br />
        <p>This is a new system from USTRANSCOM to support the relocation of families during PCS.</p>
        {context.showLoginWarning && (
          <div>
            <p>
              Right now, use of this system is by invitation only. If you haven't received an invitation, please go to{' '}
              <a href="https://eta.sddc.army.mil/ETASSOPortal/default.aspx">DPS</a> to schedule your move.
            </p>
            <p>Over the coming months, we'll be rolling this new tool out to more and more people. Stay tuned.</p>
          </div>
        )}
        <div className="align-center">
          <a href="/auth/login-gov" className="usa-button  usa-button-big ">
            Sign in
          </a>
        </div>
      </div>
    </div>
  );
};

export default withContext(SignIn);
