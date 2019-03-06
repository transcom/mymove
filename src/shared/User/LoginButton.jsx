import React from 'react';
import { connect } from 'react-redux';
import { get } from 'lodash';
import * as Cookies from 'js-cookie';

import { selectCurrentUser } from 'shared/Data/users';
import { isDevelopment } from 'shared/constants';

import './LoginButton.css';

const token = Cookies.get('masked_gorilla_csrf');

const LoginButton = props => {
  if (!props.isLoggedIn) {
    return (
      <React.Fragment>
        {props.showDevlocalButton && (
          <a data-hook="devlocal-signin" style={{ marginRight: '2em' }} href="/devlocal-auth/login">
            Local Sign In
          </a>
        )}
        <a data-hook="signin" href="/auth/login-gov">
          Sign In
        </a>
      </React.Fragment>
    );
  } else {
    return (
      <form className="logoutForm" name="logoutForm" method="post" action="/auth/logout">
        <div className="logout">
          <input type="hidden" name="gorilla.csrf.Token" value={token} />
          <input className="logoutButton" type="submit" value="Sign Out" name="logout" />
        </div>
      </form>
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
