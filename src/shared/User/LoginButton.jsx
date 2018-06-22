import React from 'react';
import { connect } from 'react-redux';
import { get } from 'lodash';

import { isDevelopment } from 'shared/constants';

const LoginButton = props => {
  if (!props.isLoggedIn) {
    return (
      <React.Fragment>
        {props.showDevlocalButton && (
          <a
            data-hook="devlocal-signin"
            style={{ marginRight: '2em' }}
            href="/devlocal-auth/login"
          >
            Local Sign In
          </a>
        )}
        <a data-hook="signin" href="/auth/login-gov">
          Sign In
        </a>
      </React.Fragment>
    );
  } else {
    return <a href="/auth/logout">Sign Out</a>;
  }
};

function mapStateToProps(state) {
  return {
    isLoggedIn: state.user.isLoggedIn,
    showDevlocalButton: get(state, 'isDevelopment', isDevelopment),
  };
}
export default connect(mapStateToProps)(LoginButton);
