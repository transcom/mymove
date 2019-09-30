import React from 'react';
import { connect } from 'react-redux';
import { get } from 'lodash';

import { selectCurrentUser } from 'shared/Data/users';
import { isDevelopment } from 'shared/constants';
import { LogoutUser } from 'shared/User/api.js';

const LoginButton = props => {
  if (!props.isLoggedIn) {
    return (
      <React.Fragment>
        {props.showDevlocalButton && (
          <a
            className="usa-nav__link"
            data-hook="devlocal-signin"
            style={{ marginRight: '2em' }}
            href="/devlocal-auth/login"
          >
            Local Sign In
          </a>
        )}
        <a className="usa-nav__link" data-hook="signin" href="/auth/login-gov">
          Sign In
        </a>
      </React.Fragment>
    );
  } else {
    return (
      <a className="usa-nav__link" href="#" onClick={LogoutUser}>
        Sign Out
      </a>
    );
  }
};

function mapStateToProps(state) {
  const user = selectCurrentUser(state);
  return {
    isLoggedIn: user.isLoggedIn,
    showDevlocalButton: get(state, 'isDevelopment', isDevelopment),
  };
}
export default connect(mapStateToProps)(LoginButton);
