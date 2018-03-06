import React from 'react';
import { connect } from 'react-redux';

const LoginButton = props => {
  if (!props.loggedIn) return <a href="/auth/login-gov">Sign In</a>;
  else return <a href="/auth/logout">Log Out</a>;
};

function mapStateToProps(state) {
  return {
    loggedIn: state.user.loggedIn,
  };
}
//export default connect(mapStateToProps)(LoginButton);
export default LoginButton;
