import React from 'react';
import { connect } from 'react-redux';

const LoginButton = props => {
  if (!props.isLoggedIn) return <a href="/auth/login-gov">Sign In</a>;
  else return <a href="/auth/logout">Log Out</a>;
};

function mapStateToProps(state) {
  return {
    isLoggedIn: state.user.isLoggedIn,
  };
}
export default connect(mapStateToProps)(LoginButton);
