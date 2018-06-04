import React from 'react';
import { connect } from 'react-redux';
import { isDevelopment } from 'shared/constants';

const LoginButton = props => {
  if (!props.isLoggedIn) {
    return (
      <React.Fragment>
        {isDevelopment && (
          <a style={{ marginRight: '2em' }} href="/devlocal-auth/login">
            Local Sign In
          </a>
        )}
        <a href="/auth/login-gov">Sign In</a>
      </React.Fragment>
    );
  } else {
    return <a href="/auth/logout">Sign Out</a>;
  }
};

function mapStateToProps(state) {
  return {
    isLoggedIn: state.user.isLoggedIn,
  };
}
export default connect(mapStateToProps)(LoginButton);
